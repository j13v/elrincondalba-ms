package types

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	decs "github.com/j13v/elrincondalba-ms/graphql/decorators"
	"github.com/j13v/elrincondalba-ms/graphql/utils"
	"github.com/j13v/elrincondalba-ms/mongodb"
)

/*
FieldListArticlesRelated ideally it would nice if the method does accept
multiple criterias in order to reference by article or user. But at first
instance it will provide only one way to obtain a list of articles related
to another article using its ID and filtering by category
*/
var FieldListArticlesRelated = &graphql.Field{
	Type:        TypeArticleConnection,
	Description: "Get articles list",
	Args:        relay.ConnectionArgs,
	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, model mongodb.Repo) (interface{}, error) {
		connArgs := relay.NewConnectionArguments(params.Args)
		articles, meta, err := model.Article.FindSlice(&params.Args)
		return utils.ConnectionFromArraySlice(articles, connArgs, meta), err
	}),
}
