package types

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	decs "github.com/jal88/elrincondalba-ms/graphql/decorators"
	"github.com/jal88/elrincondalba-ms/graphql/utils"
	"github.com/jal88/elrincondalba-ms/mongodb"
)

/*
FieldListArticles
{
  listArticles(last:3){
    pageInfo {
      endCursor
      hasNextPage
      hasPreviousPage
      startCursor
    }
    edges {
      cursor
      node {
        category
        description
        id
        images
        name
        price
        rating
      }
    }
  }
}
*/
var FieldListArticles = &graphql.Field{
	Type:        TypeArticleConnection,
	Description: "Get articles list",
	Args:        relay.ConnectionArgs,
	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, model mongodb.Repo) (interface{}, error) {
		connArgs := relay.NewConnectionArguments(params.Args)
		articles, meta, err := model.Article.FindSlice(&params.Args)
		return utils.ConnectionFromArraySlice(articles, connArgs, meta), err
	}),
}
