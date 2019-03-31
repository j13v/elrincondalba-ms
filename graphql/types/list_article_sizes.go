package types

import (
	"github.com/graphql-go/graphql"
	decs "github.com/j13v/elrincondalba-ms/graphql/decorators"
	"github.com/j13v/elrincondalba-ms/mongodb"
)

var FieldListArticleSizes = &graphql.Field{
	Description: "List the article sizes using filters",
	Type:        TypeDistinctList,
	Args: graphql.FieldConfigArgument{
		"categories": &graphql.ArgumentConfig{
			Type: graphql.NewList(graphql.String),
		},
		"priceRange": &graphql.ArgumentConfig{
			Type: graphql.NewList(graphql.Float),
		},
	},
	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
		sizes, err := repo.Article.GetSizes(&params.Args)
		return sizes, err
	}),
}
