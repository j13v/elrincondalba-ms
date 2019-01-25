package main

import (
	"encoding/json"
	"fmt"
	"time"
	"context"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincondalba-ms/schemas"
	"github.com/jal88/elrincondalba-ms/models"
	"github.com/mongodb/mongo-go-driver/mongo"
)



func executeQuery(
	schema graphql.Schema,
	query string,
	variables map[string]interface{},
	ctx context.Context) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		VariableValues: variables,
		Context: ctx,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

func initArticlesData(db *mongo.Database) {
	article1 := models.TypeArticle{
		Name:        "MiniFalda",
		Description: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		Price:       7.99,
		Images:      []string{
			"607f1f77bcf86cd799439011",
			"607f1f77bcf86cd799439012"},
		Category:    "Faldas",
		Rating:      2}
	article2 := models.TypeArticle{
		Name:        "Sandalia",
		Description: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		Price:       17.99,
		Images:      []string{
			"617f1f77bcf86cd799439011",
			"617f1f77bcf86cd799439012"},
		Category:    "Zapatos",
		Rating:      4}
	article3 := models.TypeArticle{
		Name:        "Camiseta de tirante",
		Description: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		Price:       33.99,
		Images:      []string{
			"637f1f77bcf86cd799439011",
			"637f1f77bcf86cd799439012",
			"637f1f77bcf86cd799439013"},
		Category:    "Camisetas",
		Rating:      3}

	model := models.GetModelArticle(db)
	model.CreateArticle(article1)
	model.CreateArticle(article2)
	model.CreateArticle(article3)
}

func initOrdersData(db *mongo.Database) {
	order1 := models.TypeOrder{
		Article:  "5c4b58c67cbc327aa78383fd",
		User:     "707f1f77bcf86cd799439011",
		Size:     "L",
		CreateAt: 1548427228,
		UpdateAt: 1548427228,
		State:    1}
	order2 := models.TypeOrder{
		Article:  "5c4b58c67cbc327aa78383fd",
		User:     "707f1f77bcf86cd799439012",
		Size:     "XL",
		CreateAt: 1548427221,
		UpdateAt: 1548427221,
		State:    2}
	order3 := models.TypeOrder{
		Article:  "5c4b58c67cbc327aa78383fd",
		User:     "707f1f77bcf86cd799439013",
		Size:     "S",
		CreateAt: 1548422228,
		UpdateAt: 1548427228,
		State:    3}

		model := models.GetModelOrder(db)
		model.CreateOrder(order1)
		model.CreateOrder(order2)
		model.CreateOrder(order3)
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

type BodyQueryMessage struct {
	Query string `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, "mongodb://localhost:27017")
	if err != nil {
		fmt.Print(err);
	}
	db := client.Database("elrincondalba")
	// Primary data initialization
	initArticlesData(db)
	initOrdersData(db)


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
		dbModel := models.CreateDbModel(db)
		result := executeQuery(
			schemas.Article,
			t.Query,
			t.Variables,
			context.WithValue(context.Background(), "dbModel", dbModel))
		json.NewEncoder(w).Encode(result)

	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
