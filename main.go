package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jal88/elrincondalba-ms/graphql"
	decs "github.com/jal88/elrincondalba-ms/graphql/decorators"
	"github.com/jal88/elrincondalba-ms/graphql/schema"
	"github.com/jal88/elrincondalba-ms/logger"
	"github.com/jal88/elrincondalba-ms/mongodb"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/gridfs"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logger.NewLogger("server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, "mongodb://localhost:27017")
	if err != nil {
		fmt.Print(err)
	}

	httpPort := 8080

	db := client.Database("elrincondalba")
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		panic(err)
	}
	// Primary data initialization
	if os.Getenv("INIT_DATABASE") != "" {
		mongodb.InitData(db)
	}
	repo := mongodb.CreateRepo(db)

	ctx = context.Background()
	ctx = decs.ContextRepoApply(repo)(ctx)

	rtr := mux.NewRouter()
	rtr.HandleFunc("/images/{image:[a-z0-9]+}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		image := params["image"]
		oid, err := primitive.ObjectIDFromHex(image)
		if err != nil {
			panic(err)
		}
		downstream, err := bucket.OpenDownloadStream(oid)
		if err != nil {
			panic(err)
		}

		io.Copy(w, downstream)
		defer downstream.Close()
	}))

	rtr.HandleFunc("/graphql", graphql.NewHandlerFunc(graphql.HandlerConfig{
		Schema:  &schema.Schema,
		Context: ctx,
	}))

	rtr.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))

	http.Handle("/", rtr)

	logger.WithFields(logrus.Fields{
		"port": "8080",
		"host": "0.0.0.0",
	}).Info("Server is running")

	err = http.ListenAndServe(fmt.Sprintf(":%d", httpPort), logRequest(http.DefaultServeMux, logger))

	if err != nil {
		log.Fatal(err)
	}
}

func logRequest(handler http.Handler, logger *logrus.Entry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"remote": r.RemoteAddr,
			"host":   r.URL,
		}).Info("Request")
		handler.ServeHTTP(w, r)
	})
}

//////// GRAPHiQL ////////
var page = []byte(`
<!DOCTYPE html>
<html>
    <head>
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.css" />
				<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.43.0/theme/monokai.css" />
        <script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.1.0/fetch.min.js"></script>
        <script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"></script>
        <script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"></script>
        <script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.js"></script>
    </head>
    <body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
        <div id="graphiql" style="height: 100vh;">Loading...</div>
        <script>
            function graphQLFetcher(graphQLParams) {
                return fetch("/graphql", {
                    method: "post",
                    body: JSON.stringify(graphQLParams),
                    credentials: "include",
                }).then(function (response) {
                    return response.text();
                }).then(function (responseBody) {
                    try {
                        return JSON.parse(responseBody);
                    } catch (error) {
                        return responseBody;
                    }
                });
            }
            ReactDOM.render(
                React.createElement(GraphiQL, {
									fetcher: graphQLFetcher,
									editorTheme: "monokai",
									defaultQuery: [
										"query {",
										"  getArticle(id:\"put_here_objectId\"){",
										"    ...ArticleFields",
										"  }",
										"  getOrder(id:\"put_here_objectId\"){",
										"    ...OrderFields",
										"  }",
										"  listArticles(last: 2) {",
										"    ...ArticleConnectionFields",
										"  }",
										"  listOrders(first: 2) {",
										"    ...OrderConnectionFields",
										"  }",
										"}",
										"fragment OrderConnectionFields on OrderConnection {",
										"  totalCount",
										"    edges {",
										"      cursor",
										"      node {",
										"        ...OrderFields",
										"        stock {",
										"          ...StockFields",
										"        }",
										"        user {",
										"          ...UserFields",
										"        }",
										"      }",
										"    }",
										"    pageInfo {",
										"      ...PageInfoFields",
										"    }",
										"}",
										"fragment ArticleConnectionFields on ArticleConnection {",
										"  totalCount",
										"    edges {",
										"      cursor",
										"      node {",
										"        ...ArticleFields",
										"      }",
										"    }",
										"    pageInfo {",
										"      ...PageInfoFields",
										"    }",
										"}",
										"fragment ArticleFields on Article {",
										"  id",
										"  name",
										"  category",
										"  description",
										"  images",
										"  stock {",
										"    count",
										"    size",
										"    refs",
										"  }",
										"  price",
										"  rating",
										"  createdAt",
										"  updatedAt",
										"}",
										"fragment StockFields on Stock {",
										"  article {",
										"    ...ArticleFields",
										"  }",
										"  createdAt",
										"  id",
										"  size",
										"}",
										"fragment OrderFields on Order {",
										"  id",
										"  state",
										"  createdAt",
										"  updatedAt",
										"  notes",
										"}",
										"fragment UserFields on User {",
										"  id",
										"  address",
										"  dni",
										"  email",
										"  name",
										"  phone",
										"  surname",
										"}",
										"fragment PageInfoFields on PageInfo {",
										"  endCursor",
										"  hasNextPage",
										"  hasPreviousPage",
										"  startCursor",
										"}"].join("\n")
									}),
                document.getElementById("graphiql")
            );
        </script>
    </body>
</html>
`)

// ctx := context.Background()
// ctx = context.WithValue(ctx, "model", models))
// context.WithValue(ctx, "user", users))

// // Create a subscription manager
// subscriptionManager := graphqlws.NewSubscriptionManager(&schema.Schema)
// // Create a WebSocket/HTTP handler
// graphqlwsHandler := pubsub.NewHandlerFunc(graphqlws.HandlerConfig{
// 	// Wire up the GraphqL WebSocket handler with the subscription manager
// 	SubscriptionManager: subscriptionManager,
//
// 	// Optional: Add a hook to resolve auth tokens into users that are
// 	// then stored on the GraphQL WS connections
// 	Authenticate: func(authToken string) (interface{}, error) {
// 		// This is just a dumb example
// 		return "Joe", nil
// 	},
// })
//

//
// subscriptions := subscriptionManager.Subscriptions()

// go func() {
// 	for {
// 		time.Sleep(2 * time.Second)
// 		result := db.Collection("article").FindOne(context.Background(), bson.M{
// 			"rating": bson.M{"$gte": " Math.random()"},
// 		})
// 		article := definitions.Article{}
// 		result.Decode(&article)
// 		article.Rating = (article.Rating + 1) % 5
// 		article.Price = rand.Float64() * 100
// 		repo.Article.Sync(&article)
//
// 		for conn, subs := range subscriptions {
// 			conn.ID()
// 			conn.User()
// 			for _, subscription := range subs {
// 				// Prepare an execution context for running the query
// 				ctx := context.Background()
// 				ctx = decs.ContextPubSubApply(article)(ctx)
// 				// Re-execute the subscription query
// 				params := graphql.Params{
// 					Schema:         schema.Schema, // The GraphQL schema
// 					RequestString:  subscription.Query,
// 					VariableValues: subscription.Variables,
// 					OperationName:  subscription.OperationName,
// 					Context:        ctx,
// 				}
// 				result := graphql.Do(params)
// 				// Send query results back to the subscriber at any point
// 				data := graphqlws.DataMessagePayload{
// 					// Data can be anything (interface{})
// 					Data: result.Data,
// 					// Errors is optional ([]error)
// 					Errors: graphqlws.ErrorsFromGraphQLErrors(result.Errors),
// 				}
// 				subscription.SendData(&data)
// 			}
// 		}
// 	}
// }()
//
// http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
//
// 	if r.Header["Connection"][0] == "Upgrade" {
// 		graphqlwsHandler(w, r)
// 	} else {
// 		ghandler.New(executor)
// 		setupResponse(&w, r)
// 		if (*r).Method == "OPTIONS" {
// 			return
// 		}
// 		decoder := json.NewDecoder(r.Body)
// 		var t BodyQueryMessage
// 		err := decoder.Decode(&t)
// 		if err != nil {
// 			panic(err)
// 		}
//
// 		ctx := context.Background()
// 		ctx = decs.ContextRepoApply(repo)(ctx)
// 		// ctx = decs.DecoratorContextUserApply(u)(ctx)
//
// 		result := executeQuery(
// 			schema.Schema,
// 			t.Query,
// 			t.Variables,
// 			ctx)
// 		json.NewEncoder(w).Encode(result)
// 	}
// })

// http.Handle("/images/{image:[a-z0-9]+}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	fmt.Printf("%v\n", r.URL)
// 	params := mux.Vars(r)
// 	name := params["name"]
// 	w.Write([]byte("Hello " + name))
// }))
//
// http.HandleFunc("/graphql", graphql.NewHandlerFunc(graphql.HandlerConfig{
// 	Schema:  &schema.Schema,
// 	Context: ctx,
// }))
//
// http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	w.Write(page)
// }))
