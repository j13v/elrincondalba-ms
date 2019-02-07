package types

import (
	"github.com/graphql-go/graphql"
	decs "github.com/jal88/elrincondalba-ms/graphql/decorators"
	"github.com/jal88/elrincondalba-ms/mongodb"
)

var FieldGetArticle = &graphql.Field{
	Type:        TypeArticle,
	Description: "Get article by id",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, model mongodb.Repo) (interface{}, error) {
		article, err := model.Article.FindOne(&params.Args)
		return article, err
	}),
}
