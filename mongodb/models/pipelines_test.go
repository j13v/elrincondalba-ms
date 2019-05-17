package models

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mongodb/mongo-go-driver/bson"
)

func TestPipelineHelpers(t *testing.T) {
	expected := bson.A{bson.M{"testing": true}}
	actual := assertPipeline(true, expected)
	if reflect.DeepEqual(actual, expected) {
		fmt.Printf("[v] GetObjectIDFromValue (%s)\n", "assertPipeline false assertion")
	} else {
		t.Errorf("TestPipelineHelpers (%s): expected `%v`, actual `%v`", "assertPipeline false assertion", expected, actual)
	}
}

// func TestPipelinesFragments(t *testing.T) {

// 	specs := []struct {
// 		name     string
// 		input    interface{}
// 		expected *primitive.ObjectID
// 		error    error
// 	}{{
// 		"using ObjectID",
// 		pipelineArticleStockUnwinder,
// 		&oidInput,
// 		nil,
// 	}}

// 	for _, spec := range specs {
// 		actual, err := GetObjectIDFromValue(spec.input)
// 		if err != nil {
// 			if spec.error == nil {
// 				t.Errorf("GetObjectIDFromValue (%s): expected error `%v`, actual `%v`", spec.name, spec.error, err)
// 			}
// 		} else if actual != nil {
// 			if spec.expected != nil {
// 				if actual.Hex() != spec.expected.Hex() {
// 					t.Errorf("GetObjectIDFromValue (%s): expected `%v`, actual `%v`", spec.name, spec.expected, actual)
// 				}
// 			}
// 		}
// 		fmt.Printf("[v] GetObjectIDFromValue (%s)\n", spec.name)
// 	}
// }
