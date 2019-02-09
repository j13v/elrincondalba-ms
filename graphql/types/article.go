package types

import (
	"github.com/graphql-go/graphql"
	decs "github.com/j13v/elrincondalba-ms/graphql/decorators"
	"github.com/j13v/elrincondalba-ms/graphql/utils"
	"github.com/j13v/elrincondalba-ms/mongodb"
)

/*
FieldArticle GraphQL field
*/
var FieldArticle = graphql.Field{
	Type:        TypeArticle,
	Description: "Article",
	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, model mongodb.Repo) (interface{}, error) {
		if id, ok := utils.GetValueByJSONTag(params.Source, "article"); ok {
			article, err := model.Article.FindOne(utils.NewIdArgs(id))
			return article, err

		}
		return nil, nil
	}),
}

/*
TypeArticle
*/
var TypeArticle = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Article",
		Fields: graphql.Fields{
			"id": &FieldID,
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"price": &graphql.Field{
				Type: graphql.Float,
			},
			"images": &graphql.Field{
				Type: graphql.NewList(ObjectID),
			},
			"category": &graphql.Field{
				Type: graphql.String,
			},
			"rating": &graphql.Field{
				Type: graphql.Int,
			},
			"createdAt": &graphql.Field{
				Type: graphql.Int,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.Int,
			},
			"stock": &graphql.Field{
				Type: graphql.NewList(TypeArticleStock),
				Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, model mongodb.Repo) (interface{}, error) {
					if id, ok := utils.GetValueByJSONTag(params.Source, "id"); ok {
						stock, err := model.Stock.FindByArticle(id)
						return stock, err
					}
					return nil, nil
				}),
			},
		},
	},
)

var TypeArticleConnection = utils.ConnectionDefinitions(utils.ConnectionConfig{
	Name:     "Article",
	NodeType: TypeArticle,
})
