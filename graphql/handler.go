package graphql

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"fmt"

	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	defs "github.com/j13v/elrincondalba-ms/definitions"
	"github.com/j13v/elrincondalba-ms/pubsub"
)

const (
	ContentTypeJSON               = "application/json"
	ContentTypeGraphQL            = "application/graphql"
	ContentTypeFormURLEncoded     = "application/x-www-form-urlencoded"
	ContentTypeMultiplartFormData = "multipart/form-data"
)

// AuthenticateFunc is a function that resolves an auth token
// into a user (or returns an error if that isn't possible).
type AuthenticateFunc func(token string) (interface{}, error)

// RootObjectFn allows a user to generate a RootObject per request
type RootObjectFn func(ctx context.Context, r *http.Request) map[string]interface{}

type ResultCallbackFn func(ctx context.Context, params *graphql.Params, result *graphql.Result, responseBody []byte)

// HandlerConfig stores the configuration of a GraphQL handler.
type HandlerConfig struct {
	SubscriptionManager graphqlws.SubscriptionManager
	Authenticate        AuthenticateFunc
	Secret 				string
	Schema              *graphql.Schema
	Context             context.Context
}

type Handler struct {
	Schema           *graphql.Schema
	pretty           bool
	rootObjectFn     RootObjectFn
	resultCallbackFn ResultCallbackFn
	formatErrorFn    func(err error) gqlerrors.FormattedError
}

// type Handler struct {
// 	MaxBodySize int64 // in bytes
// 	Executor    Executor
// 	Client      bool
// }

type Request struct {
	OperationName string                 `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
	Context       context.Context
}

type RequestOptions struct {
	Query         string                 `json:"query" url:"query" schema:"query"`
	Variables     map[string]interface{} `json:"variables" url:"variables" schema:"variables"`
	OperationName string                 `json:"operationName" url:"operationName" schema:"operationName"`
}

// a workaround for getting`variables` as a JSON string
type requestOptionsCompatibility struct {
	Query         string `json:"query" url:"query" schema:"query"`
	Variables     string `json:"variables" url:"variables" schema:"variables"`
	OperationName string `json:"operationName" url:"operationName" schema:"operationName"`
}

func getFromInterface(operations interface{}) *RequestOptions {
	reqOpt := RequestOptions{}
	switch data := operations.(type) {
	case map[string]interface{}:
		if value, ok := data["operationName"]; ok && value != nil {
			reqOpt.OperationName = value.(string)
		}
		if value, ok := data["query"]; ok && value != nil {
			reqOpt.Query = value.(string)
		}
		if value, ok := data["variables"]; ok && value != nil {
			reqOpt.Variables = value.(map[string]interface{})
		}
	}
	return &reqOpt
}

func getFromForm(values url.Values) *RequestOptions {
	query := values.Get("query")
	if query != "" {
		// get variables map
		variables := make(map[string]interface{}, len(values))
		variablesStr := values.Get("variables")
		json.Unmarshal([]byte(variablesStr), &variables)

		return &RequestOptions{
			Query:         query,
			Variables:     variables,
			OperationName: values.Get("operationName"),
		}
	}

	return nil
}

func set(v interface{}, m interface{}, path string) error {
	var parts []interface{}
	for _, p := range strings.Split(path, ".") {
		if isNumber, err := regexp.MatchString(`\d+`, p); err != nil {
			return err
		} else if isNumber {
			index, _ := strconv.Atoi(p)
			parts = append(parts, index)
		} else {
			parts = append(parts, p)
		}
	}
	for i, p := range parts {
		last := i == len(parts)-1
		switch idx := p.(type) {
		case string:
			if last {
				m.(map[string]interface{})[idx] = v
			} else {
				m = m.(map[string]interface{})[idx]
			}
		case int:
			if last {
				m.([]interface{})[idx] = v
			} else {
				m = m.([]interface{})[idx]
			}
		}
	}
	return nil
}

// NewRequestOptions Parses a http.Request into GraphQL request options struct
func NewRequestOptions(r *http.Request) *RequestOptions {
	if reqOpt := getFromForm(r.URL.Query()); reqOpt != nil {
		return reqOpt
	}

	if r.Method != http.MethodPost {
		return &RequestOptions{}
	}

	if r.Body == nil {
		return &RequestOptions{}
	}

	// TODO: improve Content-Type handling
	contentTypeStr := r.Header.Get("Content-Type")
	contentTypeTokens := strings.Split(contentTypeStr, ";")
	contentType := contentTypeTokens[0]
	switch contentType {
	case ContentTypeMultiplartFormData:
		// Parse multipart form
		// TODO Get MaxBodySize from config
		if err := r.ParseMultipartForm(1024); err != nil {
			panic(err)
		}

		// Unmarshal uploads
		var uploads = map[defs.File][]string{}
		var uploadsMap = map[string][]string{}
		if err := json.Unmarshal([]byte(r.Form.Get("map")), &uploadsMap); err != nil {
			panic(err)
		} else {

			for key, path := range uploadsMap {
				file, header, err := r.FormFile(key)
				// defer file.Close()
				if err != nil {
					panic(err)
				} else {
					uploads[defs.File{
						File:     file,
						Size:     header.Size,
						Filename: header.Filename,
					}] = path
				}
			}
		}
		var operations interface{}
		reqOpt := &RequestOptions{}

		// Unmarshal operations
		if err := json.Unmarshal([]byte(r.Form.Get("operations")), &operations); err != nil {
			panic(err)
		}

		reqOpt = getFromInterface(operations)

		for file, paths := range uploads {
			for _, path := range paths {
				set(file, operations, path)
			}
		}
		return reqOpt
	case ContentTypeGraphQL:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return &RequestOptions{}
		}
		return &RequestOptions{
			Query: string(body),
		}
	case ContentTypeFormURLEncoded:
		if err := r.ParseForm(); err != nil {
			return &RequestOptions{}
		}

		if reqOpt := getFromForm(r.PostForm); reqOpt != nil {
			return reqOpt
		}

		return &RequestOptions{}

	case ContentTypeJSON:
		fallthrough
	default:
		var opts RequestOptions
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return &opts
		}
		err = json.Unmarshal(body, &opts)
		if err != nil {
			// Probably `variables` was sent as a string instead of an object.
			// So, we try to be polite and try to parse that as a JSON string
			var optsCompatible requestOptionsCompatibility
			json.Unmarshal(body, &optsCompatible)
			json.Unmarshal([]byte(optsCompatible.Variables), &opts.Variables)
		}
		return &opts
	}
}

// NewHandler creates a WebSocket handler for GraphQL WebSocket connections.
// This handler takes a SubscriptionManager and adds/removes subscriptions
// as they are started/stopped by the client.
func NewHandlerFunc(config HandlerConfig) func(http.ResponseWriter, *http.Request) {
	if config.Secret == "" {
		panic(fmt.Errorf("auth: secret must be defined"))
	}

	subscriptionManager := graphqlws.NewSubscriptionManager(config.Schema)
	// Create a WebSocket/HTTP handler
	graphqlwsHandler := pubsub.NewHandlerFunc(graphqlws.HandlerConfig{
		// Wire up the GraphqL WebSocket handler with the subscription manager
		SubscriptionManager: subscriptionManager,

		// Optional: Add a hook to resolve auth tokens into users that are
		// then stored on the GraphQL WS connections
		Authenticate: func(authToken string) (interface{}, error) {
			// This is just a dumb example
			return "Joe", nil
		},
	})

	// subscriptions := subscriptionManager.Subscriptions()

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Connection"][0] == "Upgrade" {
			graphqlwsHandler(w, r)
			return
		}
		// get query
		opts := NewRequestOptions(r)

		token := ""
		if r.Header.Get("Authorization") != "" {
			token = strings.SplitN(r.Header.Get("Authorization"), " ", 2)[1]
		}
		// Join contexts
		ctx := context.WithValue(config.Context, "token", token)
		ctx = context.WithValue(ctx, "secret", "secret")
		// execute graphql query
		params := graphql.Params{
			Schema:         *config.Schema,
			RequestString:  opts.Query,
			VariableValues: opts.Variables,
			OperationName:  opts.OperationName,
			Context:        ctx,
		}

		// if h.rootObjectFn != nil {
		// 	params.RootObject = h.rootObjectFn(ctx, r)
		// }
		result := graphql.Do(params)
		//
		// if formatErrorFn := h.formatErrorFn; formatErrorFn != nil && len(result.Errors) > 0 {
		// 	formatted := make([]gqlerrors.FormattedError, len(result.Errors))
		// 	for i, formattedError := range result.Errors {
		// 		formatted[i] = formatErrorFn(formattedError.OriginalError())
		// 	}
		// 	result.Errors = formatted
		// }

		// use proper JSON Header
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		// json.NewEncoder(w).Encode(result)
		var buff []byte
		// if h.pretty {
		// 	w.WriteHeader(http.StatusOK)
		// 	buff, _ = json.MarshalIndent(result, "", "\t")
		//
		// 	w.Write(buff)
		// } else {
		w.WriteHeader(http.StatusOK)
		buff, _ = json.Marshal(result)

		w.Write(buff)
		// }

		// if h.resultCallbackFn != nil {
		// 	h.resultCallbackFn(ctx, &params, result, buff)
		// }
	}
}

// ContextHandler provides an entrypoint into executing graphQL queries with a
// user-provided context.
func (h *Handler) ContextHandler(w http.ResponseWriter, r *http.Request) {

}

//
// type Executor func(request *Request) interface{}
// type Factory func(http.ResponseWriter, *http.Request) interface{}
//
// func New(executor Executor) *Handler {
// 	return &Handler{
// 		MaxBodySize: 1024,
// 		Executor:    executor,
// 	}
// }
//
//
// func NewHandlerFunc(config graphqlws.HandlerConfig) func(http.ResponseWriter, *http.Request) bool {
//     w.Header().Set("Content-Type", "application/json; charset=utf-8")
//     var operations interface{}
//     return func(w http.ResponseWriter, r *http.Request) bool {
//
//     }
// }
//
// func (self *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json; charset=utf-8")
//
// 	switch r.Method; {
//   case "GET":
//     request := Request{Context: r.Context()}
//   case "POST":
//
//   }
//
// 	if r.Method == "GET" {
// 		request := Request{Context: r.Context()}
//
// 		// Get query
// 		if value := r.URL.Query().Get("query"); len(value) == 0 {
// 			message := fmt.Sprintf("Missing query")
// 			http.Error(w, message, http.StatusBadRequest)
// 			return
// 		} else {
// 			request.Query = value
// 		}
//
// 		// Get variables
// 		if value := r.URL.Query().Get("variables"); len(value) == 0 {
// 			request.Variables = map[string]interface{}{}
// 		} else if err := json.Unmarshal([]byte(value), &request.Variables); err != nil {
// 			message := fmt.Sprintf("Bad variables")
// 			http.Error(w, message, http.StatusBadRequest)
// 			return
// 		}
//
// 		// Get variables
// 		if value := r.URL.Query().Get("operationName"); len(value) == 0 {
// 			request.OperationName = ""
// 		} else {
// 			request.OperationName = value
// 		}
// 		result := self.Executor(&request)
// 		if err := json.NewEncoder(w).Encode(result); err != nil {
// 			panic(err)
// 		}
// 	} else if r.Method == "POST" {
// 		contentType := strings.SplitN(r.Header.Get("Content-Type"), ";", 2)[0]
//
// 		switch contentType {
// 		case "text/plain", "application/json":
// 			if err := json.NewDecoder(r.Body).Decode(&operations); err != nil {
// 				panic(err)
// 			}
// case "multipart/form-data":
// 	// Parse multipart form
// 	if err := r.ParseMultipartForm(self.MaxBodySize); err != nil {
// 		panic(err)
// 	}
//
// 	// Unmarshal uploads
// 	var uploads = map[File][]string{}
// 	var uploadsMap = map[string][]string{}
// 	if err := json.Unmarshal([]byte(r.Form.Get("map")), &uploadsMap); err != nil {
// 		panic(err)
// 	} else {
// 		for key, path := range uploadsMap {
// 			if file, header, err := r.FormFile(key); err != nil {
// 				panic(err)
// 				//w.WriteHeader(http.StatusInternalServerError)
// 				//return
// 			} else {
// 				uploads[File{
// 					File:     file,
// 					Size:     header.Size,
// 					Filename: header.Filename,
// 				}] = path
// 			}
// 		}
// 	}
//
// 			// Unmarshal operations
// 			if err := json.Unmarshal([]byte(r.Form.Get("operations")), &operations); err != nil {
// 				panic(err)
// 			}
//
// 			// set uploads to operations
// 			for file, paths := range uploads {
// 				for _, path := range paths {
// 					set(file, operations, path)
// 				}
// 			}
// 		}
// 		switch data := operations.(type) {
// 		case map[string]interface{}:
// 			request := Request{}
// 			if value, ok := data["operationName"]; ok && value != nil {
// 				request.OperationName = value.(string)
// 			}
// 			if value, ok := data["query"]; ok && value != nil {
// 				request.Query = value.(string)
// 			}
// 			if value, ok := data["variables"]; ok && value != nil {
// 				request.Variables = value.(map[string]interface{})
// 			}
// 			request.Context = r.Context()
// 			if err := json.NewEncoder(w).Encode(self.Executor(&request)); err != nil {
// 				panic(err)
// 			}
// 		case []interface{}:
// 			result := make([]interface{}, len(data))
// 			for index, operation := range data {
// 				data := operation.(map[string]interface{})
// 				request := Request{}
// 				if value, ok := data["operationName"]; ok {
// 					request.OperationName = value.(string)
// 				}
// 				if value, ok := data["query"]; ok {
// 					request.Query = value.(string)
// 				}
// 				if value, ok := data["variables"]; ok {
// 					request.Variables = value.(map[string]interface{})
// 				}
// 				request.Context = r.Context()
// 				result[index] = self.Executor(&request)
// 			}
// 			if err := json.NewEncoder(w).Encode(result); err != nil {
// 				panic(err)
// 			}
// 		default:
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 	}
