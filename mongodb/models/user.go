package models

import (
	"context"
	"time"

	defs "github.com/jal88/elrincondalba-ms/definitions"
	oprs "github.com/jal88/elrincondalba-ms/mongodb/operators"
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

func (model *ModelUser) Create(dni string, name string, surname string, email string, phone string, address string) (*defs.User, error) {
	user, err := defs.NewUser(dni, name, surname, email, phone, address)
	if err != nil {
		return nil, err
	}
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

func (model *ModelUser) FindOne(args *map[string]interface{}) (interface{}, error) {
	user := defs.User{}
	cursor, err := oprs.FindOne(model.collection, context.Background(), args)
	if err != nil {
		return nil, err
	}
	err = cursor.Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (model *ModelUser) FindById(id primitive.ObjectID) (interface{}, error) {
	user, err := model.FindOne(&map[string]interface{}{"_id": id})
	return user, err
}
