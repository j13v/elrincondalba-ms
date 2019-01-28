package mongodb

import (
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Model struct {
	Article *ModelArticle
	Order   *ModelOrder
}

func CreateModel(db *mongo.Database) Model {
	models := Model{}
	models.Article = NewModelArticle(db)
	models.Order = NewModelOrder(db)
	return models
}
