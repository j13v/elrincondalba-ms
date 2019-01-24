package queries

import (
	"github.com/jal88/elrincondalba-ms/types"

	"github.com/graphql-go/graphql"
)

var Order = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			/* Get (read) single order by id
			   http://localhost:8080/order?query={order(id:1){name,info,price}}
			*/
			"order": &graphql.Field{
				Type:        types.Order,
				Description: "Get order by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(ord graphql.ResolveParams) (interface{}, error) {
					id, ok := ord.Args["id"].(int)
					if ok {
						// Find order
						for _, order := range types.OrdersMock {
							if int(order.ID) == id {
								return order, nil
							}
						}
					}
					return nil, nil
				},
			},
			/* Get (read) order list
			   http://localhost:8080/order?query={list{id,name,info,price}}
			*/
			"list": &graphql.Field{
				Type:        graphql.NewList(types.Order),
				Description: "Get order list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return types.OrdersMock, nil
				},
			},
		},
	})
