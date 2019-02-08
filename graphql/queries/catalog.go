package queries

import (
	"fmt"

	"github.com/graphql-go/graphql"
	decs "github.com/jal88/elrincondalba-ms/graphql/decorators"
	"github.com/jal88/elrincondalba-ms/mongodb"
)

var QueryCatalog = graphql.Fields{
	"listCatalog": &graphql.Field{
		Type:        graphql.NewList(graphql.String),
		Description: "List catalog",
		Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
			articles, meta, err := repo.Stock.FindSlice(&params.Args)
			fmt.Printf("%v %v %v\n", articles, meta, err)
			return articles, err
		}),
	},
}
