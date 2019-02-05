package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincondalba-ms/graphql/decorators"
	"github.com/jal88/elrincondalba-ms/graphql/mutations"
	"github.com/jal88/elrincondalba-ms/graphql/queries"
	"github.com/jal88/elrincondalba-ms/graphql/types"
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
		Subscription: graphql.NewObject(graphql.ObjectConfig{
			Name: "Subscription",
			Fields: graphql.Fields{
				"postLikesSubscribe": &graphql.Field{
					Type: types.TypeArticle,
					Resolve: decorators.ContextPubSubConsumer(func(value interface{}) (interface{}, error) {
						return value, nil
					}),
				},
			},
		}),
	},
)
