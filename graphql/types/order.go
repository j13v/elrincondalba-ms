package types

import (
	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincondalba-ms/graphql/enums"
	"github.com/jal88/elrincondalba-ms/graphql/utils"
)

/*
TypeOrder
*/
var TypeOrder = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Order",
		Description: "A ship in the Star Wars saga",
		Fields: graphql.Fields{
			"id":    &FieldID,
			"stock": &FieldStock,
			"user":  &FieldUser,
			"state": &graphql.Field{
				Type: enums.EnumOrderStatus,
			},
			"notes": &graphql.Field{
				Type: graphql.String,
			},
			"createdAt": &graphql.Field{
				Type: graphql.Int,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

/*
TypeOrderConnection
*/
var TypeOrderConnection = utils.ConnectionDefinitions(utils.ConnectionConfig{
	Name:     "Order",
	NodeType: TypeOrder,
})

// var OrderEdge = graphql.NewObject(
// 	graphql.ObjectConfig{
// 		Name: "OrderEdge",
// 		Fields: graphql.Fields{
// 			"cursor": &graphql.Field{
// 				Type: graphql.String,
// 			},
// 			"node": &graphql.Field{
// 				Type: Order,
// 			},
// 		},
// 	},
// )
//
// var OrderConnection = graphql.NewObject(
// 	graphql.ObjectConfig{
// 		Name: "OrderConnection",
// 		Fields: graphql.Fields{
// 			"edges": &graphql.Field{
// 				Type: graphql.NewList(ArticleEdge),
// 			},
// 		},
// 	},
// )
