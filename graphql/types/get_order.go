package types

import (
	"github.com/graphql-go/graphql"
	decs "github.com/j13v/elrincondalba-ms/graphql/decorators"
	"github.com/j13v/elrincondalba-ms/mongodb"
)

var FieldGetOrder = &graphql.Field{
	Type:        TypeOrder,
	Description: "Get order by id",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: ObjectID,
		},
	},
	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, model mongodb.Repo) (interface{}, error) {
		order, err := model.Order.FindOne(&params.Args)
		return order, err
	}),
}
