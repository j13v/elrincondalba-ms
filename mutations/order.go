package mutations

import (
	"math/rand"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincondalba-ms/types"
)

var Order = graphql.NewObject(graphql.ObjectConfig{
	Name: "MutationOrder",
	Fields: graphql.Fields{
		/* Create new order item
		http://localhost:8080/order?query=mutation+_{create(name:"Inca Kola",info:"Inca Kola is a soft drink that was created in Peru in 1935 by British immigrant Joseph Robinson Lindley using lemon verbena (wiki)",price:1.99){id,name,info,price}}
		*/
		"create": &graphql.Field{
			Type:        types.Order,
			Description: "Create new order",
			Args: graphql.FieldConfigArgument{
				"article": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"user": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"size": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"createAt": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"updateAt": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"state": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				rand.Seed(time.Now().UnixNano())
				order := types.OrderMock{
					ID:       string(rand.Intn(100000)), // generate random ID
					Article:  params.Args["article"].(string),
					User:     params.Args["description"].(string),
					Size:     params.Args["size"].(string),
					CreateAt: params.Args["createAt"].(int32),
					UpdateAt: params.Args["updateAt"].(int32),
					State:    params.Args["state"].(int8),
				}
				types.OrdersMock = append(types.OrdersMock, order)
				return order, nil
			},
		},

		/* Update product by id
		   http://localhost:8080/product?query=mutation+_{update(id:1,price:3.95){id,name,info,price}}
		*/
		"update": &graphql.Field{
			Type:        types.Article,
			Description: "Update order by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"article": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"user": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"size": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"createAt": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"updateAt": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"state": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(string)
				article, articleOK := params.Args["article"].(string)
				user, userOK := params.Args["user"].(string)
				size, sizeOK := params.Args["size"].(string)
				createAt, createAtOK := params.Args["createAt"].(int32)
				updateAt, updateAtOK := params.Args["updateAt"].(int32)
				state, stateOK := params.Args["state"].(int8)
				order := types.OrderMock{}
				for i, ord := range types.OrdersMock {
					if string(id) == ord.ID {
						if articleOK {
							types.OrdersMock[i].Article = article
						}
						if userOK {
							types.OrdersMock[i].User = user
						}
						if sizeOK {
							types.OrdersMock[i].Size = size
						}
						if createAtOK {
							types.OrdersMock[i].CreateAt = createAt
						}
						if updateAtOK {
							types.OrdersMock[i].UpdateAt = updateAt
						}
						if stateOK {
							types.OrdersMock[i].State = state
						}
						order = types.OrdersMock[i]
						break
					}
				}
				return order, nil
			},
		},

		/* Delete product by id
		   http://localhost:8080/product?query=mutation+_{delete(id:1){id,name,info,price}}
		*/
		"delete": &graphql.Field{
			Type:        types.Order,
			Description: "Delete order by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(string)
				order := types.OrderMock{}
				for i, ord := range types.OrdersMock {
					if string(id) == ord.ID {
						order = types.OrdersMock[i]
						// Remove from product list
						types.OrdersMock = append(types.OrdersMock[:i], types.OrdersMock[i+1:]...)
					}
				}

				return order, nil
			},
		},
	},
})
