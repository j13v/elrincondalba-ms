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
		orders, _, err := model.Order.FindSlice(params.Args)
		return utils.ConnectionFromArraySlice(orders, connArgs, relay.ArraySliceMetaInfo{
			SliceStart:  0,
			ArrayLength: 100,
		}), err
	}),
}

// import (
// 	"github.com/globalsign/mgo/bson"
// 	"github.com/graphql-go/graphql"
// 	"github.com/graphql-go/relay"
// 	"github.com/jal88/elrincondalba-ms/graphql/types"
// )
//
// var ListOrders = &graphql.Field{
// 	Type:        types.TypeOrderConnection,
// 	Description: "Get orders list",
// 	Args:        relay.ConnectionArgs,
// 	Resolve: decs.ContextModelConsumer(func(params graphql.ResolveParams, model inte.Model) (interface{}, error) {
// 		args := relay.NewConnectionArguments(params.Args)
// 		orders, err := model.Order().FindWithPagination(params.Args, bson.M{})
// 		if err != nil {
// 			return nil, err
// 		}
// 		return relay.ConnectionFromArraySlice(orders, args, relay.ArraySliceMetaInfo{
// 			SliceStart:  0,
// 			ArrayLength: 100,
// 		}), err
// 	}),
// }
