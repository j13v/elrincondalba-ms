package queries

import (
	"github.com/jal88/elrincodalba-ms/types"

	"github.com/graphql-go/graphql"
)

var Article = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			/* Get (read) single product by id
			   http://localhost:8080/product?query={product(id:1){name,info,price}}
			*/
			"article": &graphql.Field{
				Type:        types.Article,
				Description: "Get article by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if ok {
						// Find article
						for _, article := range types.ArticlesMock {
							if int(article.ID) == id {
								return article, nil
							}
						}
					}
					return nil, nil
				},
			},
			/* Get (read) product list
			   http://localhost:8080/product?query={list{id,name,info,price}}
			*/
			"list": &graphql.Field{
				Type:        graphql.NewList(types.Article),
				Description: "Get product list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return types.ArticlesMock, nil
				},
			},
		},
	})
