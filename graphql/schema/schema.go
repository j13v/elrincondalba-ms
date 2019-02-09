package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/j13v/elrincondalba-ms/graphql/decorators"
	"github.com/j13v/elrincondalba-ms/graphql/mutations"
	"github.com/j13v/elrincondalba-ms/graphql/queries"
	"github.com/j13v/elrincondalba-ms/graphql/types"
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
