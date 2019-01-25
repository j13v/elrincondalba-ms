package queries

import (
	"github.com/jal88/elrincondalba-ms/types"
	"github.com/jal88/elrincondalba-ms/models"
	"github.com/graphql-go/graphql"
)

var Order = graphql.Fields{
	/* Get (read) single order by id
		 http://localhost:8080/order?query={order(id:1){name,info,price}}
	*/
	"getOrder": &graphql.Field{
		Type:        types.Order,
		Description: "Get order by id",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			id, idOk := params.Args["id"].(string)
			dbModel := params.Context.Value("dbModel").(models.DbModel)
			if idOk {
				order, err := dbModel.Order.GetOrderById(id)
				return order, err
			}
			return nil, nil
		},
	},
	/* Get (read) order list
		 http://localhost:8080/order?query={list{id,name,info,price}}
	*/
	"listOrders": &graphql.Field{
		Type:        graphql.NewList(types.Order),
		Description: "Get order list",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			dbModel := params.Context.Value("dbModel").(models.DbModel)
			orders, err := dbModel.Order.ListOrders()
			return orders, err
		},
	},
}
