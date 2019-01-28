package mongodb

import (
	"context"
	"log"
	"reflect"
	"strings"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

/*
FindSliceArguments
*/
type FindSliceArguments struct {
	Filter      interface{}
	FindOptions *options.FindOptions
}

/*
ParseConnectionArgs
*/
func NewFindArgs(args map[string]interface{}, count int64) *FindSliceArguments {
	filter := bson.D{}
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
			}
		}

		if skip != 0 {
			fopts.Skip = &skip
		}

		if limit != 0 {
			fopts.Limit = &limit
		}
		// if before, ok := args["before"]; ok {
		// oid, _ := primitive.ObjectIDFromHex(before.(string))
		// filter["_id"] = bson.M{filter["_id"], "$gt": oid}
		// filter.Append("_id", bsonx.Val{
		// 	Key:   "$lt",
		// 	Value: oid,
		// })
		// }
		// if after, ok := args["after"]; ok {
		// 	oid, _ := primitive.ObjectIDFromHex(after.(string))
		// 	filter = append(filter, bson.M{"_id": bson.M{"$gt": oid}})
		// }
	}

	// newone, _ := bson.Marshal(filter)
	return &FindSliceArguments{
		Filter:      &filter,
		FindOptions: fopts,
	}
}

func FindSliceData(coll *mongo.Collection, ctx context.Context, args *FindSliceArguments) ([]bson.Raw, error) {
	cursor, err := coll.Find(ctx, args.Filter, args.FindOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(ctx)
	data := []bson.Raw{}
	for cursor.Next(ctx) {
		item := bson.Raw{}
		if err = cursor.Decode(&item); err != nil {
			log.Fatal(err)
			return nil, err
		}
		data = append(data, item)
	}
	return data, nil
}

func GetCount(coll *mongo.Collection, ctx context.Context) (int64, error) {
	count, err := coll.Count(ctx, bson.D{})
	return count, err
}

func GetHexFromObjectID(oid primitive.ObjectID) string {
	return oid.Hex()
}

func GetObjectIDFromValue(value interface{}) primitive.ObjectID {
	oid := primitive.ObjectID{}
	elem := reflect.ValueOf(value)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	for i := 0; i < elem.NumField(); i++ {
		elemField := elem.Field(i)
		if typeField := elem.Type().Field(i); typeField.Name == "ID" {
			if str, ok := elemField.Interface().(primitive.ObjectID); ok {
				return str
			}
			if tagArgs := strings.Split(typeField.Tag.Get("bson"), ","); tagArgs[0] == "_id" {
				if str, ok := elemField.Interface().(primitive.ObjectID); ok {
					return str
				}
			}
			if typeField.Type == reflect.TypeOf(primitive.ObjectID{}) {
				if str, ok := elemField.Interface().(primitive.ObjectID); ok {
					return str
				}
			}
		}

	}
	return oid
}

//
// func LimitQueryWithId(query bson.M, before interface{}, after interface{}) bson.M {
//
// 	query["_id"] = []bson.M{}
//
// 	if before != nil {
// 		oid, _ := primitive.ObjectIDFromHex(before.(string))
// 		query["_id"] = append(query["_id"].([]bson.M), bson.M{"$lt": oid})
// 	}
//
// 	if after != nil {
// 		oid, _ := primitive.ObjectIDFromHex(after.(string))
// 		query["_id"] = append(query["_id"].([]bson.M), bson.M{"$gt": oid})
// 	}
//
// 	return query
// }

//
// var err error
// first := int64(pagArgs.First)
// last := int64(pagArgs.Last)
// order := pagArgs.Order
// pageInfo := PageInfo{}
// findOptions := options.FindOptions{}
// totalCount := 0
//
// if first != 0 || last != 0 {
// 	var count int64
// 	var limit int64
// 	var skip int64
//
// 	count, err = collection.Count(context.Background(), filter)
//
// 	if err != nil {
// 		return totalCount, findOptions, pageInfo, err
// 	}
//
// 	totalCount = int(count)
// 	pageInfo.HasNextPage = first != 0 && count > first
// 	pageInfo.HasPreviousPage = last != 0 && count > last
//
// if first != 0 && count > first {
// 	limit = first
// }
//
// if last != 0 {
// 	if limit != 0 && limit > last {
// 		skip = limit - last
// 		limit = limit - skip
// 	} else if limit == 0 && count > last {
// 		skip = count - last
// 	}
// }
//
// 	if skip != 0 {
// 		findOptions.Skip = &skip
// 	}
//
// 	if limit != 0 {
// 		findOptions.Limit = &limit
// 	}
//
// 	if order != nil {
// 		findOptions.Sort = order
// 	}
// }
//
// return totalCount, findOptions, pageInfo, err

// if first, ok := filters["first"]; ok {
// 	if first, ok := first.(int); ok {
// 		conn.First = first
// 	}
// }
// if first, ok := filters["first"]; ok {
// 	if first, ok := first.(int); ok {
// 		conn.First = first
// 	}
// }
//
// if first, ok := filters["first"]; ok {
// 	if first, ok := first.(int); ok {
// 		conn.First = first
// 	}
// }
// if last, ok := filters["last"]; ok {
// 	if last, ok := last.(int); ok {
// 		conn.Last = last
// 	}
// }
// if before, ok := filters["before"]; ok {
// 	conn.Before = ConnectionCursor(fmt.Sprintf("%v", before))
// }
// if after, ok := filters["after"]; ok {
// 	conn.After = ConnectionCursor(fmt.Sprintf("%v", after))
// }
