package types

import (
	"github.com/graphql-go/graphql"
	decs "github.com/jal88/elrincondalba-ms/graphql/decorators"
	"github.com/jal88/elrincondalba-ms/mongodb"
)

var FieldGetOrder = &graphql.Field{
	Type:        TypeOrder,
	Description: "Get order by id",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: decs.ContextModelConsumer(func(params graphql.ResolveParams, model mongodb.Model) (interface{}, error) {
		if id, ok := params.Args["id"].(string); ok {
			order, err := model.Order.FindById(id)
			return order, err
		}
		return nil, nil
	}),
}
