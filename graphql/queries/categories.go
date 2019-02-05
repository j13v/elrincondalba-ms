package queries

import (
	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincondalba-ms/graphql/types"
)

var QueryCategories = graphql.Fields{
	"listCategories": types.FieldListCategories,
}
