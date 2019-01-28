package types

import (
	"github.com/graphql-go/graphql"
)

/*
Cursor GraphQL field definition used in pagination in order to establish a point reference
*/
var FieldCursor = &graphql.Field{
	Type:        graphql.ID,
	Description: "Cursor for use in pagination",
}
