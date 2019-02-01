package enums

import "github.com/graphql-go/graphql"

var EnumOrderStatus = graphql.NewEnum(graphql.EnumConfig{
	Name:        "Status",
	Description: "One of the states in an order",
	Values: graphql.EnumValueConfigMap{
		"PENDING": &graphql.EnumValueConfig{
			Value:       0,
			Description: "A user request an article.",
		},
		"PURCHASED": &graphql.EnumValueConfig{
			Value:       1,
			Description: "User pay the article.",
		},
		"PREPARING": &graphql.EnumValueConfig{
			Value:       2,
			Description: "Admin prepares the article to be shipped.",
		},
		"SHIPPING": &graphql.EnumValueConfig{
			Value:       3,
			Description: "Admin/User cancel the order flow.",
		},
		"RECEIVED": &graphql.EnumValueConfig{
			Value:       4,
			Description: "User receive the order.",
		},
		"CANCELLED": &graphql.EnumValueConfig{
			Value:       5,
			Description: "Admin/User cancel the order flow.",
		},
	},
})
