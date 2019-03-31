package types

import (
	"github.com/graphql-go/graphql"
)

var TypeDistinct = graphql.ObjectConfig{
	Name: "Pepesito",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"count": &graphql.Field{
			Type: graphql.Int,
		},
	},
}

var TypeDistinctList = graphql.NewList(graphql.NewObject(TypeDistinct))
