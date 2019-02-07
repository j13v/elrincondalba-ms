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
			"count": &graphql.Field{
				Type: graphql.Int,
			},
			"refs": &graphql.Field{
				Type: graphql.NewList(ObjectID),
			},
		},
	},
)
