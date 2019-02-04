package queries

import (
	"github.com/graphql-go/graphql"
	decs "github.com/jal88/elrincondalba-ms/graphql/decorators"
	"github.com/jal88/elrincondalba-ms/graphql/types"
	"github.com/jal88/elrincondalba-ms/mongodb"
)

var TypeMinMax = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "MinMaxPrice",
		Fields: graphql.Fields{
			"min": &graphql.Field{
				Type: graphql.Float,
			},
			"max": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)
var FieldArticleMinMaxPrice = &graphql.Field{
	Type:        TypeMinMax,
	Description: "Get min max price",

	Resolve: decs.ContextRepoConsumer(func(params graphql.ResolveParams, repo mongodb.Repo) (interface{}, error) {
		rangePrice, err := repo.Article.GetMinMaxPrice()
		return rangePrice, err
	}),
}
var QueryArticle = graphql.Fields{
	"getArticle":            types.FieldGetArticle,
	"listArticles":          types.FieldListArticles,
	"getArticleMinMaxPrice": FieldArticleMinMaxPrice,
}
