package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincondalba-ms/definitions"
	decs "github.com/jal88/elrincondalba-ms/graphql/decorators"
	"github.com/jal88/elrincondalba-ms/graphql/schema"
	"github.com/jal88/elrincondalba-ms/graphql/types"
	"github.com/jal88/elrincondalba-ms/mongodb"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
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

type Post struct {
	ID    int `json:"id"`
	Likes int `json:"count"`
}

type ConnectionACKMessage struct {
	OperationID string `json:"id,omitempty"`
	Type        string `json:"type"`
	Payload     struct {
		Query string `json:"query"`
	} `json:"payload,omitempty"`
}

type Subscriber struct {
	ID            int
	Conn          *websocket.Conn
	RequestString string
	OperationID   string
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

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{"graphql-ws"},
	}
	var articles = []*definitions.Article{
		&definitions.Article{
			ID:          primitive.ObjectID{},
			Rating:      0,
			Name:        "Falda burdeos",
			Description: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
			Price:       33.99,
			Images: []string{
				"637f1f77bcf86cd799439011",
				"637f1f77bcf86cd799439012",
				"637f1f77bcf86cd799439013"},
			Category: "Faldas",
		},
		&definitions.Article{
			ID:          primitive.ObjectID{},
			Rating:      0,
			Name:        "Camiseta de tirante",
			Description: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
			Price:       13.99,
			Images: []string{
				"637f1f77bcf86cd799439011",
				"637f1f77bcf86cd799439012",
				"637f1f77bcf86cd799439013"},
			Category: "Camisetas",
		},
	}
	var subscribers sync.Map

	schm, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"articles": &graphql.Field{
					Type: graphql.NewList(types.TypeArticle),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return articles, nil
					},
				},
			},
		}),
		Subscription: graphql.NewObject(graphql.ObjectConfig{
			Name: "Subscription",
			Fields: graphql.Fields{
				"postLikesSubscribe": &graphql.Field{
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Type: graphql.NewList(types.TypeArticle),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return articles, nil
					},
				},
			},
		}),
	})
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))

	http.HandleFunc("/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("failed to do websocket upgrade: %v", err)
			return
		}
		connectionACK, err := json.Marshal(map[string]string{
			"type": "connection_ack",
		})
		if err != nil {
			log.Printf("failed to marshal ws connection ack: %v", err)
		}
		if err := conn.WriteMessage(websocket.TextMessage, connectionACK); err != nil {
			log.Printf("failed to write to ws connection: %v", err)
			return
		}
		go func() {
			for {
				_, p, err := conn.ReadMessage()
				if websocket.IsCloseError(err, websocket.CloseGoingAway) {
					return
				}
				if err != nil {
					log.Println("failed to read websocket message: %v", err)
					return
				}
				var msg ConnectionACKMessage
				if err := json.Unmarshal(p, &msg); err != nil {
					log.Printf("failed to unmarshal: %v", err)
					return
				}
				if msg.Type == "start" {
					length := 0
					subscribers.Range(func(key, value interface{}) bool {
						length++
						return true
					})
					var subscriber = Subscriber{
						ID:            length + 1,
						Conn:          conn,
						RequestString: msg.Payload.Query,
						OperationID:   msg.OperationID,
					}
					subscribers.Store(subscriber.ID, &subscriber)
				}
			}
		}()
	})
	go func() {
		for {
			time.Sleep(5 * time.Second)
			for _, post := range articles {
				post.Price = post.Price + 0.01
			}
			subscribers.Range(func(key, value interface{}) bool {
				subscriber, ok := value.(*Subscriber)
				if !ok {
					return true
				}
				payload := graphql.Do(graphql.Params{
					Schema:        schm,
					RequestString: subscriber.RequestString,
				})
				message, err := json.Marshal(map[string]interface{}{
					"type":    "data",
					"id":      subscriber.OperationID,
					"payload": payload,
				})
				if err != nil {
					log.Printf("failed to marshal message: %v", err)
					return true
				}
				if err := subscriber.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
					if err == websocket.ErrCloseSent {
						subscribers.Delete(key)
						return true
					}
					log.Printf("failed to write to ws connection: %v", err)
					return true
				}
				return true
			})
		}
	}()

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
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
		repo := mongodb.CreateRepo(db)

		// u := struct{
		// 	GetName()
		// 	GetScopes: func() []string {
		// 		return []string{"admin"}
		// 	},
		// }

		ctx := context.Background()
		ctx = decs.ContextRepoApply(repo)(ctx)
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
