package mongodb

import (
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Model struct {
	Article *ModelArticle
	Order   *ModelOrder
	User    *ModelUser
	Stock   *ModelStock
}

func CreateModel(db *mongo.Database) Model {
	models := Model{}
	models.Article = NewModelArticle(db)
	models.User = NewModelUser(db)
	models.Stock = NewModelStock(db)
	models.Order = NewModelOrder(db, models.Stock, models.User)
	return models
}
