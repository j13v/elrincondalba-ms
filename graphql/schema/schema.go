package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincondalba-ms/graphql/mutations"
	"github.com/jal88/elrincondalba-ms/graphql/queries"
)

var Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: queries.Root,
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: mutations.Root,
		}),
	},
)
