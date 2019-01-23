package mutations

import (
	"math/rand"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincodalba-ms/types"
)

var Article = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		/* Create new product item
		http://localhost:8080/product?query=mutation+_{create(name:"Inca Kola",info:"Inca Kola is a soft drink that was created in Peru in 1935 by British immigrant Joseph Robinson Lindley using lemon verbena (wiki)",price:1.99){id,name,info,price}}
		*/
		"create": &graphql.Field{
			Type:        types.Article,
			Description: "Create new article",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"price": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Float),
				},
				"images": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"category": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"rating": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				rand.Seed(time.Now().UnixNano())
				article := types.ArticleMock{
					ID:          int64(rand.Intn(100000)), // generate random ID
					Name:        params.Args["name"].(string),
					Description: params.Args["description"].(string),
					Price:       params.Args["price"].(float64),
					Images:      params.Args["images"].(string),
					Category:    params.Args["category"].(string),
					Rating:      params.Args["rating"].(int8),
				}
				types.ArticlesMock = append(types.ArticlesMock, article)
				return article, nil
			},
		},

		/* Update product by id
		   http://localhost:8080/product?query=mutation+_{update(id:1,price:3.95){id,name,info,price}}
		*/
		"update": &graphql.Field{
			Type:        types.Article,
			Description: "Update product by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"price": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Float),
				},
				"images": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"category": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"rating": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(int)
				name, nameOk := params.Args["name"].(string)
				description, descriptionOK := params.Args["description"].(string)
				price, priceOk := params.Args["price"].(float64)
				images, imagesOK := params.Args["images"].(string)
				category, categoryOK := params.Args["category"].(string)
				rating, ratingOK := params.Args["rating"].(int8)
				article := types.ArticleMock{}
				for i, art := range types.ArticlesMock {
					if int64(id) == art.ID {
						if nameOk {
							types.ArticlesMock[i].Name = name
						}
						if descriptionOK {
							types.ArticlesMock[i].Description = description
						}
						if priceOk {
							types.ArticlesMock[i].Price = price
						}
						if imagesOK {
							types.ArticlesMock[i].Images = images
						}
						if categoryOK {
							types.ArticlesMock[i].Category = category
						}
						if ratingOK {
							types.ArticlesMock[i].Rating = rating
						}
						article = types.ArticlesMock[i]
						break
					}
				}
				return article, nil
			},
		},

		/* Delete product by id
		   http://localhost:8080/product?query=mutation+_{delete(id:1){id,name,info,price}}
		*/
		"delete": &graphql.Field{
			Type:        types.Article,
			Description: "Delete article by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(int)
				article := types.ArticleMock{}
				for i, art := range types.ArticlesMock {
					if int64(id) == art.ID {
						article = types.ArticlesMock[i]
						// Remove from product list
						types.ArticlesMock = append(types.ArticlesMock[:i], types.ArticlesMock[i+1:]...)
					}
				}

				return article, nil
			},
		},
	},
})
