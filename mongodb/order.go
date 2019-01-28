package mongodb

import (
	"context"
	"log"
	"time"

	defs "github.com/jal88/elrincondalba-ms/definitions"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

/*
ModelOrder asd
*/
type ModelOrder struct {
	collection *mongo.Collection
}

/*
NewOrderModel
*/
func NewModelOrder(db *mongo.Database) *ModelOrder {
	return &ModelOrder{collection: db.Collection("order")}
}

/*
Create
*/
func (model *ModelOrder) Create(order *defs.Order) (interface{}, error) {
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

func (model *ModelOrder) FindOne(args map[string]interface{}) (interface{}, error) {
	order := defs.Article{}
	cursor, err := FindOne(model.collection, context.Background(), args)
	if err != nil {
		return nil, err
	}
	err = cursor.Decode(&order)
	if err != nil {
		return nil, err
	}
	return order, err
}

func (model *ModelOrder) FindSlice(args map[string]interface{}) ([]interface{}, *FindSliceMetadata, error) {

	data, meta, err := FindSlice(model.collection, context.Background(), args)
	if err != nil {
		log.Fatal(err)
		return nil, meta, err
	}
	orders := []defs.Article{}
	for _, v := range data {
		order := defs.Article{}
		bson.Unmarshal(v, &order)
		orders = append(orders, order)
	}

	var interfaceSlice []interface{} = make([]interface{}, len(orders))
	for i, d := range orders {
		interfaceSlice[i] = d
	}

	return interfaceSlice, meta, nil
}

func (model *ModelOrder) GetCount() (int64, error) {
	count, err := GetCount(model.collection, context.Background())
	return count, err
}
