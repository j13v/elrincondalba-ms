package operators

import (
	"context"

	"github.com/jal88/elrincondalba-ms/mongodb/helpers"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

/*
GetIndexById Obtain the index from collection using id as referece
*/
func GetIndexById(coll *mongo.Collection, ctx context.Context, id interface{}) (int64, error) {
	oid, err := helpers.GetObjectIDFromValue(id)
	if err != nil {
		return -1, err
	}
	index, err := coll.Count(ctx, bson.D{
		{Key: "_id", Value: bson.D{{Key: "$lt", Value: oid}}},
	})
	return index, err
}
