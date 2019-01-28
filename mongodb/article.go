package mongodb

import (
	"context"
	"log"
	"time"

	defs "github.com/jal88/elrincondalba-ms/definitions"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type ModelArticle struct {
	collection *mongo.Collection
}

func NewModelArticle(db *mongo.Database) *ModelArticle {
	return &ModelArticle{collection: db.Collection("article")}
}

func (model *ModelArticle) Create(article *defs.Article) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	val, err := bson.Marshal(article)
	if err != nil {
		return nil, err
	}
	res, err := model.collection.InsertOne(ctx, val)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, err
}

func (model *ModelArticle) FindOne(args map[string]interface{}) (interface{}, error) {
	article := defs.Article{}
	cursor, err := FindOne(model.collection, context.Background(), args)
	if err != nil {
		return nil, err
	}
	err = cursor.Decode(&article)
	if err != nil {
		return nil, err
	}
	return article, err
}

func (model *ModelArticle) FindSlice(args map[string]interface{}) ([]interface{}, *FindSliceMetadata, error) {

	data, meta, err := FindSlice(model.collection, context.Background(), args)
	if err != nil {
		log.Fatal(err)
		return nil, meta, err
	}
	articles := []defs.Article{}
	for _, v := range data {
		article := defs.Article{}
		bson.Unmarshal(v, &article)
		articles = append(articles, article)
	}

	var interfaceSlice []interface{} = make([]interface{}, len(articles))
	for i, d := range articles {
		interfaceSlice[i] = d
	}

	return interfaceSlice, meta, nil
}

func (model *ModelArticle) GetCount() (int64, error) {
	count, err := GetCount(model.collection, context.Background())
	return count, err
}
