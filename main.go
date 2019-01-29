package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
	decs "github.com/jal88/elrincondalba-ms/graphql/decorators"
	"github.com/jal88/elrincondalba-ms/graphql/schema"
	"github.com/jal88/elrincondalba-ms/mongodb"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func executeQuery(
	schema graphql.Schema,
	query string,
	variables map[string]interface{},
	ctx context.Context) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  query,
		VariableValues: variables,
		Context:        ctx,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

type BodyQueryMessage struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, "mongodb://localhost:27017")
	if err != nil {
		fmt.Print(err)
	}
	db := client.Database("elrincondalba")
	// Primary data initialization
	mongodb.InitData(db)

	http.Handle("/graphiql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		decoder := json.NewDecoder(r.Body)
		var t BodyQueryMessage
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		model := mongodb.CreateModel(db)

		// u := struct{
		// 	GetName()
		// 	GetScopes: func() []string {
		// 		return []string{"admin"}
		// 	},
		// }

		ctx := context.Background()
		ctx = decs.ContextModelApply(model)(ctx)
		// ctx = decs.DecoratorContextUserApply(u)(ctx)

		result := executeQuery(
			schema.Schema,
			t.Query,
			t.Variables,
			ctx)
		json.NewEncoder(w).Encode(result)

	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
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
                return fetch("/", {
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
										"  listArticles(last: 2) {",
										"    totalCount",
										"    edges {",
										"      cursor",
										"      node {",
										"        ...ArticleFields",
										"      }",
										"    }",
										"    pageInfo {",
										"      ...PageInfoFields",
										"    }",
										"  }",
										"  listOrders(first: 2) {",
										"    totalCount",
										"    pageInfo {",
										"      ...PageInfoFields",
										"    }",
										"    edges {",
										"      cursor",
										"      node {",
										"        article {",
										"          ...ArticleFields",
										"        }",
										"        user {",
										"          ...UserFields",
										"        }",
										"        ...OrderFields",
										"      }",
										"    }",
										"  }",
										"}",
										"fragment ArticleFields on Article {",
										"  category",
										"  description",
										"  id",
										"  images",
										"  name",
										"  price",
										"  rating",
										"}",
										"fragment OrderFields on Order {",
										"  id",
										"  size",
										"  createAt",
										"  updateAt",
										"  state",
										"}",
										"fragment UserFields on User {",
										"  id",
										"  address",
										"  dni",
										"  email",
										"  name",
										"  notes",
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
