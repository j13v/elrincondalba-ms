package types

import (
	"github.com/graphql-go/graphql"
	decs "github.com/j13v/elrincondalba-ms/graphql/decorators"
	"github.com/j13v/elrincondalba-ms/mongodb"
)

var FieldListSizes = &graphql.Field{
	Type: graphql.NewList(graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Pepesito",
			Fields: graphql.Fields{
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"count": &graphql.Field{
					Type: graphql.Int,
				},
			},
		},
	)),
	Args: graphql.FieldConfigArgument{
		"categories": &graphql.ArgumentConfig{
			Type: graphql.NewList(graphql.String),
		},
	},
	Description: "List the sizes",

	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
		sizes, err := repo.Article.GetSizes(&params.Args)
		return sizes, err
	}),
}
