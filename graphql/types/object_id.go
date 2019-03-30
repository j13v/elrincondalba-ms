package types

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/j13v/elrincondalba-ms/graphql/utils"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

func serializeObjectID(value interface{}) interface{} {
	if oid, ok := value.(primitive.ObjectID); ok {
		hash, err := utils.HexToBase58(oid.Hex())
		if err != nil {
			return err
		}
		return hash
	}
	return nil
}

func unserializeObjectID(value interface{}) interface{} {
	oid, err := primitive.ObjectIDFromHex(utils.Base58ToHex(value.(string)))
	if err != nil {
		return err
	}
	return oid
}

var ObjectID = graphql.NewScalar(graphql.ScalarConfig{
	Name: "ObjectID",
	Description: "The `ID` scalar type represents a unique identifier, often used to " +
		"refetch an object or as key for a cache. The ID type appears in a JSON " +
		"response as a String; however, it is not intended to be human-readable. " +
		"a 4-byte value representing the seconds since the Unix epoch" +
		"a 5-byte random value, and" +
		"a 3-byte counter, starting with a random value.",
	Serialize:  serializeObjectID,
	ParseValue: unserializeObjectID,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return unserializeObjectID(valueAST.Value)
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
