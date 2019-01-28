package queries

import (
	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincondalba-ms/graphql/types"
)

var QueryOrder = graphql.Fields{
	"getOrder":   types.FieldGetOrder,
	"listOrders": types.FieldListOrders,
}
