package types

import (
	"github.com/graphql-go/graphql"
	decs "github.com/jal88/elrincondalba-ms/graphql/decorators"
	"github.com/jal88/elrincondalba-ms/mongodb"
)

var FieldGetCategories = &graphql.Field{
	Type:        graphql.NewList(graphql.String),
	Description: "Get categories",

	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
		categories, err := repo.Article.GetCategories()
		return categories, err
	}),
}
