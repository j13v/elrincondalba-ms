package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincondalba-ms/definitions"
	decs "github.com/jal88/elrincondalba-ms/graphql/decorators"
	"github.com/jal88/elrincondalba-ms/graphql/schema"
	"github.com/jal88/elrincondalba-ms/mongodb"
	"github.com/jal88/elrincondalba-ms/pubsub"
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
	db.Drop(ctx)
	// Primary data initialization
	mongodb.InitData(db)
	repo := mongodb.CreateRepo(db)

	// Create a subscription manager
	subscriptionManager := graphqlws.NewSubscriptionManager(&schema.Schema)
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

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))

	subscriptions := subscriptionManager.Subscriptions()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			result, _ := repo.Article.FindOne(&map[string]interface{}{
				"name": "MiniFalda",
			})
			article := result.(definitions.Article)
			article.Rating = (article.Rating + 1) % 5
			article.Price = rand.Float64() * 100
			repo.Article.Sync(&article)

			for conn, subs := range subscriptions {
				conn.ID()
				conn.User()
				for _, subscription := range subs {
					// Prepare an execution context for running the query
					ctx := context.Background()
					ctx = decs.ContextPubSubApply(article)(ctx)
					// Re-execute the subscription query
					params := graphql.Params{
						Schema:         schema.Schema, // The GraphQL schema
						RequestString:  subscription.Query,
						VariableValues: subscription.Variables,
						OperationName:  subscription.OperationName,
						Context:        ctx,
					}
					result := graphql.Do(params)
					// Send query results back to the subscriber at any point
					data := graphqlws.DataMessagePayload{
						// Data can be anything (interface{})
						Data: result.Data,
						// Errors is optional ([]error)
						Errors: graphqlws.ErrorsFromGraphQLErrors(result.Errors),
					}
					subscription.SendData(&data)
				}
			}
		}
	}()

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Connection"][0] == "Upgrade" {
			graphqlwsHandler(w, r)
		} else {
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

			ctx := context.Background()
			ctx = decs.ContextRepoApply(repo)(ctx)
			// ctx = decs.DecoratorContextUserApply(u)(ctx)

			result := executeQuery(
				schema.Schema,
				t.Query,
				t.Variables,
				ctx)
			json.NewEncoder(w).Encode(result)
		}
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
