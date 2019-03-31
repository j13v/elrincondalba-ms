package helpers

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

func Base64ToHex(strB64 string) (string, error) {
	src, err := base64.StdEncoding.DecodeString(strB64)
	if err != nil {
		return "", err
	}
	dst := make([]byte, hex.EncodedLen(len(src)))
	_ = hex.Encode(dst, src)

	return fmt.Sprintf("%s", dst), nil
}

func NewFindOptionsFromArgs(args map[string]interface{}, count int64) *options.FindOptions {
	fopts := &options.FindOptions{}
	if args != nil {
		var first, last, limit, skip int64
		if firstArg, ok := args["first"]; ok {
			if firstArg, ok := firstArg.(int); ok {
				first = int64(firstArg)
			}
		}

		if lastArg, ok := args["last"]; ok {
			if lastArg, ok := lastArg.(int); ok {
				last = int64(lastArg)
			}
		}

		if first != 0 && count > first {
			limit = first
		}

		if last != 0 {
			if limit != 0 && limit > last {
				skip = limit - last
				limit = limit - skip
			} else if limit == 0 && count > last {
				skip = count - last
				limit = last
			}
		}
		fopts.Skip = &skip
		fopts.Limit = &limit
		// fmt.Printf("first(%d), last(%d), skip(%d), limit(%d), count(%d)\n", first, last, skip, limit, count)
	}

	return fopts
}

/*
NewFindFilterFromArgs
*/
func NewFindFilterFromArgs(args *map[string]interface{}) (*bson.D, error) {
	filter := bson.D{}
	var err error
	if args != nil {
		for key, value := range *args {
			switch {
			case key == "id":
				key = "_id"
				value, err = GetObjectIDFromValue(value)
				if err != nil {
					return nil, err
				}
			case key == "after" || key == "before":
				hid, err := Base64ToHex(value.(string))
				if err != nil {
					return nil, err
				}
				oid, err := GetObjectIDFromValue(hid)
				if err != nil {
					return nil, err
				}
				if key == "after" {
					value = bson.M{"$gt": oid}
				} else {
					value = bson.M{"$lt": oid}
				}
				key = "_id"
			case key == "first" || key == "last":
				continue
			}
			filter = append(filter, bson.E{
				Key:   key,
				Value: value,
			})
		}
	}
	return &filter, nil
}

/*
GetObjectIDFromValue tries to find a valid ObjectID whatever the value was passed
handeling 4 types (string, primitive.ObjectID, bson.Raw or a valid struct
with a nested field tagged `bson:"_id"`
*/
func GetObjectIDFromValue(value interface{}) (*primitive.ObjectID, error) {
	elem := reflect.ValueOf(value)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	iface := elem.Interface()
	switch iface.(type) {
	case bson.Raw:
		bsn := iface.(bson.Raw)
		sid := &struct {
			ID primitive.ObjectID `bson:"_id,omitempty"`
		}{}
		bson.Unmarshal(bsn, sid)
		return &sid.ID, nil
	case []byte:
		bsn := iface.([]byte)
		sid := &struct {
			ID primitive.ObjectID `bson:"_id,omitempty"`
		}{}
		bson.Unmarshal(bsn, sid)
		return &sid.ID, nil
	case string:
		str := iface.(string)
		oid, err := primitive.ObjectIDFromHex(str)
		if err != nil {
			return nil, err
		}
		return &oid, nil
	case primitive.ObjectID:
		oid, _ := iface.(primitive.ObjectID)
		return &oid, nil
	}

	if elem.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Invalid struct value %v", value)
	}

	for i := 0; i < elem.NumField(); i++ {
		elemField := elem.Field(i)
		typeField := elem.Type().Field(i)
		if tagArgs := strings.Split(typeField.Tag.Get("bson"), ","); tagArgs[0] == "_id" || tagArgs[0] == "id" {
			oid, err := GetObjectIDFromValue(elemField.Interface())
			return oid, err
		}
	}

	return nil, fmt.Errorf("Missing ObjectID inside the struct trying catch from %v", value)
}

func GetHexFromObjectID(oid *primitive.ObjectID) string {
	return oid.Hex()
}

func GetObjectIDFromHex(hex string) *primitive.ObjectID {
	oid, _ := primitive.ObjectIDFromHex(hex)
	return &oid
}

func CombineBsonArrays(args ...bson.A) bson.A {
	res := bson.A{}
	for _, stage := range args {
		res = append(res, stage...)
	}
	return res
}

func MergeBsonArrays(args ...bson.A) bson.A {
	res := bson.A{}
	for _, stage := range args {
		res = append(res, stage...)
	}
	return res
}

func AssertBsonArray(assertion bool, arr bson.A) bson.A {
	if assertion {
		return arr
	}
	return bson.A{}
}
