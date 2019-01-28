package types

import (
	"github.com/graphql-go/graphql"
	decs "github.com/jal88/elrincondalba-ms/graphql/decorators"
	"github.com/jal88/elrincondalba-ms/graphql/utils"
	"github.com/jal88/elrincondalba-ms/mongodb"
)

/*
FieldArticle GraphQL field
*/
var FieldArticle = graphql.Field{
	Type:        TypeArticle,
	Description: "Article",
	Resolve: decs.ContextModelConsumer(func(params graphql.ResolveParams, model mongodb.Model) (interface{}, error) {
		if id, ok := utils.GetValueByJSONTag(params.Source, "article"); ok {
			if article, err := model.Article.FindOne(map[string]interface{}{"id": id}); err != nil {
				return article, err
			}
		}
		return nil, nil
	}),
}

/*
TypeArticle
*/
var TypeArticle = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "article",
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
				Type: graphql.NewList(graphql.String),
			},
			"category": &graphql.Field{
				Type: graphql.String,
			},
			"rating": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

var TypeArticleConnection = utils.ConnectionDefinitions(utils.ConnectionConfig{
	Name:     "article",
	NodeType: TypeArticle,
})

// var ArticleEdge = graphql.NewObject(
// 	graphql.ObjectConfig{
// 		Name: "ArticleEdge",
// 		Fields: graphql.Fields{
// 			"cursor": &graphql.Field{
// 				Type: graphql.String,
// 			},
// 			"node": &graphql.Field{
// 				Type: Article,
// 			},
// 		},
// 	},
// )