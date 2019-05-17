package types

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	decs "github.com/j13v/elrincondalba-ms/graphql/decorators"
	"github.com/j13v/elrincondalba-ms/graphql/utils"
	"github.com/j13v/elrincondalba-ms/mongodb"
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
	Args: relay.NewConnectionArgs(graphql.FieldConfigArgument{
		"categories": &graphql.ArgumentConfig{
			Type: graphql.NewList(graphql.String),
		},
		"priceRange": &graphql.ArgumentConfig{
			Type: graphql.NewList(graphql.Float),
		},
		"sizes": &graphql.ArgumentConfig{
			Type: graphql.NewList(graphql.String),
		},
		"sorting": &graphql.ArgumentConfig{
			Type: graphql.NewEnum(graphql.EnumConfig{
				Name: "ArticleSortingEnum",
				Values: graphql.EnumValueConfigMap{
					"POPULAR": &graphql.EnumValueConfig{
						Value: 1,
					},
					"PURCHASES": &graphql.EnumValueConfig{
						Value: 2,
					},
					"RECENTS": &graphql.EnumValueConfig{
						Value: 3,
					},
				},
			}),
		},
	}),
	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, model mongodb.Repo) (interface{}, error) {
		connArgs := relay.NewConnectionArguments(params.Args)
		articles, meta, err := model.Article.FindSlice(&params.Args)
		return utils.ConnectionFromArraySlice(articles, connArgs, meta), err
	}),
}
