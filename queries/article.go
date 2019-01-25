package queries

import (
	"github.com/jal88/elrincondalba-ms/types"
	"github.com/jal88/elrincondalba-ms/models"
	"github.com/graphql-go/graphql"
)

var Article = graphql.Fields{
	"getArticle": &graphql.Field{
		Type:        types.Article,
		Description: "Get article by id",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			id, idOk := params.Args["id"].(string)
			dbModel := params.Context.Value("dbModel").(models.DbModel)
			if idOk {
				article, err := dbModel.Article.GetArticleById(id)
				return article, err
			}
			return nil, nil
		},
	},
	"listArticles": &graphql.Field{
		Type:        graphql.NewList(types.Article),
		Description: "Get product list",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			dbModel := params.Context.Value("dbModel").(models.DbModel)
			articles, err := dbModel.Article.ListArticles()
			return articles, err
		},
	},
}
