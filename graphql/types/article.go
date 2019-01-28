package types

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
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
			article, err := model.Article.FindById(id)
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

var TypeArticleConnection = relay.ConnectionDefinitions(relay.ConnectionConfig{
	Name:     "Article",
	NodeType: TypeArticle,
}).ConnectionType

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
