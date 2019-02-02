package helpers

import (
	"errors"
	"fmt"
	"testing"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

func TestGetObjectIDFromValue(t *testing.T) {

	strInput := "5c5582f950c7ebd97e2cb0e6"
	oidInput, _ := primitive.ObjectIDFromHex(strInput)
	bsonInput, _ := bson.Marshal(struct {
		ID primitive.ObjectID `bson:"_id"`
	}{ID: oidInput})
	specs := []struct {
		name     string
		input    interface{}
		expected *primitive.ObjectID
		error    error
	}{{
		"using ObjectID",
		oidInput,
		&oidInput,
		nil,
	}, {
		"using pointer to ObjectID",
		&oidInput,
		&oidInput,
		nil,
	}, {
		"using string",
		strInput,
		&oidInput,
		nil,
	}, {
		"using pointer to string",
		&strInput,
		&oidInput,
		nil,
	}, {
		"using string witch wrong format",
		"as",
		nil,
		errors.New("encoding/hex: invalid byte: U+0073 's'"),
	}, {
		"using struct with a field ID of type pointer ObjectID",
		struct {
			ID *primitive.ObjectID `bson:"_id"`
		}{&oidInput},
		&oidInput,
		nil,
	}, {
		"using struct with a field ID of type ObjectID",
		struct {
			ID primitive.ObjectID `bson:"_id"`
		}{oidInput},
		&oidInput,
		nil,
	}, {
		"using struct with a field ID of type string",
		struct {
			ID string `bson:"_id"`
		}{strInput},
		&oidInput,
		nil,
	}, {
		"using struct with a field ID of type pointer to string",
		struct {
			ID *string `bson:"_id"`
		}{&strInput},
		&oidInput,
		nil,
	}, {
		"using empty struct",
		struct{}{},
		&oidInput,
		errors.New("Missing ObjectID inside the struct trying catch from {}"),
	}, {
		"using bson.Raw",
		bsonInput,
		&oidInput,
		nil,
	}, {
		"using non allow value",
		true,
		&oidInput,
		errors.New("Invalid struct value true"),
	}}

	for _, spec := range specs {
		actual, err := GetObjectIDFromValue(spec.input)
		if err != nil {
			if spec.error == nil {
				t.Errorf("GetObjectIDFromValue (%s): expected error `%v`, actual `%v`", spec.name, spec.error, err)
			}
		} else if actual != nil {
			if spec.expected != nil {
				if actual.Hex() != spec.expected.Hex() {
					t.Errorf("GetObjectIDFromValue (%s): expected `%v`, actual `%v`", spec.name, spec.expected, actual)
				}
			}
		}
		fmt.Printf("[v] GetObjectIDFromValue (%s)\n", spec.name)
	}
}
