package types

import (
	"github.com/graphql-go/graphql"
	decs "github.com/j13v/elrincondalba-ms/graphql/decorators"
	"github.com/j13v/elrincondalba-ms/mongodb"
)

var FieldListCategories = &graphql.Field{
	Type:        graphql.NewList(graphql.String),
	Description: "List categories",

	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
		categories, err := repo.Article.GetCategories()
		return categories, err
	}),
}
