package operators

import (
	"context"
	"fmt"

	"github.com/j13v/elrincondalba-ms/mongodb/helpers"
	"github.com/mongodb/mongo-go-driver/mongo/options"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type SliceMetadata interface {
	GetTotal() int64
	GetSliceCount() int64
	GetSliceStart() int64
}

/*
FindSliceMetadata
*/
type AggregateSliceMetadata struct {
	Total      int64 `json:"total"`
	SliceCount int64 `json:"sliceCount"`
	SliceStart int64 `json:"sliceStart"`
}

func (meta *AggregateSliceMetadata) GetTotal() int64 {
	return meta.Total
}

func (meta *AggregateSliceMetadata) GetSliceCount() int64 {
	return meta.SliceCount
}

func (meta *AggregateSliceMetadata) GetSliceStart() int64 {
	return meta.SliceStart
}

func AggregateSlice(
	coll *mongo.Collection,
	ctx context.Context,
	pipeline bson.A,
) (
	result []bson.Raw,
	meta SliceMetadata,
	err error,
) {
	var count int64
	var cursor *mongo.Cursor

	count, err = AggregateCount(coll, ctx, pipeline, nil)
	if err != nil {
		return nil, nil, err
	}

	cursor, err = coll.Aggregate(ctx, pipeline)
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		item := bson.Raw{}
		cursor.Decode(&item)
		result = append(result, item)
	}

	var start int64
	if len(result) > 0 {
		oid, err := helpers.GetObjectIDFromValue(result[0])
		if err != nil {
			return nil, nil, err
		}
		index, err := AggregateGetIndexById(coll, ctx, pipeline, oid)
		if err != nil {
			return nil, nil, err
		}
		start = index
	}
	fmt.Printf("%v\n", start)
	return result, &AggregateSliceMetadata{
		Total:      count,
		SliceCount: int64(len(result)),
		SliceStart: start,
	}, nil
}

/*
GetIndexById Obtain the index from collection using id as referece
*/
func AggregateGetIndexById(coll *mongo.Collection, ctx context.Context, pipeline bson.A, id interface{}) (int64, error) {
	oid, err := helpers.GetObjectIDFromValue(id)
	if err != nil {
		return -1, err
	}
	pipeline = append(pipeline, bson.M{"$match": bson.M{"_id": bson.M{"$lt": oid}}})
	index, err := AggregateCount(coll, ctx, pipeline)
	return index, err
}

// AggregateCount takes a pipeline and count the results
func AggregateCount(
	coll *mongo.Collection,
	ctx context.Context,
	pipeline bson.A,
	restOpts ...*options.CountOptions) (int64, error) {

	var opts *options.CountOptions
	if restOpts != nil {
		opts = restOpts[0]
	}

	if opts != nil {
		if opts.Skip != nil {
			pipeline = append(pipeline, bson.M{"$skip": opts.Skip})
		}
		if opts.Limit != nil {
			pipeline = append(pipeline, bson.M{"$limit": opts.Limit})
		}
	}

	pipeline = append(pipeline, bson.M{
		"$group": bson.M{
			"_id": nil,
			"n":   bson.M{"$sum": 1},
		},
	})

	result, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, nil
	}
	countStruct := struct {
		N int64 `bson:"n"`
	}{}

	result.Next(ctx)
	result.Decode(&countStruct)

	return countStruct.N, nil
}

// count, _ := oprs.AggregateCount(coll, ctx, pipeline, nil)
// 	fmt.Printf("%v\n", count)

// 	filter, err := helpers.NewFilterFromArgs(args)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	count, err := coll.CountDocuments(ctx, args)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	fopts := helpers.NewFindOptionsFromArgs(*args, count)
// 	cursor, err := coll.Find(ctx, filter, fopts)
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, nil, err
// 	}
// 	defer cursor.Close(ctx)
// 	data := []bson.Raw{}
// 	for cursor.Next(ctx) {
// 		item := bson.Raw{}
// 		if err = cursor.Decode(&item); err != nil {
// 			log.Fatal(err)
// 			return nil, nil, err
// 		}
// 		data = append(data, item)
// 	}
// 	var start int64
// 	if len(data) > 0 {
// 		oid, err := helpers.GetObjectIDFromValue(data[0])
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		index, err := GetIndexById(coll, ctx, oid)
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		start = index
// 	}
// 	scount := *fopts.Limit
// 	return data, &FindSliceMetadata{
// 		Total:      int(count),
// 		SliceCount: int(scount),
// 		SliceStart: int(start),
// 	}, nil
