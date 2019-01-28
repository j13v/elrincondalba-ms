package definitions

import "github.com/mongodb/mongo-go-driver/bson/primitive"

/*
Article definition
*/
type Article struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Price       float64            `bson:"price" json:"price"`
	Images      []string           `bson:"images" json:"images"`
	Category    string             `bson:"category" json:"category"`
	Rating      int8               `bson:"rating" json:"rating"`
}
