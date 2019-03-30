package types

import (
	"github.com/graphql-go/graphql"
	decs "github.com/j13v/elrincondalba-ms/graphql/decorators"
	"github.com/j13v/elrincondalba-ms/graphql/utils"
	"github.com/j13v/elrincondalba-ms/mongodb"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

/*
FieldArticle GraphQL field
*/
var FieldStock = graphql.Field{
	Type:        TypeStock,
	Description: "Stock",
	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, model mongodb.Repo) (interface{}, error) {
		if id, ok := utils.GetValueByJSONTag(params.Source, "stock"); ok {
			stock, err := model.Article.FindStockById(id.(primitive.ObjectID))
			return stock, err

		}
		return nil, nil
	}),
}

/*
TypeStock
*/
var TypeStock = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Stock",
		Fields: graphql.Fields{
			"id":      &FieldID,
			"article": &FieldArticle,
			"size": &graphql.Field{
				Type: graphql.String,
			},
			"createdAt": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
