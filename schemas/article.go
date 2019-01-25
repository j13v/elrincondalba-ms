package schemas

import (
	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincondalba-ms/helpers"
	"github.com/jal88/elrincondalba-ms/mutations"
	"github.com/jal88/elrincondalba-ms/queries"
)
// queries.Article,
var Article, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    graphql.NewObject(graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: helpers.CombineFields(queries.Article, queries.Order),
		}),
		Mutation: mutations.Article,
	},
)
