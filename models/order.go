package models

import (
	"context"
  "log"
	"github.com/mongodb/mongo-go-driver/bson"
  "github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"time"
)

type TypeOrder struct {
	ID          primitive.ObjectID    `bson:"_id,omitempty" json:"id,omitempty"`
  Article     string                `bson:"article" json:"article"`
  User        string                `bson:"user" json:"user"`
  Size        string                `bson:"size" json:"size"`
  CreateAt    int32                 `bson:"createAt" json:"createAt"`
  UpdateAt    int32                 `bson:"updateAt" json:"updateAt"`
  State       int8                  `bson:"state" json:"state"`
}

type ModelOrder struct {
	collection *mongo.Collection
}

func GetModelOrder(db *mongo.Database) ModelOrder {
	return ModelOrder{collection: db.Collection("order")}
}

func (model *ModelOrder) CreateOrder(order TypeOrder) (interface{}, interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	val, err := bson.Marshal(order)
	if err != nil {
		return nil, err
	}
	res, err := model.collection.InsertOne(ctx, val)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, err
}

func (model *ModelOrder) GetOrderById(id string) (TypeOrder, error) {
	order := TypeOrder{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
  oid, _ := primitive.ObjectIDFromHex(id)
	err := model.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&order)
	return order, err
}

func (model *ModelOrder) ListOrders() ([]TypeOrder, error) {
  orders := []TypeOrder{}
	cursor, err := model.collection.Find(context.Background(), bson.M{})
  if err != nil {
    log.Fatal(err)
    return nil, err
  }
  defer cursor.Close(context.Background())
  for cursor.Next(context.Background()) {
    order := TypeOrder{}
    err := cursor.Decode(&order)
    if err != nil {
      log.Fatal(err)
      return nil, err;
    }
    orders = append(orders, order)
  }
	return orders, nil
}
