package mongodb

import (
	"context"
	"time"

	defs "github.com/jal88/elrincondalba-ms/definitions"
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

func (model *ModelStock) FindOne(args map[string]interface{}) (interface{}, error) {
	stock := defs.Stock{}
	cursor, err := FindOne(model.collection, context.Background(), args)
	if err != nil {
		return nil, err
	}
	err = cursor.Decode(&stock)
	if err != nil {
		return nil, err
	}
	return stock, err
}

func (model *ModelStock) FindById(id primitive.ObjectID) (interface{}, error) {
	stock, err := model.FindOne(map[string]interface{}{"_id": id})
	return stock, err
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

func (model *ModelStock) GetCount() (int64, error) {
	count, err := GetCount(model.collection, context.Background())
	return count, err
}
