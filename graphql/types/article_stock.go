package types

import "github.com/graphql-go/graphql"

var TypeArticleStock = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "ArticleStock",
		Description: "ArticleStock",
		Fields: graphql.Fields{
			"size": &graphql.Field{
				Type: graphql.String,
			},
			"id": &graphql.Field{
				Type: ObjectID,
			},
			"createdAt": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
