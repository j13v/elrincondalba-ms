package queries

import (
	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincondalba-ms/graphql/types"
)

var QueryArticle = graphql.Fields{
	"getArticle":           types.FieldGetArticle,
	"getArticlePriceRange": types.FieldArticlePriceRange,
	"listArticles":         types.FieldListArticles,
	"listSizes":            types.FieldListSizes,
}
