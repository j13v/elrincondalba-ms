package queries

import (
	"github.com/graphql-go/graphql"
	"github.com/j13v/elrincondalba-ms/graphql/types"
)

var QueryArticle = graphql.Fields{
	"getArticle":           types.FieldGetArticle,
	"getArticlePriceRange": types.FieldArticlePriceRange,
	"listArticles":         types.FieldListArticles,
	"listSizes":            types.FieldListSizes,
}
