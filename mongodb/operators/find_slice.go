package operators

import (
	"context"
	"log"

	"github.com/j13v/elrincondalba-ms/mongodb/helpers"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

/*
FindSliceMetadata
*/
type FindSliceMetadata struct {
	Total      int `json:"total"`
	SliceCount int `json:"sliceCount"`
	SliceStart int `json:"sliceStart"`
}

func FindSlice(coll *mongo.Collection, ctx context.Context, args *map[string]interface{}) ([]bson.Raw, *FindSliceMetadata, error) {
	filter, err := helpers.NewFindFilterFromArgs(args)
	if err != nil {
		return nil, nil, err
	}
	count, err := coll.Count(ctx, bson.D{})
	if err != nil {
		return nil, nil, err
	}
	fopts := helpers.NewFindOptionsFromArgs(*args, count)
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
		oid, err := helpers.GetObjectIDFromValue(data[0])
		if err != nil {
			return nil, nil, err
		}
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
