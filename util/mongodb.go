package util

import (
	"context"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

// oid, _ := primitive.ObjectIDFromHex(id)
// query["$or"] = append(query["$or"].([]bson.M), bson.M{"abc": "1"})
// query["$or"] = append(query["$or"].([]bson.M), bson.M{"def": "2"})
// query._id = bson.M{}
// oid, _ := primitive.ObjectIDFromHex(id)
// bson.D{{"_id", "world"}}
// query["origin"] = "test"

func LimitQueryWithId(query bson.M, before interface{}, after interface{}) bson.M {

	query["_id"] = []bson.M{}

	if before != nil {
		oid, _ := primitive.ObjectIDFromHex(before.(string))
		query["_id"] = append(query["_id"].([]bson.M), bson.M{"$lt": oid})
	}

	if after != nil {
		oid, _ := primitive.ObjectIDFromHex(after.(string))
		query["_id"] = append(query["_id"].([]bson.M), bson.M{"$gt": oid})
	}

	return query
}

type PaginationArgs struct {
	First int64
	Last  int64
	Order map[string]interface{}
}

type ConnectionArgs struct {
	Before string
	After  string
	PaginationArgs
}

func NewConnectionArgs(args map[string]interface{}) (ConnectionArgs, error) {
	conArgs := ConnectionArgs{}
	err := conArgs.Parse(args)
	return conArgs, err
}

func (ca *ConnectionArgs) Parse(args map[string]interface{}) error {
	err := FillStruct(&ca, args)
	return err
}

// Information about pagination in a connection.
type PageInfo struct {
	// When paginating forwards, are there more items?
	HasNextPage bool `json:"hasNextPage,omitempty"`
	// When paginating backwards, are there more items?
	HasPreviousPage bool `json:"hasPreviousPage,omitempty"`
	// When paginating backwards, the cursor to continue.
	StartCursor primitive.ObjectID `json:"startCursor"`
	// When paginating forwards, the cursor to continue.
	EndCursor primitive.ObjectID `json:"endCursor"`
}

// https://www.reindex.io/blog/relay-graphql-pagination-with-mongodb/
// https://stackoverflow.com/questions/51956805/how-to-page-a-cursor-in-official-mongo-go-driver
func GetPaginationParams(collection *mongo.Collection, filter interface{}, pagArgs PaginationArgs) (int, options.FindOptions, PageInfo, error) {
	var err error
	first := int64(pagArgs.First)
	last := int64(pagArgs.Last)
	order := pagArgs.Order
	pageInfo := PageInfo{}
	findOptions := options.FindOptions{}
	totalCount := 0

	if first != 0 || last != 0 {
		var count int64
		var limit int64
		var skip int64

		count, err = collection.Count(context.Background(), filter)

		if err != nil {
			return totalCount, findOptions, pageInfo, err
		}

		totalCount = int(count)
		pageInfo.HasNextPage = first != 0 && count > first
		pageInfo.HasPreviousPage = last != 0 && count > last

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
			findOptions.Skip = &skip
		}

		if limit != 0 {
			findOptions.Limit = &limit
		}

		if order != nil {
			findOptions.Sort = order
		}
	}

	return totalCount, findOptions, pageInfo, err
}

type ConnectionInfo struct {
	TotalCount int
	PageInfo   PageInfo
	Data       []bson.Raw
}

func FindWithPagination(ctx context.Context, collection *mongo.Collection, filter interface{}, conArgs ConnectionArgs) (ConnectionInfo, error) {
	ret := ConnectionInfo{}

	totalCount, findOptions, pageInfo, err := GetPaginationParams(collection, filter, PaginationArgs{
		First: conArgs.First,
		Last:  conArgs.Last,
		Order: conArgs.Order,
	})

	if err != nil {
		return ret, err
	}

	if filter == nil {
		filter = bson.M{}
	}

	cursor, err := collection.Find(ctx, filter, &findOptions)
	if err != nil {
		log.Fatal(err)
		return ret, err
	}
	defer cursor.Close(ctx)
	data := []bson.Raw{}
	for cursor.Next(ctx) {
		item := bson.Raw{}
		if err = cursor.Decode(&item); err != nil {
			log.Fatal(err)
			return ret, err
		}
		data = append(data, item)
	}

	dataLen := len(data)
	if dataLen > 0 {
		// start := data[0]
		// end := data[dataLen-1].(bson.M)
		// pageInfo.StartCursor = start["_id"].(primitive.ObjectID)
		// pageInfo.EndCursor = end["_id"].(primitive.ObjectID)
	}

	ret.Data = data
	ret.TotalCount = totalCount
	ret.PageInfo = pageInfo

	return ret, err
}
