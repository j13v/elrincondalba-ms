package types

import "github.com/graphql-go/graphql"

var TypeUpload = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Upload",
	Description: "Scalar upload object",
})
