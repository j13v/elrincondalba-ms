package types

import (
	"github.com/graphql-go/graphql"
	decs "github.com/j13v/elrincondalba-ms/graphql/decorators"
	"github.com/j13v/elrincondalba-ms/mongodb"
)

var FieldListArticleCategories = &graphql.Field{
	Description: "List article categories using filters",
	Type:        TypeDistinctList,
	Args: graphql.FieldConfigArgument{
		"sizes": &graphql.ArgumentConfig{
			Type: graphql.NewList(graphql.String),
		},
		"priceRange": &graphql.ArgumentConfig{
			Type: graphql.NewList(graphql.Float),
		},
	},
	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
		categories, err := repo.Article.GetCategories(&params.Args)
		return categories, err
	}),
}
