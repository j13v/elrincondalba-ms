package operators

import (
	"context"

	"github.com/j13v/elrincondalba-ms/mongodb/helpers"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func FindOne(coll *mongo.Collection, ctx context.Context, args *map[string]interface{}) (*mongo.SingleResult, error) {
	filter, err := helpers.NewFindFilterFromArgs(args)
	if err != nil {
		return nil, err
	}
	cursor := coll.FindOne(ctx, filter)
	return cursor, nil
}
