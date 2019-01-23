package schemas

import (
	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincodalba-ms/mutations"
	"github.com/jal88/elrincodalba-ms/queries"
)

var Article, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queries.Article,
		Mutation: mutations.Article,
	},
)
