package schemas

import (
	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincondalba-ms/mutations"
	"github.com/jal88/elrincondalba-ms/queries"
)

var Article, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queries.Article,
		Mutation: mutations.Article,
	},
)
