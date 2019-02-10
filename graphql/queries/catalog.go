package queries

import (
	"github.com/graphql-go/graphql"
	decs "github.com/j13v/elrincondalba-ms/graphql/decorators"
	"github.com/j13v/elrincondalba-ms/mongodb"
)

var QueryCatalog = graphql.Fields{
	"listCatalog": &graphql.Field{
		Type:        graphql.NewList(graphql.String),
		Description: "List catalog",
		Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
			articles, _, err := repo.Stock.FindSlice(&params.Args)
			return articles, err
		}),
	},
}
