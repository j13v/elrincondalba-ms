package mutations

import (
	"github.com/graphql-go/graphql"
	defs "github.com/jal88/elrincondalba-ms/definitions"
	decs "github.com/jal88/elrincondalba-ms/graphql/decorators"
	"github.com/jal88/elrincondalba-ms/graphql/types"
	"github.com/jal88/elrincondalba-ms/mongodb"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

func parseImages(input []interface{}) (out []defs.File) {
	for _, v := range input {
		// using a type assertion, convert v to a string
		out = append(out, v.(defs.File))
	}
	return out
}

var MutationArticle = graphql.Fields{
	"createArticle": &graphql.Field{
		Type:        types.TypeArticle,
		Description: "Create new article",
		Args: graphql.FieldConfigArgument{
			"name": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"description": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"price": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Float),
			},
			"images": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.NewList(types.TypeUpload)),
			},
			"category": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"rating": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
			article, err := repo.Article.Create(
				params.Args["name"].(string),
				params.Args["description"].(string),
				params.Args["price"].(float64),
				parseImages(params.Args["images"].([]interface{})),
				params.Args["category"].(string),
				int8(params.Args["rating"].(int)),
			)
			return article, err
		}),
	},

	/* Update product by id
	http://localhost:8080/product?query=mutation+_{update(id:1,price:3.95){id,name,info,price}}
	*/
	"updateArticle": &graphql.Field{
		Type:        types.TypeArticle,
		Description: "Update product by id",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"name": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"description": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"price": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Float),
			},
			"images": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
			},
			"category": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"rating": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
			// aid, err := primitive.ObjectIDFromHex(params.Args["id"].(string))
			// if err != nil {
			// 	return nil, err
			// }

			// id, _ := params.Args["id"].(string)
			// name, nameOk := params.Args["name"].(string)
			// description, descriptionOK := params.Args["description"].(string)
			// price, priceOk := params.Args["price"].(float64)
			// images, imagesOK := params.Args["images"].([]string)
			// category, categoryOK := params.Args["category"].(string)
			// rating, ratingOK := params.Args["rating"].(int8)
			// article := repo.Article.updateArticle
			// err = repo.Article.updateArticle(oid, state)
			return nil, nil

			return nil, nil
		}),
	},

	/* Delete product by id
	http://localhost:8080/product?query=mutation+_{delete(id:1){id,name,info,price}}
	*/
	"deleteArticle": &graphql.Field{
		Type:        types.TypeArticle,
		Description: "Delete article by id",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
			aid, err := primitive.ObjectIDFromHex(params.Args["id"].(string))

			err = repo.Article.DeleteArticle(aid)
			return nil, err
		}),
	},
}
