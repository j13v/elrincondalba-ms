package operators

/*
Sync synchronize a document in the collection using an struct
*/import (
	"context"
	"github.com/j13v/elrincondalba-ms/mongodb/helpers"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func Sync(coll *mongo.Collection, ctx context.Context, mixed interface{}) error {
	oid, err := helpers.GetObjectIDFromValue(mixed)
	if err != nil {
		return err
	}
	if _, err := coll.UpdateOne(
		context.Background(),
		bson.M{"_id": oid},
		bson.M{
			"$set": mixed,
		}); err != nil {
		return err
	}
	return nil
}
