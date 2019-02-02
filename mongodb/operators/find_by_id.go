package operators

import (
	"context"

	"github.com/jal88/elrincondalba-ms/mongodb/helpers"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func FindById(coll *mongo.Collection, ctx context.Context, id interface{}) (*mongo.SingleResult, error) {
	oid, err := helpers.GetObjectIDFromValue(id)
	if err != nil {
		return nil, err
	}
	cursor := coll.FindOne(ctx, bson.M{"_id": oid})
	return cursor, err
}
