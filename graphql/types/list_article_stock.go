package types

import (
	"github.com/graphql-go/graphql"
	decs "github.com/j13v/elrincondalba-ms/graphql/decorators"
	"github.com/j13v/elrincondalba-ms/mongodb"
)

// "id": &FieldID,
// "status": &graphql.Field{
//   Type: graphql.Int,
// },
// "size": &graphql.Field{
//   Type: graphql.Int,
// },

var TypeArticleStockOrderItem = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "ArticleStockOrderItem",
		Description: "ArticleStock",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: ObjectID,
			},
			"state": &graphql.Field{
				Type: graphql.Int,
			},
			"order": &graphql.Field{
				Type: ObjectID,
			},
		},
	},
)

var TypeArticleStockOrder = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "ArticleStockOrder",
		Description: "ArticleStock",
		Fields: graphql.Fields{
			"size": &graphql.Field{
				Type: graphql.String,
			},
			"count": &graphql.Field{
				Type: graphql.Int,
			},
			"refs": &graphql.Field{
				Type: graphql.NewList(TypeArticleStockOrderItem),
			},
		},
	},
)

var FieldArticleStock = &graphql.Field{
	Type: graphql.NewList(TypeArticleStockOrder),
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: ObjectID,
		},
	},
	Description: "Get article stock list",
	Resolve: decs.ContextAuthIsAdmin(decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
		// stock, err := repo.Article.FindByArticle(params.Args["id"], false)
		return nil, nil
		// return claims, nil
	})),
	// Resolve: decs.ContextRepoConsumer(
	// 	[]string{"pepe"},
	// 	decs.ContextRepoConsumer(
	// 		func(params graphql.ResolveParams, model mongodb.Repo) (interface{}, error) {
	// 			stock, err := model.Stock.FindByArticle(params.Args["id"])
	// 			return stock, err
	// 		})),
}

// var FieldListCategories = &graphql.Field{
// 	Type:        graphql.NewList(graphql.String),
// 	Description: "List categories",

// 	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
// 		categories, err := repo.Article.GetCategories()
// 		return categories, err
// 	}),
// }

// &graphql.Field{
// 	Type: graphql.NewList(TypeArticleStock),
// 	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, model mongodb.Repo) (interface{}, error) {
// if id, ok := utils.GetValueByJSONTag(params.Source, "id"); ok {
// 	stock, err := model.Stock.FindByArticle(id)
// 	return stock, err
// }
// return nil, nil
// 	}),
// },
