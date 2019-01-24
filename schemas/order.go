package schemas

import (
	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincondalba-ms/mutations"
	"github.com/jal88/elrincondalba-ms/queries"
)

var Order, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queries.Order,
		Mutation: mutations.Order,
	},
)
