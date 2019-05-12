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
	"github.com/j13v/elrincondalba-ms/graphql"
	decs "github.com/j13v/elrincondalba-ms/graphql/decorators"
	"github.com/j13v/elrincondalba-ms/graphql/schema"
	"github.com/j13v/elrincondalba-ms/graphql/utils"
	"github.com/j13v/elrincondalba-ms/logger"
	"github.com/j13v/elrincondalba-ms/mongodb"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/gridfs"
	"github.com/mongodb/mongo-go-driver/x/network/connstring"
	"github.com/sirupsen/logrus"
)

func corsHandler(h http.Handler, origin string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	}
}

func main() {
	httpPort := os.Getenv("PORT")
	httpCorsOrigin := os.Getenv("CORS_ORIGIN")

	if httpCorsOrigin == "" {
		httpCorsOrigin = "*"
	}

	if httpPort == "" {
		httpPort = "8080"
	}

	mongoURI := os.Getenv("MONGODB_URI")
	mongoDbName := os.Getenv("MONGODB_DB")

	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017/elrincondalba"
	}

	mongoURIParams, err := connstring.Parse(mongoURI)
	if err == nil {
		mongoDbName = mongoURIParams.Database
	}

	logger := logger.NewLogger("server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, mongoURI)
	if err != nil {
		panic(err)
	}

	db := client.Database(mongoDbName)
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		panic(err)
	}

	if os.Getenv("API_SECRET") == "" {
		panic("API_SECRET enviroment var must be defined but empty value was detected")
	}
	if os.Getenv("API_ADMIN_EMAIL") == "" {
		panic("API_ADMIN_EMAIL enviroment var must be defined but empty value was detected")
	}
	if os.Getenv("API_ADMIN_PASSWORD") == "" {
		panic("API_ADMIN_PASSWORD enviroment var must be defined but empty value was detected")
	}
	// Primary data initialization
	if os.Getenv("INIT_DATABASE") != "" {
		mongodb.InitData(db)
	}
	repo := mongodb.CreateRepo(db)

	ctx = context.Background()
	ctx = decs.ContextRepoApply(repo)(ctx)
	ctx = decs.ContextAuthApply()(ctx)

	rtr := mux.NewRouter()
	rtr.HandleFunc("/images/{image:[A-Za-z0-9]+}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		image := params["image"]
		oid, err := primitive.ObjectIDFromHex(utils.Base58ToHex(image))
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
		Secret:  os.Getenv("API_SECRET"),
		Schema:  &schema.Schema,
		Context: ctx,
	}))

	rtr.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))

	http.Handle("/", corsHandler(rtr, httpCorsOrigin))

	logger.WithFields(logrus.Fields{
		"port": httpPort,
		"host": "0.0.0.0",
	}).Info("Server is running")

	err = http.ListenAndServe(fmt.Sprintf(":%s", httpPort), logRequest(http.DefaultServeMux, logger))

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
