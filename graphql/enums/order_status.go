package enums

import "github.com/graphql-go/graphql"

var EnumOrderStatus = graphql.NewEnum(graphql.EnumConfig{
	Name:        "Status",
	Description: "One of the states in an order",
	Values: graphql.EnumValueConfigMap{
		"PENDING": &graphql.EnumValueConfig{
			Value:       1,
			Description: "A user request an article.",
		},
		"APPROVED": &graphql.EnumValueConfig{
			Value:       2,
			Description: "Admin conciliate the purchase and shipping.",
		},
		"COMPLETED": &graphql.EnumValueConfig{
			Value:       3,
			Description: "Admin sends article and user receives it.",
		},
	},
})
