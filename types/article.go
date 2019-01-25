package types

import (
	"github.com/graphql-go/graphql"
)

type ArticleMock struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Images      []string  `json:"images"`
	Category    string  `json:"category"`
	Rating      int8    `json:"rating"`
}

var ArticlesMock []ArticleMock

/*
Article
*/
var Article = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Article",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"price": &graphql.Field{
				Type: graphql.Float,
			},
			"images": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"category": &graphql.Field{
				Type: graphql.String,
			},
			"rating": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
