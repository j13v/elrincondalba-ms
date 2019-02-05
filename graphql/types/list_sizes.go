package types

import (
	"github.com/graphql-go/graphql"
	decs "github.com/jal88/elrincondalba-ms/graphql/decorators"
	"github.com/jal88/elrincondalba-ms/mongodb"
)

var FieldListSizes = &graphql.Field{
	Type:        graphql.NewList(graphql.String),
	Description: "List the sizes",

	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
		sizes, err := repo.Stock.ListSizes()
		return sizes, err
	}),
}
