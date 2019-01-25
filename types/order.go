package types

import (
	"github.com/jal88/elrincondalba-ms/models"
	"github.com/graphql-go/graphql"
)

type OrderMock struct {
	ID       string  `json:"id"`
	Article  string `json:"article"`
	User     string `json:"user"`
	Size     string `json:"size"`
	CreateAt int32 `json:"createAt"`
	UpdateAt int32 `json:"updateAt"`
	State    int8 `json:"state"`
}

var OrdersMock []OrderMock

var OrderStatusEnum = graphql.NewEnum(graphql.EnumConfig{
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


/*
Article
*/
var Order = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Order",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"article": &graphql.Field{
				Type: Article,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					if order, ok := params.Source.(models.TypeOrder); ok {
						dbModel := params.Context.Value("dbModel").(models.DbModel)
						article, err := dbModel.Article.GetArticleById(order.Article)
						return article, err
					}
					return nil, nil
				},
			},
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
				Type: OrderStatusEnum,
			},
		},
	},
)
