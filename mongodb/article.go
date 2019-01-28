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
	// fmt.Printf("%+v\n", article)
	res, err := model.collection.InsertOne(ctx, val)
	// fmt.Printf("%+v\n", err)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, err
}

func (model *ModelArticle) FindById(id string) (interface{}, error) {
	article := defs.Article{}
	oid, _ := primitive.ObjectIDFromHex(id)
	err := model.collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&article)
	return article, err
}

func (model *ModelArticle) FindSlice(args *FindSliceArguments) ([]interface{}, error) {

	data, err := FindSliceData(model.collection, context.Background(), args)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	articles := []defs.Article{}
	for _, v := range data {
		article := defs.Article{}
		bson.Unmarshal(v, &article)
		// util.FillStruct(&article, v)
		articles = append(articles, article)
		// fmt.Printf("%+v\n", order)
	}

	var interfaceSlice []interface{} = make([]interface{}, len(articles))
	for i, d := range articles {
		interfaceSlice[i] = d
	}

	return interfaceSlice, nil
}

func (model *ModelArticle) GetCount() (int64, error) {
	count, err := GetCount(model.collection, context.Background())
	return count, err
}

// func (model *ModelArticle) FindWithPagination(args map[string]interface{}) ([]interface{}, error) {
// 	conArgs, err := util.NewConnectionArgs(args)
// 	ctx := context.Background()
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}
//
// 	connectionData, err := util.FindWithPagination(ctx, model.collection, bson.M{}, conArgs)
// 	orders := []defs.Article{}
// 	for _, v := range connectionData.Data {
// 		order := defs.Article{}
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

// func (model *ModelArticle) ListArticles() ([]defs.Article, error) {
// 	articles := []defs.Article{}
// 	cursor, err := model.collection.Find(context.Background(), bson.M{})
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}
// 	defer cursor.Close(context.Background())
// 	for cursor.Next(context.Background()) {
// 		article := defs.Article{}
// 		err := cursor.Decode(&article)
// 		if err != nil {
// 			log.Fatal(err)
// 			return nil, err
// 		}
// 		articles = append(articles, article)
// 	}
// 	return articles, nil
// }
