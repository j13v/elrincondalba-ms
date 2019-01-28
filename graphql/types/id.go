package types

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

func coerceString(value interface{}) interface{} {
	if v, ok := value.(primitive.ObjectID); ok {
		return v.Hex()
	}
	return fmt.Sprintf("%v", value)
}

var ObjectID = graphql.NewScalar(graphql.ScalarConfig{
	Name: "ObjectID",
	Description: "The `ID` scalar type represents a unique identifier, often used to " +
		"refetch an object or as key for a cache. The ID type appears in a JSON " +
		"response as a String; however, it is not intended to be human-readable. " +
		"a 4-byte value representing the seconds since the Unix epoch" +
		"a 5-byte random value, and" +
		"a 3-byte counter, starting with a random value.",
	Serialize:  coerceString,
	ParseValue: coerceString,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return valueAST.Value
		}
		return nil
	},
})

/*
ID GraphQL field definition used to identify a resource
*/
var FieldID = graphql.Field{
	Type: ObjectID,
}
