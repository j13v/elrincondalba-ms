package types

import (
	"github.com/graphql-go/graphql"
	decs "github.com/jal88/elrincondalba-ms/graphql/decorators"
	"github.com/jal88/elrincondalba-ms/mongodb"
)

var FieldArticlePriceRange = &graphql.Field{
	Type:        graphql.NewList(graphql.Float),
	Description: "Get min max price",
	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
		rangePrice, err := repo.Article.GetPriceRange()
		return rangePrice, err
	}),
}
