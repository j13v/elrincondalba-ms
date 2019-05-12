package models

import (
	"context"
	"io"
	"log"
	"time"
	"fmt"

	defs "github.com/j13v/elrincondalba-ms/definitions"
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

func (model *ModelArticle) FindStockById(stockId primitive.ObjectID) (*defs.StockArticle, error) {
	pipeline := bson.A{
		bson.M{
			"$unwind": bson.M{
				"path": "$stock",
			},
		},
		bson.M{
			"$replaceRoot": bson.M{
				"newRoot": bson.M{"$mergeObjects": bson.A{
					"$stock",
					bson.M{
						"article": "$$ROOT",
					},
				}},
			},
		},
		bson.M{
			"$match": bson.M{
				"_id": stockId,
			},
		},
		bson.M{
			"$project": bson.M{
				"article.stock": false,
			},
		},
	}

	ctx := context.Background()
	cursor, err := model.collection.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	stockArticle := defs.StockArticle{}
	cursor.Next(ctx)
	if err = cursor.Decode(&stockArticle); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &stockArticle, err
}

func (model *ModelArticle) FindStockBySize(stockSize string) (*defs.StockArticle, error) {
	// pipeline := bson.A{
	// 	bson.M{
	// 		"$unwind": bson.M{
	// 			"path": "$stock",
	// 		},
	// 	},
	// 	bson.M{
	// 		"$replaceRoot": bson.M{
	// 			"newRoot": bson.M{"$mergeObjects": bson.A{
	// 				"$stock",
	// 				bson.M{
	// 					"article": "$$ROOT",
	// 				},
	// 			}},
	// 		},
	// 	},
	// 	bson.M{
	// 		"$match": bson.M{
	// 			"size": "M",
	// 		},
	// 	},
	// 	bson.M{
	// 		"$project": bson.M{
	// 			"article.stock": false,
	// 		},
	// 	},
	// }


	// ctx := context.Background()
	// cursor, err := model.collection.Aggregate(ctx, pipeline)

	// if err != nil {
		// return nil, err
	// }
	// defer cursor.Close(ctx)
	stockArticle := defs.StockArticle{}
	// cursor.Next(ctx)

	// if err = cursor.Decode(&stockArticle); err != nil {
	// 	log.Fatal(err)
	// 	return nil, err
	// }
	return &stockArticle, nil
}

func (model *ModelArticle) FindSlice(args *map[string]interface{}) (
	result []defs.Article,
	meta oprs.SliceMetadata,
	err error,
) {
	var bsonData []bson.Raw
	ctx := context.Background()
	filterArgs := NewArticleFiltersFromArgs(args)
	fmt.Printf("%v", filterArgs)
	bsonData, meta, err = oprs.AggregateSlice(model.collection, ctx, combinePipelines(
	bson.A{
		bson.M{
			"$unwind": bson.M{
				"path": "$stock",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}, 
	assertPipeline(filterArgs != nil, bson.A{
		bson.M{
			"$match": filterArgs,
		},
	}), 
	bson.A{
		bson.M{
			"$group": bson.M{
				"_id": "$_id",
				"article": bson.M{"$first":"$$ROOT"},
				"stock": bson.M{ "$push": "$stock"},
			},
	}}, 
	bson.A{
		bson.M{
			"$replaceRoot": bson.M{
				"newRoot": bson.M{
					"$mergeObjects": bson.A{
						"$article",
						bson.M{
							"stock": "$stock",
						},
					},
				},
			},
	}}))


	for _, v := range bsonData {
		article := defs.Article{}
		bson.Unmarshal(v, &article)
		result = append(result, article)
	}
	return result, meta, err
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

type DistincItem struct {
	Name  string `bson:"name"`
	Count int64  `bson:"count"`
}

func (model *ModelArticle) DistincUsingFilters(path interface{}, filter ...interface{}) ([]DistincItem, error) {
	if len(filter) == 0 || filter[0] == nil {
		filter = []interface{}{1}
	}

	pipeline := combinePipelines(createStockEntriesPipeline(path, filter[0]), bson.A{
		bson.M{
			"$project": bson.M{
				"name":  1,
				"count": bson.M{"$size": "$entries"},
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

func (model *ModelArticle) GetCategories(args *map[string]interface{}) ([]DistincItem, error) {
	data, err := model.DistincUsingFilters("$category", NewArticleDistinctFiltersFromArgs(args))

	return data, err
}

func (model *ModelArticle) GetSizes(args *map[string]interface{}) ([]DistincItem, error) {
	data, err := model.DistincUsingFilters("$stock.size", NewArticleDistinctFiltersFromArgs(args))
	return data, err
}

type MaxMinItem = struct {
	Min float64 `bson:"min"`
	Max float64 `bson:"max"`
}

func (model *ModelArticle) GetPriceRange(args *map[string]interface{}) (*MaxMinItem, error) {
	filterArgs := NewArticleFiltersFromArgs(args)
	pipeline := combinePipelines(bson.A{
		bson.M{
			"$unwind": bson.M{
				"path": "$stock",
				"preserveNullAndEmptyArrays": true,
			},
		},
	},
		assertPipeline(filterArgs != nil, bson.A{
			bson.M{
				"$match": filterArgs,
			},
		}),
		bson.A{
			bson.M{
				"$group": bson.M{
					"_id": nil,
					"max": bson.M{"$max": "$price"},
					"min": bson.M{"$min": "$price"},
				},
			},
		})

	ctx := context.Background()
	cursor, err := model.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.Next(ctx)
	data := &MaxMinItem{}
	cursor.Decode(data)
	return data, err
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

func NewArticleDistinctFiltersFromArgs(args *map[string]interface{}, omitted ...string) interface{} {
	conds := bson.A{}
	for argName, argValue := range *args {
		switch argName {
		case "sizes":
			bsonArr := castArrayString(argValue)
			if len(bsonArr) == 0 {
				continue
			}
			conds = append(conds, bson.M{"$in": bson.A{"$$this.stock.size", bsonArr}})
		case "categories":
			bsonArr := castArrayString(argValue)
			if len(bsonArr) == 0 {
				continue
			}
			conds = append(conds, bson.M{"$in": bson.A{"$$this.category", bsonArr}})
		case "priceRange":
			priceRange := castArrayFloat(argValue)
			if len(priceRange) > 0 {
				conds = append(conds, bson.M{"$gt": bson.A{"$$this.price", priceRange[0]}})
			}
			if len(priceRange) > 1 {
				conds = append(conds, bson.M{"$lt": bson.A{"$$this.price", priceRange[1]}})
			}
		}

	}

	if len(conds) == 0 {
		return nil
	}

	return &bson.M{
		"$and": conds,
	}
}

func NewArticleFiltersFromArgs(args *map[string]interface{}) interface{} {
	conds := bson.M{}
	if args != nil {
		for argName, argValue := range *args {
			switch {
			case argName == "categories":
				argName = "category"
				bsonArr := castArrayString(argValue)
				if len(bsonArr) == 0 {
					continue
				}
				argValue = bson.M{
					"$in": argValue,
				}
			case argName == "priceRange":
				argName = "price"
				bsonValue := bson.M{}
				priceRange := castArrayFloat(argValue)
				if len(priceRange) == 0 {
					continue
				}
				if len(priceRange) > 0 {
					bsonValue["$gte"] = priceRange[0]
				}
				if len(priceRange) > 1 {
					bsonValue["$lte"] = priceRange[1]
				}
				argValue = bsonValue
			case argName == "sizes":
				argName = "stock.size"
				bsonArr := castArrayString(argValue)
				if len(bsonArr) == 0 {
					continue
				}
				argValue = bson.M{
					"$in": argValue,
				}
			case argName == "after" || argName == "before" || argName == "first" || argName == "last" || argName == "id":
			default:
				continue
			}
			conds[argName] = argValue

		}
	}

	if len(conds) == 0 {
		return nil
	}

	return &conds
}
