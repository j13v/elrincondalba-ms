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
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, "mongodb://localhost:27017")
	if err != nil {
		fmt.Print(err)
	}
	db := client.Database("elrincondalba")
	// Primary data initialization
	// mongodb.InitData(db)

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

// ctx := context.Background()
// ctx = context.WithValue(ctx, "model", models))
// context.WithValue(ctx, "user", users))
