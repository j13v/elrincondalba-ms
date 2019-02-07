package definitions

import "github.com/mongodb/mongo-go-driver/bson/primitive"

type StockArticle struct {
	Refs  []primitive.ObjectID `bson:"refs,omitempty" json:"refs,omitempty"`
	Size  string               `bson:"size" json:"size"`
	Count int32                `bson:"count" json:"count"`
}
