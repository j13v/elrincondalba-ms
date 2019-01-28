package types

import (
	"github.com/graphql-go/graphql"
	decs "github.com/jal88/elrincondalba-ms/graphql/decorators"
	"github.com/jal88/elrincondalba-ms/graphql/utils"
	"github.com/jal88/elrincondalba-ms/mongodb"
)

/*
FieldUser GraphQL field
*/
var FieldUser = graphql.Field{
	Type:        TypeUser,
	Description: "User",
	Resolve: decs.ContextModelConsumer(func(params graphql.ResolveParams, model mongodb.Model) (interface{}, error) {
		if id, ok := utils.GetValueByJSONTag(params.Source, "user"); ok {
			user, err := model.User.FindOne(map[string]interface{}{"id": id})
			return user, err

		}
		return nil, nil
	}),
}

/*
TypeUser
*/
var TypeUser = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &FieldID,
			"dni": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"surname": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"phone": &graphql.Field{
				Type: graphql.String,
			},
			"address": &graphql.Field{
				Type: graphql.String,
			},
			"notes": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
