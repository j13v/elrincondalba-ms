package types

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

var Cursor = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Cursor",
	Description: "",
	Serialize: func(value interface{}) interface{} {
		if oid, ok := value.(primitive.ObjectID); ok {
			return oid
		}
		return nil
	},
	ParseValue: func(value interface{}) interface{} {
		if v, ok := value.(primitive.ObjectID); ok {
			return v.Hex()
		}
		return fmt.Sprintf("%v", value)
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return valueAST.Value
		}
		return nil
	},
})

/*
Cursor GraphQL field definition used in pagination in order to establish a point reference
*/
var FieldCursor = &graphql.Field{
	Type:        graphql.ID,
	Description: "Cursor for use in pagination",
}
