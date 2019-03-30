package models

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	defs "github.com/j13v/elrincondalba-ms/definitions"
	"github.com/j13v/elrincondalba-ms/mongodb/helpers"
	oprs "github.com/j13v/elrincondalba-ms/mongodb/operators"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/gridfs"
	"github.com/mongodb/mongo-go-driver/mongo/options"
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

func (model *ModelArticle) AddStock(article primitive.ObjectID, size string) (*defs.Stock, error) {
	// Create struct data
	stock, err := defs.NewStock(size)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = model.collection.UpdateOne(ctx, bson.M{
		"_id": article,
	}, bson.M{
		"$push": bson.M{"stock": stock},
	})

	if err != nil {
		return nil, err
	}

	return stock, err
}

func (model *ModelArticle) UploadImage(filename string, source io.Reader) (primitive.ObjectID, error) {
	id, err := model.bucket.UploadFromStream(filename, source)
	return id, err
}
func (model *ModelArticle) FindOne(args *map[string]interface{}, opts ...*options.FindOneOptions) (interface{}, error) {
	fmt.Printf("%x", args)
	article := defs.Article{}
	cursor, err := oprs.FindOne(model.collection, context.Background(), args, opts...)
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

func (model *ModelArticle) FindStockById(stockId primitive.ObjectID) (interface{}, error) {
	cursor, err := oprs.FindOne(model.collection, context.Background(), &map[string]interface{}{"stock._id": stockId}, &options.FindOneOptions{
		Projection: bson.M{"stock.$": 1},
	})
	if err != nil {
		return nil, err
	}
	res := struct {
		Stock []defs.Stock `bson: "stock"`
	}{}

	err = cursor.Decode(&res)
	if err != nil {
		return nil, err
	}
	if res.Stock == nil || len(res.Stock) == 0 {
		return nil, fmt.Errorf("Stock not found with id %v", stockId)
	}
	stock := res.Stock[0]
	return stock, err
}

func (model *ModelArticle) FindSlice(args *map[string]interface{}) (
	result []defs.Article,
	meta oprs.SliceMetadata,
	err error,
) {
	var bsonData []bson.Raw
	ctx := context.Background()
	bsonData, meta, err = oprs.AggregateSlice(model.collection, ctx, bson.A{})
	for _, v := range bsonData {
		article := defs.Article{}
		bson.Unmarshal(v, &article)
		result = append(result, article)
	}
	return result, meta, err
}

// func (model *ModelArticle) FindSlice(args *map[string]interface{}) (interfaceSlice []interface{}, meta *oprs.FindSliceMetadata, err error) {
// 	var bsonData []bson.Raw
// 	if args, err = NewFindFilterSliceFromArgs(args); err != nil {
// 		return nil, nil, err
// 	}
// 	bsonData, meta, err = oprs.FindSlice(model.collection, context.Background(), args)
// 	if err != nil {
// 		return nil, meta, err
// 	}
// 	articles := []defs.Article{}
// 	for _, v := range bsonData {
// 		article := defs.Article{}
// 		bson.Unmarshal(v, &article)
// 		articles = append(articles, article)
// 	}

// 	interfaceSlice = make([]interface{}, len(articles))
// 	for i, d := range articles {
// 		interfaceSlice[i] = d
// 	}

// 	return interfaceSlice, meta, nil
// }

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

type DistincItem struct {
	Name  string `bson:"name"`
	Count int64  `bson:"count"`
}

func (model *ModelArticle) DistincUsingFilters(path string, filterPipline bson.A) ([]DistincItem, error) {
	pipeline := helpers.CombineBsonArrays(bson.A{
		bson.M{
			"$unwind": bson.M{
				"path": "$stock",
				"preserveNullAndEmptyArrays": true,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id":   path,
				"stock": bson.M{"$push": "$$ROOT"},
			},
		},
	},
		helpers.AssertBsonArray(true, filterPipline),
		bson.A{
			bson.M{
				"$project": bson.M{
					"name":  path,
					"count": bson.M{"$size": "$stock"},
				},
			},
			bson.M{
				"$sort": bson.M{
					"name": 1,
				},
			},
		})

	ctx := context.Background()
	cursor, err := model.collection.Aggregate(ctx, pipeline)

	if err != nil {

		return nil, err
	}
	defer cursor.Close(ctx)
	data := []DistincItem{}
	for cursor.Next(ctx) {
		stock := &DistincItem{}
		if err = cursor.Decode(&stock); err != nil {
			log.Fatal(err)
			return nil, err
		}
		data = append(data, *stock)
	}
	return data, err
}

func (model *ModelArticle) GetCategories(args *map[string]interface{}) (interface{}, error) {
	data, err := model.DistincUsingFilters("$category", assertPipeline((*args)["sizes"] != nil, bson.A{
		bson.M{
			"$project": bson.M{
				"_id":      0,
				"category": "$_id",
				"stock": bson.M{
					"$filter": bson.M{
						"input": "$stock",
						"as":    "article",
						"cond": bson.M{
							"$and": bson.A{
								bson.M{"$in": bson.A{"$$article.stock.size", (*args)["sizes"]}},
							},
						},
					},
				},
			},
		},
	}))
	return data, err
}

func (model *ModelArticle) GetSizes(args *map[string]interface{}) (interface{}, error) {
	data, err := model.DistincUsingFilters("$size", assertPipeline((*args)["categories"] != nil, bson.A{
		bson.M{
			"$project": bson.M{
				"_id":      0,
				"category": "$_id",
				"stock": bson.M{
					"$filter": bson.M{
						"input": "$stock",
						"as":    "article",
						"cond": bson.M{
							"$and": bson.A{
								bson.M{"$in": bson.A{"$$article.category", (*args)["categories"]}},
							},
						},
					},
				},
			},
		},
	}))
	return data, err
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

func castArrayString(input interface{}) (out bson.A) {
	value := input.([]interface{})
	for _, v := range value {
		// using a type assertion, convert v to a string
		out = append(out, v.(string))
	}
	return out
}

func castArrayFloat(input interface{}) (out []float64) {
	value := input.([]interface{})
	for _, v := range value {
		out = append(out, v.(float64))
	}
	return out
}

func NewFindFilterSliceFromArgs(args *map[string]interface{}) (*map[string]interface{}, error) {
	resArgs := map[string]interface{}{}
	if args != nil {
		for key, value := range *args {
			switch {
			case key == "categories":
				key = "category"
				bsonArr := castArrayString(value)
				if len(bsonArr) == 0 {
					continue
				}
				value = bson.M{
					"$in": value,
				}
			case key == "priceRange":
				key = "price"
				price := castArrayFloat(value)
				if len(price) < 2 {
					continue
				}
				value = bson.M{
					"$gte": price[0],
					"$lte": price[1],
				}
			case key == "sizes":
				key = "stock.size"
				bsonArr := castArrayString(value)
				if len(bsonArr) == 0 {
					continue
				}
				value = bson.M{
					"$in": value,
				}
			case key == "after" || key == "before" || key == "first" || key == "last" || key == "id":
			default:
				continue
			}
			resArgs[key] = value

		}
	}
	return &resArgs, nil
}
