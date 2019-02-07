package models

import (
	"context"
	"log"
	"time"

	defs "github.com/jal88/elrincondalba-ms/definitions"
	oprs "github.com/jal88/elrincondalba-ms/mongodb/operators"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type ModelStock struct {
	collection *mongo.Collection
	article    *ModelArticle
}

func NewModelStock(db *mongo.Database, modelArticle *ModelArticle) *ModelStock {
	return &ModelStock{
		collection: db.Collection("stock"),
		article:    modelArticle,
	}
}

func (model *ModelStock) Create(article primitive.ObjectID, size string) (*defs.Stock, error) {
	// Check if article exists then if not raise an error
	if _, err := model.article.FindById(article); err != nil {
		return nil, err
	}
	// Create struct data
	stock, err := defs.NewStock(article, size)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	val, err := bson.Marshal(stock)
	if err != nil {
		return nil, err
	}
	// Insert stock as bson.Document
	res, err := model.collection.InsertOne(ctx, val)
	if err != nil {
		return nil, err
	}
	// Return insertedID and set to struct data
	stock.SetID(res.InsertedID.(primitive.ObjectID))
	return stock, err
}

func (model *ModelStock) FindOne(args *map[string]interface{}) (interface{}, error) {
	stock := defs.Stock{}
	cursor, err := oprs.FindOne(model.collection, context.Background(), args)
	if err != nil {
		return nil, err
	}
	err = cursor.Decode(&stock)
	if err != nil {
		return nil, err
	}
	return stock, err
}

func (model *ModelStock) FindById(id interface{}) (interface{}, error) {
	stock, err := model.FindOne(&map[string]interface{}{"_id": id.(primitive.ObjectID)})
	return stock, err
}

func (model *ModelStock) FindByArticle(article interface{}) (interface{}, error) {
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"article": article.(primitive.ObjectID),
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "order",
				"localField":   "_id",
				"foreignField": "stock",
				"as":           "order",
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path": "$order",
				"preserveNullAndEmptyArrays": true,
			},
		},
		bson.M{
			"$match": bson.M{
				"$or": bson.A{
					bson.M{"order": bson.M{"$exists": false}},
					bson.M{"order.state": bson.M{"$eq": defs.ORDER_STATUS_CANCELLED}},
				},
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"article": "$article",
					"size":    "$size",
				},
				"refs":  bson.M{"$push": "$_id"},
				"count": bson.M{"$sum": 1},
			},
		},
		bson.M{
			"$sort": bson.M{
				"_id.article": 1,
				"_id.size":    1,
			},
		},
	}
	ctx := context.Background()
	cursor, err := model.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(ctx)
	data := []defs.StockArticle{}
	for cursor.Next(ctx) {
		stock := &defs.StockArticle{}
		p := bson.M{}
		if err = cursor.Decode(&p); err != nil {
			log.Fatal(err)
			return nil, err
		}
		refs := p["refs"].(bson.A)
		for _, item := range refs {
			stock.Refs = append(stock.Refs, item.(primitive.ObjectID))
		}

		stock.Size = p["_id"].(bson.M)["size"].(string)
		stock.Count = p["count"].(int32)
		data = append(data, *stock)
	}

	return data, err
}

func (model *ModelStock) GetCount() (int64, error) {
	count, err := oprs.GetCount(model.collection, context.Background())
	return count, err
}

func (model *ModelStock) GetSizes() ([]interface{}, error) {
	sizes, err := model.collection.Distinct(context.Background(), "size", bson.D{})

	return sizes, err
}

//
// func (model *ModelStock) FindByStock(id string) ([]interface{}, error) {
// 	article := defs.Stock{}
// 	oid, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	cursor, err := FindOne(model.collection, context.Background(), map[string]interface{}{"article": oid})
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = cursor.Decode(&article)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return article, err
// }

// func (model *ModelStock) Create(article *defs.Article) (interface{}, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	val, err := bson.Marshal(article)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res, err := model.collection.InsertOne(ctx, val)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return res.InsertedID, err
// }

// func (model *ModelStock) FindOne(args map[string]interface{}) (interface{}, error) {
// 	article := defs.Article{}
// 	cursor, err := FindOne(model.collection, context.Background(), args)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = cursor.Decode(&article)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return article, err
// }

// func (model *ModelStock) FindSlice(args map[string]interface{}) ([]interface{}, *FindSliceMetadata, error) {
//
// 	data, meta, err := FindSlice(model.collection, context.Background(), args)
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, meta, err
// 	}
// 	articles := []defs.Article{}
// 	for _, v := range data {
// 		article := defs.Article{}
// 		bson.Unmarshal(v, &article)
// 		articles = append(articles, article)
// 	}
//
// 	var interfaceSlice []interface{} = make([]interface{}, len(articles))
// 	for i, d := range articles {
// 		interfaceSlice[i] = d
// 	}
//
// 	return interfaceSlice, meta, nil
// }
