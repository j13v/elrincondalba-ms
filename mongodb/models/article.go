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

type ModelArticle struct {
	collection *mongo.Collection
}

func NewModelArticle(db *mongo.Database) *ModelArticle {
	return &ModelArticle{collection: db.Collection("article")}
}

func (model *ModelArticle) Create(name string, description string, price float64, images []string, category string, rating int8) (*defs.Article, error) {
	article, err := defs.NewArticle(name, description, price, images, category, rating)
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
	article.ID = res.InsertedID.(primitive.ObjectID)
	return article, err
}

func (model *ModelArticle) FindOne(args *map[string]interface{}) (interface{}, error) {
	article := defs.Article{}
	cursor, err := oprs.FindOne(model.collection, context.Background(), args)
	if err != nil {
		return nil, err
	}
	err = cursor.Decode(&article)
	if err != nil {
		return nil, err
	}
	return article, err
}

func (model *ModelArticle) FindById(id primitive.ObjectID) (interface{}, error) {
	article, err := model.FindOne(&map[string]interface{}{"_id": id})
	return article, err
}

func (model *ModelArticle) FindSlice(args *map[string]interface{}) ([]interface{}, *oprs.FindSliceMetadata, error) {
	data, meta, err := oprs.FindSlice(model.collection, context.Background(), args)
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

func (model *ModelArticle) DeleteArticle(id primitive.ObjectID) error {

	if _, err := model.collection.DeleteOne(
		context.Background(),
		map[string]interface{}{"_id": id},
	); err != nil {
		return err
	}
	return nil
}

func (model *ModelArticle) UpdateArticle(id primitive.ObjectID, name string, description string, price float64, images []string, category string, rating int8) error {

	if _, err := model.collection.UpdateOne(
		context.Background(),
		map[string]interface{}{"_id": id},
		bson.M{
			"$set": bson.M{"name": name},
		},
	); err != nil {
		return err
	}
	return nil
}

// func (model *ModelArticle) SetStock(id sring, size string, count int) error {
//
// }

func (model *ModelArticle) GetCount() (int64, error) {
	count, err := oprs.GetCount(model.collection, context.Background())
	return count, err
}

func (model *ModelArticle) GetCategories() ([]interface{}, error) {
	categories, err := model.collection.Distinct(context.Background(), "category", bson.D{})

	return categories, err
}

//TODO define filters params
func (model *ModelArticle) GetMinMaxPrice() (interface{}, error) {

	pipeline := bson.A{
		bson.M{
			"$group": bson.M{
				"_id": nil,
				"max": bson.M{"$max": "$price"},
				"min": bson.M{"$min": "$price"},
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
	cursor.Next(ctx)
	data := struct {
		Min float64 `bson:"min" json:"min"`
		Max float64 `bson:"max" json:"max"`
	}{}
	cursor.Decode(&data)
	return data, err
}
