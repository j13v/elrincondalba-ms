package models

import (
	"context"
  "log"
	"github.com/mongodb/mongo-go-driver/bson"
  "github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"time"
)

type TypeArticle struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string   `bson:"name" json:"name"`
	Description string   `bson:"description" json:"description"`
	Price       float64  `bson:"price" json:"price"`
	Images      []string `bson:"images" json:"images"`
	Category    string   `bson:"category" json:"category"`
	Rating      int8     `bson:"rating" json:"rating"`
}

type ModelArticle struct {
	collection *mongo.Collection
}

func GetModelArticle(db *mongo.Database) ModelArticle {
	return ModelArticle{collection: db.Collection("article")}
}

func (model *ModelArticle) CreateArticle(article TypeArticle) (interface{}, interface{}) {
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

func (model *ModelArticle) GetArticleById(id string) (TypeArticle, error) {
	article := TypeArticle{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
  oid, _ := primitive.ObjectIDFromHex(id)
	err := model.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&article)
	return article, err
}

func (model *ModelArticle) ListArticles() ([]TypeArticle, error) {
  articles := []TypeArticle{}
	cursor, err := model.collection.Find(context.Background(), bson.M{})
  if err != nil {
    log.Fatal(err)
    return nil, err
  }
  defer cursor.Close(context.Background())
  for cursor.Next(context.Background()) {
    article := TypeArticle{}
    err := cursor.Decode(&article)
    if err != nil {
      log.Fatal(err)
      return nil, err;
    }
    articles = append(articles, article)
  }
	return articles, nil
}
