package queries

import (
	"github.com/graphql-go/graphql"
	"github.com/j13v/elrincondalba-ms/graphql/types"
)

var QueryArticle = graphql.Fields{
	"getArticle":            types.FieldGetArticle,
	"getArticlePriceRange":  types.FieldArticlePriceRange,
	"listArticles":          types.FieldListArticles,
	"listArticleStock":      types.FieldArticleStock,
	"listArticlesRelated":   types.FieldListArticlesRelated,
	"listArticleSizes":      types.FieldListArticleSizes,
	"listArticleCategories": types.FieldListArticleCategories,
}
