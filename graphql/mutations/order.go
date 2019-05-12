package mutations

import (
	"fmt"

	"github.com/graphql-go/graphql"
	decs "github.com/j13v/elrincondalba-ms/graphql/decorators"
	"github.com/j13v/elrincondalba-ms/graphql/types"
	"github.com/j13v/elrincondalba-ms/mongodb"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

/*
Mutation Create Order
*/
var MutationOrder = graphql.Fields{
	"createOrder": &graphql.Field{
		Type:        types.TypeOrder,
		Description: "Create new order",
		Args: graphql.FieldConfigArgument{
			"stockSize": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"name": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"surname": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"email": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"phone": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"address": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"notes": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
			notes := ""
			if params.Args["notes"] != nil {
				notes = params.Args["notes"].(string)
			}
			user, err := repo.User.Create(
				params.Args["name"].(string),
				params.Args["surname"].(string),
				params.Args["email"].(string),
				params.Args["phone"].(string),
				params.Args["address"].(string),
			)
			if err != nil {
				return nil, err
			}
			stockSize := params.Args["stockSize"].(string)
			stockArticle, err := repo.Article.FindStockBySize(stockSize)
			if err != nil {
				return nil, fmt.Errorf("No stock found by size %s", stockSize)
			}
			order, err := repo.Order.Create(
				stockArticle.ID,
				user.ID,
				notes,
			)
			return order, err
		}),
	},
	/*
		Purchase order
	*/
	"purchaseOrder": &graphql.Field{
		Type:        types.TypeOrder,
		Description: "Purchase order",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(types.ObjectID),
			},
			"paymentMethod": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"paymentRef": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
			paymentMethod := int8(params.Args["paymentMethod"].(int))
			order, err := repo.Order.Purchase(
				params.Args["id"].(primitive.ObjectID),
				paymentMethod,
				params.Args["paymentRef"].(string))
			return order, err
		}),
	},
	"prepareOrder": &graphql.Field{
		Type:        types.TypeOrder,
		Description: "Purchase order",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(types.ObjectID),
			},
		},
		Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
			order, err := repo.Order.Prepare(params.Args["id"].(primitive.ObjectID))
			return order, err
		}),
	},
	"shipOrder": &graphql.Field{
		Type:        types.TypeOrder,
		Description: "Purchase order",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(types.ObjectID),
			},
			"trackingRef": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
			order, err := repo.Order.Ship(params.Args["id"].(primitive.ObjectID), params.Args["trackingRef"].(string))
			return order, err
		}),
	},
	"confirmReceived": &graphql.Field{
		Type:        types.TypeOrder,
		Description: "Purchase order",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(types.ObjectID),
			},
		},
		Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {

			order, err := repo.Order.ConfirmReceived(params.Args["id"].(primitive.ObjectID))
			return order, err
		}),
	},
	"cancelOrder": &graphql.Field{
		Type:        types.TypeOrder,
		Description: "Purchase order",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(types.ObjectID),
			},
		},
		Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
			order, err := repo.Order.Cancel(params.Args["id"].(primitive.ObjectID))
			return order, err
		}),
	},
}
