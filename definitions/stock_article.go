package definitions

import "github.com/mongodb/mongo-go-driver/bson/primitive"

type StockArticle struct {
	Refs  []primitive.ObjectID `bson:"refs,omitempty" json:"refs,omitempty"`
	Size  string               `bson:"size" json:"size"`
	Count int32                `bson:"count" json:"count"`
}

type StockOrderArticleItem struct {
	ID    primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
	Order primitive.ObjectID `bson:"order,omitempty" json:"order,omitempty"`
	State int8               `bson:"state,omitempty" json:"state,omitempty"`
}

type StockOrderArticle struct {
	Refs  []StockOrderArticleItem `bson:"refs,omitempty" json:"refs,omitempty"`
	Size  string                  `bson:"size" json:"size"`
	Count int32                   `bson:"count" json:"count"`
}
