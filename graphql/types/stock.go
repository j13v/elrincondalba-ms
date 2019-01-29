package types

import (
	"github.com/graphql-go/graphql"
)

/*
TypeStock
*/
var TypeStock = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Stock",
		Fields: graphql.Fields{
			"size": &graphql.Field{
				Type: graphql.String,
			},
			"count": &graphql.Field{
				Type: graphql.Int,
			},
			"createAt": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
