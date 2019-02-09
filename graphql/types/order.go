package types

import (
	"github.com/graphql-go/graphql"
	"github.com/j13v/elrincondalba-ms/graphql/utils"
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
				Type: graphql.Int,
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
