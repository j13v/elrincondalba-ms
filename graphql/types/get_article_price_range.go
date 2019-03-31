package types

import (
	"github.com/graphql-go/graphql"
	decs "github.com/j13v/elrincondalba-ms/graphql/decorators"
	"github.com/j13v/elrincondalba-ms/mongodb"
)

var FieldArticlePriceRange = &graphql.Field{
	Type:        graphql.NewList(graphql.Float),
	Description: "Get min max price",
	Args: graphql.FieldConfigArgument{
		"sizes": &graphql.ArgumentConfig{
			Type: graphql.NewList(graphql.String),
		},
		"categories": &graphql.ArgumentConfig{
			Type: graphql.NewList(graphql.String),
		},
	},
	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
		rangePrice, err := repo.Article.GetPriceRange(&params.Args)
		return []float64{rangePrice.Min, rangePrice.Max}, err
	}),
}
