package models

import (
	"context"
	"log"
	"time"

	defs "github.com/j13v/elrincondalba-ms/definitions"
	oprs "github.com/j13v/elrincondalba-ms/mongodb/operators"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
)

type ModelStock struct {
	collection *mongo.Collection
	article    *ModelArticle
}

func NewModelStock(db *mongo.Database, modelArticle *ModelArticle) *ModelStock {
	model := &ModelStock{
		collection: db.Collection("stock"),
		article:    modelArticle,
	}

	model.ensureIndex()
	return model
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

func (model *ModelStock) FindAvailableByArticle(article interface{}) (interface{}, error) {
	ctx := context.Background()
	cursor, err := model.collection.Aggregate(ctx,
		combinePipeline(
			bson.A{
				bson.M{
					"$match": bson.M{
						"article": bson.M{"$eq": article.(primitive.ObjectID)},
					},
				},
			},
			pipelineStockOrder,
			pipelineStockOrderAvailable,
			pipelineStockArticleGroup))

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(ctx)
	data := []defs.StockArticle{}
	for cursor.Next(ctx) {
		stock := &defs.StockArticle{}
		if err = cursor.Decode(&stock); err != nil {
			log.Fatal(err)
			return nil, err
		}
		data = append(data, *stock)
	}
	return data, err
}

func (model *ModelStock) FindByArticle(article interface{}) (interface{}, error) {
	ctx := context.Background()
	cursor, err := model.collection.Aggregate(ctx,
		combinePipeline(
			bson.A{
				bson.M{
					"$match": bson.M{
						"article": bson.M{"$eq": article.(primitive.ObjectID)},
					},
				},
			},
			pipelineStockOrder,
			pipelineStockOrderArticleGroup))

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(ctx)
	// TODO make function that group this array result
	data := []defs.StockOrderArticle{}
	for cursor.Next(ctx) {
		stock := &defs.StockOrderArticle{}
		if err = cursor.Decode(&stock); err != nil {
			log.Fatal(err)
			return nil, err
		}
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
// func (model *ModelStock) GetAvailableStockByArticle(article interface{}) ([]interface{}, error) {
// 	// article.(primitive.ObjectID)
// 	return nil, nil
// }

// func (model *ModelStock) GetAvailableStock() ([]interface{}, error) {
// 	model.collection.Database().Collection("availableStock").Find(ctx, filter)
// 	return nil, nil
// }

func (model *ModelStock) FindSlice(args *map[string]interface{}) ([]interface{}, *oprs.FindSliceMetadata, error) {
	collection := model.collection.Database().Collection("availableStock")
	data, meta, err := oprs.FindSlice(collection, context.Background(), args)
	if err != nil {
		log.Fatal(err)
		return nil, meta, err
	}
	articles := []interface{}{}
	for _, v := range data {
		article := bson.D{}
		bson.Unmarshal(v, &article)
		articles = append(articles, article)
	}

	return articles, meta, nil
}

func (model *ModelStock) ensureIndex() error {
	indexView := model.collection.Indexes()
	_, err := indexView.CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{"article": bsonx.Int32(1), "order": bsonx.Int32(1)},
	})
	return err
}

// "_id": {
// 	"article": "$article",
// 	"size":    "$size",
//  },
//  "count": {"$sum": 1},
//  "refs": {"$push": "$_id"},

// bson.M{
// 	"$group": bson.M{
// 		"_id": bson.M{
// 			"article": "$article",
// 			"size":    "$size",
// 		},
// 		"refs":  bson.M{"$push": "$_id"},
// 		"count": bson.M{"$sum": 1},
// 	},
// },
// bson.M{
// 	"$project": bson.M{
// 		"article": "$_id.article",
// 		"size":    "$_id.size",
// 		"_id":     0,
// 		"refs":    1,
// 		"count":   1,
// 	},
// },
// bson.M{
// 	"$group": bson.M{
// 		"_id": bson.M{
// 			"article": "$article",
// 			"size":    "$size",
// 		},
// 		"refs":  bson.M{"$push": "$_id"},
// 		"count": bson.M{"$sum": 1},
// 	},
// },
// bson.M{
// 	"$sort": bson.M{
// 		"_id.article": 1,
// 		"_id.size":    1,
// 	},
// },
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
