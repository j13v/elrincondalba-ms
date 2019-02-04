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
		// Subscription: graphql.NewObject(graphql.ObjectConfig{
		// 	Name: "Subscription",
		// 	Fields: graphql.Fields{
		// 		"articleSubscribe": &graphql.Field{
		// 			Type: graphql.NewList(types.TypeArticle),
		// 			Resolve: decorators.ContextRepoConsumer(func(params graphql.ResolveParams, model mongodb.Repo) (interface{}, error) {
		// 				connArgs := relay.NewConnectionArguments(params.Args)
		// 				articles, meta, err := model.Article.FindSlice(&params.Args)
		// 				return utils.ConnectionFromArraySlice(articles, connArgs, meta), err
		// 			}),
		// 		},
		// 	},
		// }),
	},
)
