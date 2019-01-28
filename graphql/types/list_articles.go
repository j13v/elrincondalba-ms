package types

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	decs "github.com/jal88/elrincondalba-ms/graphql/decorators"
	"github.com/jal88/elrincondalba-ms/graphql/utils"
	"github.com/jal88/elrincondalba-ms/mongodb"
)

var FieldListArticles = &graphql.Field{
	Type:        TypeArticleConnection,
	Description: "Get articles list",
	Args:        relay.ConnectionArgs,
	Resolve: decs.ContextModelConsumer(func(params graphql.ResolveParams, model mongodb.Model) (interface{}, error) {
		connArgs := relay.NewConnectionArguments(params.Args)
		count, err := model.Article.GetCount()
		if err != nil {
			return nil, err
		}
		findArgs := mongodb.NewFindArgs(params.Args, count)
		articles, err := model.Article.FindSlice(findArgs)
		if err != nil {
			return nil, err
		}
		return utils.ConnectionFromArraySlice(articles, connArgs, relay.ArraySliceMetaInfo{
			SliceStart:  0,
			ArrayLength: 100,
		}), err
	}),
}

//
//
// Resolve: func(params graphql.ResolveParams) (interface{}, error) {
//   args := relay.NewConnectionArguments(params.Args)
//   dbModel := params.Context.Value("model").(model.Models)
//   orders, err := dbModel.Order.ListOrders(params.Args)
//   return relay.ConnectionFromArraySlice(orders, args, relay.ArraySliceMetaInfo{
//     SliceStart:  0,
//     ArrayLength: 100,
//   }), err
// },
//
// Type:        graphql.NewList(qltype.ArticleConnection),
// Description: "Get product list",
// Resolve: func(params graphql.ResolveParams) (interface{}, error) {
//   dbModel := params.Context.Value("model").(model.Models)
//   articles, err := dbModel.Article.ListArticles()
//   return articles, err
// },
