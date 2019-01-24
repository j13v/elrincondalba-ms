package types

import (
	"github.com/graphql-go/graphql"
)

type OrderMock struct {
	ID       int64  `json:"id"`
	Article  string `json:"article"`
	User     string `json:"user"`
	Size     string `json:"size"`
	CreateAt string `json:"createAt"`
	UpdateAt string `json:"updateAt"`
	State    string `json:"state"`
}

var OrdersMock []OrderMock

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
				Type: graphql.String,
			},
			"user": &graphql.Field{
				Type: graphql.String,
			},
			"size": &graphql.Field{
				Type: graphql.String,
			},
			"createAt": &graphql.Field{
				Type: graphql.String,
			},
			"updateAt": &graphql.Field{
				Type: graphql.String,
			},
			"state": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
