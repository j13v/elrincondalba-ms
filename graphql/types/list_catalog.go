package types

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	decs "github.com/j13v/elrincondalba-ms/graphql/decorators"
	"github.com/j13v/elrincondalba-ms/graphql/utils"
	"github.com/j13v/elrincondalba-ms/mongodb"
)

var FieldListCatalog = &graphql.Field{
	Type:        TypeArticleConnection,
	Description: "Get catalog list by filters",
	Args:        relay.ConnectionArgs,
	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, model mongodb.Repo) (interface{}, error) {
		connArgs := relay.NewConnectionArguments(params.Args)
		articles, meta, err := model.Article.FindSlice(&params.Args)
		return utils.ConnectionFromArraySlice(articles, connArgs, meta), err
	}),
}
