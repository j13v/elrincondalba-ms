package util

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql"
)

func CombineFields(fields ...graphql.Fields) graphql.Fields {
	resultFields := graphql.Fields{}

	for _, field := range fields {
		for nameProp, fieldDefinition := range field {
			resultFields[nameProp] = fieldDefinition
		}
	}

	return resultFields
}

func InArray(arr []string, value string) bool {
	for _, item := range arr {
		if item == value {
			return true
		}
	}
	return false
}

func Contains(source []string, matches []string, any bool) (bool, string) {
	for _, match := range matches {
		if ok := InArray(source, match); !ok && !any {
			return false, match
		}
	}
	return true, ""
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		invalidTypeError := errors.New("Provided value type didn't match obj field type")
		return invalidTypeError
	}

	structFieldValue.Set(val)
	return nil
}

func FillStruct(s interface{}, m map[string]interface{}) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// func InArray(array interface{}, val interface{}) (exists bool, index int) {
//     exists = false
//     index = -1
//
//     switch reflect.TypeOf(array).Kind() {
//     case reflect.Slice:
//         s := reflect.ValueOf(array)
//
//         for i := 0; i < s.Len(); i++ {
//             if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
//                 index = i
//                 exists = true
//                 return
//             }
//         }
//     }
//
//     return
// }
