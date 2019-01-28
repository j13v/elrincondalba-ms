package mongodb

import (
	"context"
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

func (model *ModelOrder) FindById(id string) (interface{}, error) {
	order := defs.Order{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	oid, _ := primitive.ObjectIDFromHex(id)
	err := model.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&order)
	return order, err
}

func (model *ModelOrder) FindSlice(args *FindSliceArguments) ([]interface{}, error) {

	data, err := FindSliceData(model.collection, context.Background(), args)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	orders := []defs.Order{}
	for _, v := range data {
		order := defs.Order{}
		bson.Unmarshal(v, &order)
		// util.FillStruct(&order, v)
		orders = append(orders, order)
		// fmt.Printf("%+v\n", order)
	}

	var interfaceSlice []interface{} = make([]interface{}, len(orders))
	for i, d := range orders {
		interfaceSlice[i] = d
	}

	return interfaceSlice, nil
}

func (model *ModelOrder) GetCount() (int64, error) {
	count, err := GetCount(model.collection, context.Background())
	return count, err
}

// func (model *ModelOrder) ListOrders() ([]TypeOrder, error) {
//   orders := []TypeOrder{}
// 	cursor, err := model.collection.Find(context.Background(), bson.M{})
//   if err != nil {
//     log.Fatal(err)
//     return nil, err
//   }
//   defer cursor.Close(context.Background())
//   for cursor.Next(context.Background()) {
//     order := TypeOrder{}
//     err := cursor.Decode(&order)
//     if err != nil {
//       log.Fatal(err)
//       return nil, err;
//     }
//     orders = append(orders, order)
//   }
// 	return orders, nil
// }

// func (model *ModelOrder) FindWithPagination(args map[string]interface{}) ([]interface{}, error) {
// 	conArgs, err := util.NewConnectionArgs(args)
// 	ctx := context.Background()
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}
//
// 	connectionData, err := util.FindWithPagination(ctx, model.collection, bson.M{}, conArgs)
// 	orders := []defs.Order{}
// 	for _, v := range connectionData.Data {
// 		order := defs.Order{}
// 		bson.Unmarshal(v, &order)
// 		// util.FillStruct(&order, v)
// 		orders = append(orders, order)
// 		// fmt.Printf("%+v\n", order)
// 	}
//
// 	var interfaceSlice []interface{} = make([]interface{}, len(orders))
// 	for i, d := range orders {
// 		interfaceSlice[i] = d
// 	}
//
// 	fmt.Printf("%+v\n", orders)
//
// 	return interfaceSlice, nil
// }
