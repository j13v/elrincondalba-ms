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
FindSliceMetadata
*/
type FindSliceMetadata struct {
	Total      int `json:"total"`
	SliceCount int `json:"sliceCount"`
	SliceStart int `json:"sliceStart"`
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
func NewFindFilterFromArgs(args map[string]interface{}) (*bson.D, error) {
	filter := bson.D{}

	if args != nil {
		for key, value := range args {
			switch {
			case key == "id":
				oid, err := primitive.ObjectIDFromHex(value.(string))
				if err != nil {
					return &filter, err
				}
				key = "_id"
				value = oid
			case key == "after" || key == "before":
				oid, err := primitive.ObjectIDFromHex(value.(string))
				if err != nil {
					return &filter, err
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

func FindById(coll *mongo.Collection, ctx context.Context, id string) (*mongo.SingleResult, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	cursor := coll.FindOne(ctx, bson.M{"_id": oid})
	return cursor, err
}

func FindOne(coll *mongo.Collection, ctx context.Context, args map[string]interface{}) (*mongo.SingleResult, error) {
	filter, err := NewFindFilterFromArgs(args)
	if err != nil {
		return nil, err
	}
	cursor := coll.FindOne(ctx, filter)
	return cursor, nil
}

func FindSlice(coll *mongo.Collection, ctx context.Context, args map[string]interface{}) ([]bson.Raw, *FindSliceMetadata, error) {
	filter, err := NewFindFilterFromArgs(args)
	if err != nil {
		return nil, nil, err
	}
	count, err := coll.Count(ctx, bson.D{})
	if err != nil {
		return nil, nil, err
	}
	fopts := NewFindOptionsFromArgs(args, count)
	cursor, err := coll.Find(ctx, filter, fopts)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	defer cursor.Close(ctx)
	data := []bson.Raw{}
	for cursor.Next(ctx) {
		item := bson.Raw{}
		if err = cursor.Decode(&item); err != nil {
			log.Fatal(err)
			return nil, nil, err
		}
		data = append(data, item)
	}
	var start int64
	if len(data) > 0 {
		oid := GetIdFromRawBson(data[0])
		index, err := GetIndexById(coll, ctx, oid)
		if err != nil {
			return nil, nil, err
		}
		start = index
	}
	scount := *fopts.Limit
	return data, &FindSliceMetadata{
		Total:      int(count),
		SliceCount: int(scount),
		SliceStart: int(start),
	}, nil
}

func GetIdFromRawBson(d bson.Raw) *primitive.ObjectID {
	sid := &struct {
		ID primitive.ObjectID `bson:"_id,omitempty"`
	}{}
	bson.Unmarshal(d, sid)
	return &sid.ID
}

func GetIndexById(coll *mongo.Collection, ctx context.Context, id *primitive.ObjectID) (int64, error) {
	index, err := coll.Count(ctx, bson.D{
		{Key: "_id", Value: bson.D{{Key: "$lt", Value: id}}},
	})
	return index, err
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
