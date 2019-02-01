package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	defs "github.com/jal88/elrincondalba-ms/definitions"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

/*
ModelOrder asd
*/
type ModelOrder struct {
	collection *mongo.Collection
	stock      *ModelStock
	user       *ModelUser
}

/*
NewOrderModel
*/
func NewModelOrder(db *mongo.Database, modelStock *ModelStock, modelUser *ModelUser) *ModelOrder {
	return &ModelOrder{
		collection: db.Collection("order"),
		stock:      modelStock,
		user:       modelUser,
	}
}

/*
Create
*/
func (model *ModelOrder) Create(stock primitive.ObjectID, user primitive.ObjectID, notes string) (*defs.Order, error) {
	// Check if article exists then if not raise an error
	if _, err := model.stock.FindById(stock); err != nil {
		return nil, fmt.Errorf("No stock found by id %s", stock)
	}
	if _, err := model.user.FindById(user); err != nil {
		return nil, fmt.Errorf("No user found by id %s", user)
	}
	order, err := defs.NewOrder(stock, user, notes)
	if err != nil {
		return nil, err
	}
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
	order.ID = res.InsertedID.(primitive.ObjectID)
	return order, err
}

func (model *ModelOrder) FindOne(args map[string]interface{}) (interface{}, error) {
	order := defs.Order{}
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

func (model *ModelOrder) FindById(id primitive.ObjectID) (interface{}, error) {
	order, err := model.FindOne(map[string]interface{}{"_id": id})
	return order, err
}

func (model *ModelOrder) FindSlice(args map[string]interface{}) ([]interface{}, *FindSliceMetadata, error) {

	data, meta, err := FindSlice(model.collection, context.Background(), args)
	if err != nil {
		log.Fatal(err)
		return nil, meta, err
	}
	orders := []defs.Order{}
	for _, v := range data {
		order := defs.Order{}
		bson.Unmarshal(v, &order)
		orders = append(orders, order)
	}

	var interfaceSlice []interface{} = make([]interface{}, len(orders))
	for i, d := range orders {
		interfaceSlice[i] = d
	}

	return interfaceSlice, meta, nil
}

func (model *ModelOrder) UpdateState(id primitive.ObjectID, state int8) error {
	currentState, err := model.GetState(id)
	if err != nil {
		return err
	}
	if err := defs.ValidateNextOrderState(currentState, state); err != nil {
		return err
	}
	if _, err := model.collection.UpdateOne(
		context.Background(),
		map[string]interface{}{"_id": id},
		bson.M{
			"$set": bson.M{"state": state},
		},
	); err != nil {
		return err
	}
	return nil
}
func (model *ModelOrder) GetState(id primitive.ObjectID) (int8, error) {
	order, err := model.FindById(id)
	fmt.Printf("Yo pasaba por acá en la jurisdiccion %v %v\n", order, err)
	if err != nil {
		return 0, err
	}

	orderStruct := order.(defs.Order)
	return orderStruct.State, nil
}

func (model *ModelOrder) GetCount() (int64, error) {
	count, err := GetCount(model.collection, context.Background())
	return count, err
}
