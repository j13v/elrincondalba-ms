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
var FieldStock = graphql.Field{
	Type:        TypeStock,
	Description: "Stock",
	Resolve: decs.ContextModelConsumer(func(params graphql.ResolveParams, model mongodb.Model) (interface{}, error) {
		if id, ok := utils.GetValueByJSONTag(params.Source, "stock"); ok {
			stock, err := model.Stock.FindOne(map[string]interface{}{"id": id})
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
