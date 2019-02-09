package queries

import (
	"github.com/graphql-go/graphql"
	"github.com/j13v/elrincondalba-ms/graphql/types"
)

var QueryCategories = graphql.Fields{
	"listCategories": types.FieldListCategories,
}
