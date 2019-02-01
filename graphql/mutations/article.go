package mutations

// import (
// 	"math/rand"
// 	"time"

// 	"github.com/graphql-go/graphql"
// 	"github.com/jal88/elrincondalba-ms/graphql/types"
// )

// var Article = graphql.Fields{
// 	"createArticle": &graphql.Field{
// 		Type:        types.TypeArticle,
// 		Description: "Create new article",
// 		Args: graphql.FieldConfigArgument{
// 			"name": &graphql.ArgumentConfig{
// 				Type: graphql.NewNonNull(graphql.String),
// 			},
// 			"description": &graphql.ArgumentConfig{
// 				Type: graphql.String,
// 			},
// 			"price": &graphql.ArgumentConfig{
// 				Type: graphql.NewNonNull(graphql.Float),
// 			},
// 			"images": &graphql.ArgumentConfig{
// 				Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
// 			},
// 			"category": &graphql.ArgumentConfig{
// 				Type: graphql.NewNonNull(graphql.String),
// 			},
// 			"rating": &graphql.ArgumentConfig{
// 				Type: graphql.NewNonNull(graphql.Int),
// 			},
// 		},
// 		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
// 			rand.Seed(time.Now().UnixNano())
// 			article := types.FieldArticle{
// 				ID:          string(rand.Intn(100000)), // generate random ID
// 				Name:        params.Args["name"].(string),
// 				Description: params.Args["description"].(string),
// 				Price:       params.Args["price"].(float64),
// 				Images:      params.Args["images"].([]string),
// 				Category:    params.Args["category"].(string),
// 				Rating:      params.Args["rating"].(int8),
// 			}
// 			//qltype.ArticlesMock = append(qltype.ArticlesMock, article)
// 			return article, nil
// 		},
// 	},

// 	/* Update product by id
// 	http://localhost:8080/product?query=mutation+_{update(id:1,price:3.95){id,name,info,price}}
// 	*/
// 	"updateArticle": &graphql.Field{
// 		Type:        types.TypeArticle,
// 		Description: "Update product by id",
// 		Args: graphql.FieldConfigArgument{
// 			"id": &graphql.ArgumentConfig{
// 				Type: graphql.NewNonNull(graphql.String),
// 			},
// 			"name": &graphql.ArgumentConfig{
// 				Type: graphql.NewNonNull(graphql.String),
// 			},
// 			"description": &graphql.ArgumentConfig{
// 				Type: graphql.String,
// 			},
// 			"price": &graphql.ArgumentConfig{
// 				Type: graphql.NewNonNull(graphql.Float),
// 			},
// 			"images": &graphql.ArgumentConfig{
// 				Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
// 			},
// 			"category": &graphql.ArgumentConfig{
// 				Type: graphql.NewNonNull(graphql.String),
// 			},
// 			"rating": &graphql.ArgumentConfig{
// 				Type: graphql.NewNonNull(graphql.Int),
// 			},
// 		},
// 		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
// 			id, _ := params.Args["id"].(string)
// 			name, nameOk := params.Args["name"].(string)
// 			description, descriptionOK := params.Args["description"].(string)
// 			price, priceOk := params.Args["price"].(float64)
// 			images, imagesOK := params.Args["images"].([]string)
// 			category, categoryOK := params.Args["category"].(string)
// 			rating, ratingOK := params.Args["rating"].(int8)
// 			article := qltype.ArticleMock{}
// 			for i, art := range qltype.ArticlesMock {
// 				if string(id) == art.ID {
// 					if nameOk {
// 						qltype.ArticlesMock[i].Name = name
// 					}
// 					if descriptionOK {
// 						qltype.ArticlesMock[i].Description = description
// 					}
// 					if priceOk {
// 						qltype.ArticlesMock[i].Price = price
// 					}
// 					if imagesOK {
// 						qltype.ArticlesMock[i].Images = images
// 					}
// 					if categoryOK {
// 						qltype.ArticlesMock[i].Category = category
// 					}
// 					if ratingOK {
// 						qltype.ArticlesMock[i].Rating = rating
// 					}
// 					article = qltype.ArticlesMock[i]
// 					break
// 				}
// 			}
// 			return article, nil
// 		},
// 	},

// 	/* Delete product by id
// 	http://localhost:8080/product?query=mutation+_{delete(id:1){id,name,info,price}}
// 	*/
// 	"deleteArticle": &graphql.Field{
// 		Type:        types.TypeArticle,
// 		Description: "Delete article by id",
// 		Args: graphql.FieldConfigArgument{
// 			"id": &graphql.ArgumentConfig{
// 				Type: graphql.NewNonNull(graphql.String),
// 			},
// 		},
// 		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
// 			id, _ := params.Args["id"].(string)
// 			article := types.FieldArticle{}
// 			for i, art := range types.FieldListArticles {
// 				if string(id) == art.ID {
// 					article = types.FieldListArticles[i]
// 					// Remove from product list
// 					qltype.ArticlesMock = append(qltype.ArticlesMock[:i], qltype.ArticlesMock[i+1:]...)
// 				}
// 			}

// 			return article, nil
// 		},
// 	},
// }
