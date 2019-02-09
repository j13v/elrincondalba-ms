package models

import (
	"context"
	"io"
	"log"
	"time"

	defs "github.com/j13v/elrincondalba-ms/definitions"
	oprs "github.com/j13v/elrincondalba-ms/mongodb/operators"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/gridfs"
)

type ModelArticle struct {
	collection *mongo.Collection
	bucket     *gridfs.Bucket
}

func NewModelArticle(db *mongo.Database) *ModelArticle {
	bucket, _ := gridfs.NewBucket(db)
	return &ModelArticle{
		collection: db.Collection("article"),
		bucket:     bucket,
	}
}

func (model *ModelArticle) Create(name string, description string, price float64, images []defs.File, category string, rating int8) (*defs.Article, error) {
	imageIds := []primitive.ObjectID{}
	for _, image := range images {
		imageId, err := model.UploadImage(image.Filename, image.File)
		if err != nil {
			return nil, err
		}
		imageIds = append(imageIds, imageId)
	}
	article, err := defs.NewArticle(name, description, price, imageIds, category, rating)
	if err != nil {
		return nil, err
	}
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

func (model *ModelArticle) UploadImage(filename string, source io.Reader) (primitive.ObjectID, error) {
	id, err := model.bucket.UploadFromStream(filename, source)
	return id, err
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

func (model *ModelArticle) Sync(article *defs.Article) error {
	if _, err := model.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": article.ID},
		bson.M{
			"$set": article,
		}); err != nil {
		return err
	}

	cursor, err := oprs.FindById(model.collection, context.Background(), article.ID)
	cursor.Decode(article)

	return err
}

func (model *ModelArticle) GetCount() (int64, error) {
	count, err := oprs.GetCount(model.collection, context.Background())
	return count, err
}

func (model *ModelArticle) GetCategories() ([]interface{}, error) {
	categories, err := model.collection.Distinct(context.Background(), "category", bson.D{})

	return categories, err
}

func (model *ModelArticle) GetPriceRange() (interface{}, error) {

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
		Min float64 `bson:"min"`
		Max float64 `bson:"max"`
	}{}
	cursor.Decode(&data)
	return []float64{data.Min, data.Max}, err
}
