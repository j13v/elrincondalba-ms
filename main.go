package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincodalba-ms/schemas"
	"github.com/jal88/elrincodalba-ms/types"
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

func initArticlesData(art *[]types.ArticleMock) {
	article1 := types.ArticleMock{
		ID:          1,
		Name:        "MiniFalda",
		Description: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		Price:       7.99,
		Images:      "dsuihfsuizfdhuishg",
		Category:    "Faldas",
		Rating:      2}
	article2 := types.ArticleMock{
		ID:          2,
		Name:        "Sandalia",
		Description: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		Price:       17.99,
		Images:      "dsuihfsuizfdhuishg",
		Category:    "Zapatos",
		Rating:      4}
	article3 := types.ArticleMock{
		ID:          3,
		Name:        "Camiseta de tirante",
		Description: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		Price:       33.99,
		Images:      "szdfszdfszf",
		Category:    "Camisetas",
		Rating:      3}
	*art = append(*art, article1, article2, article3)
}

func main() {
	// Primary data initialization
	initArticlesData(&types.ArticlesMock)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schemas.Article)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
