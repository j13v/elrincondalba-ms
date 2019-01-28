package queries

import (
	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincondalba-ms/graphql/types"
)

var QueryOrder = graphql.Fields{
	"getOrder":   types.FieldGetOrder,
	"listOrders": types.FieldListOrders,
}

// "listOrders": &graphql.Field{
// 	Type:        qltype.OrderConnection,
// 	Args:        relay.ConnectionArgs,
// 	Description: "Get order list",
// 	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
// 		args := relay.NewConnectionArguments(params.Args)
// 		dbModel := params.Context.Value("model").(model.Models)
// 		orders, err := dbModel.Order.ListOrders(params.Args)
// 		return relay.ConnectionFromArraySlice(orders, args, relay.ArraySliceMetaInfo{
// 			SliceStart:  0,
// 			ArrayLength: 100,
// 		}), err
// 	},
// },
// func ConnectionFromData(data []interface{}) *relay.Connection{
// 	edges := []*relay.Edge{}
// 	for index, value := range data {
// 		edges = append(edges, &relay.Edge{
// 			Cursor: "value.ID",
// 			Node:   value,
// 		})
// 	}
//
// 	conn := relay.NewConnection()
// 	conn.Edges = edges
// 	conn.PageInfo = relay.PageInfo{
// 		StartCursor:     firstEdgeCursor,
// 		EndCursor:       lastEdgeCursor,
// 		HasPreviousPage: hasPreviousPage,
// 		HasNextPage:     hasNextPage,
// 	}
//
// 	return conn
// }
