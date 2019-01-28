package types

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"github.com/jal88/elrincondalba-ms/graphql/enums"
)

/*
TypeOrder
*/
var TypeOrder = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Order",
		Description: "A ship in the Star Wars saga",
		Fields: graphql.Fields{
			"id":      &FieldID,
			"article": &FieldArticle,
			"user": &graphql.Field{
				Type: graphql.String,
			},
			"size": &graphql.Field{
				Type: graphql.String,
			},
			"createAt": &graphql.Field{
				Type: graphql.Int,
			},
			"updateAt": &graphql.Field{
				Type: graphql.Int,
			},
			"state": &graphql.Field{
				Type: enums.EnumOrderStatus,
			},
		},
	},
)

/*
TypeOrderConnection
*/
var TypeOrderConnection = relay.ConnectionDefinitions(relay.ConnectionConfig{
	Name:     "Order",
	NodeType: TypeOrder,
}).ConnectionType

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
