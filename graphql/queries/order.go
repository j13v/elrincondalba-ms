package queries

import (
	"github.com/graphql-go/graphql"
	"github.com/j13v/elrincondalba-ms/graphql/types"
)

var QueryOrder = graphql.Fields{
	"getOrder":   types.FieldGetOrder,
	"listOrders": types.FieldListOrders,
}
