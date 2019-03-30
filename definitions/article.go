package definitions

import (
	"errors"
	"strings"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

/*
Article definition
*/
type Article struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string               `bson:"name" json:"name"`
	Description string               `bson:"description" json:"description"`
	Price       float64              `bson:"price" json:"price"`
	Images      []primitive.ObjectID `bson:"images" json:"images"`
	Category    string               `bson:"category" json:"category"`
	Rating      int8                 `bson:"rating" json:"rating"`
	Stock       []Stock              `bson:"stock,omitempty" json:"stock"`
	CreatedAt   int32                `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt   int32                `bson:"updatedAt,omitempty" json:"updatedAt"`
	Disabled    bool                 `bson:"disabled"`
}

/*
NewArticle with validations
*/
func NewArticle(
	name string,
	description string,
	price float64,
	images []primitive.ObjectID,
	category string,
	rating int8,
) (*Article, error) {
	article := &Article{}
	if name == "" {
		return nil, errors.New("Empty name in article creation")
	}
	article.Name = name
	if description == "" {
		return nil, errors.New("Empty description in article creation")
	}
	if len(description) < 50 {
		return nil, errors.New("Invalid length description in actricle creation must have a minimun of 50 characteres")
	}
	article.Description = description
	if price == 0 {
		return nil, errors.New("Empty price in article creation")
	}
	article.Price = price
	if len(images) == 0 {
		return nil, errors.New("Empty images in article creation")
	}
	article.Images = images
	if category == "" {
		return nil, errors.New("Empty category in article creation")
	}
	article.Category = strings.ToLower(category)
	if rating < 0 && rating > 5 {
		return nil, errors.New("Not valid rating in article creation")
	}
	article.Rating = rating
	now := int32(time.Now().Unix())
	article.CreatedAt, article.UpdatedAt = now, now
	return article, nil
}
