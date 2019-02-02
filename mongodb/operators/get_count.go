package operators

/*
GetCount return the number of documents in a collection
*/import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func GetCount(coll *mongo.Collection, ctx context.Context) (int64, error) {
	count, err := coll.Count(ctx, bson.D{})
	return count, err
}
