package types

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	decs "github.com/jal88/elrincondalba-ms/graphql/decorators"
	"github.com/jal88/elrincondalba-ms/graphql/utils"
	"github.com/jal88/elrincondalba-ms/mongodb"
)

var FieldListOrders = &graphql.Field{
	Type:        TypeOrderConnection,
	Description: "Get orders list",
	Args:        relay.ConnectionArgs,
	Resolve: decs.ContextModelConsumer(func(params graphql.ResolveParams, model mongodb.Model) (interface{}, error) {
		connArgs := relay.NewConnectionArguments(params.Args)
		orders, meta, err := model.Order.FindSlice(params.Args)
		return utils.ConnectionFromArraySlice(orders, connArgs, meta), err
	}),
}
