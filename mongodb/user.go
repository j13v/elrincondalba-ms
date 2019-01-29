package mongodb

import (
	"context"
	"time"

	defs "github.com/jal88/elrincondalba-ms/definitions"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type ModelUser struct {
	collection *mongo.Collection
}

func NewModelUser(db *mongo.Database) *ModelUser {
	return &ModelUser{collection: db.Collection("user")}
}

func (model *ModelUser) Create(user *defs.User) (*defs.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	val, err := bson.Marshal(user)
	if err != nil {
		return nil, err
	}
	res, err := model.collection.InsertOne(ctx, val)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, err
}

func (model *ModelUser) FindOne(args map[string]interface{}) (interface{}, error) {
	user := defs.User{}
	cursor, err := FindOne(model.collection, context.Background(), args)
	if err != nil {
		return nil, err
	}
	err = cursor.Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, err
}
