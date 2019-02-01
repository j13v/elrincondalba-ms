package mongodb

import (
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Repo struct {
	Article *ModelArticle
	Order   *ModelOrder
	User    *ModelUser
	Stock   *ModelStock
}

func CreateRepo(db *mongo.Database) Repo {
	models := Repo{}
	models.Article = NewModelArticle(db)
	models.User = NewModelUser(db)
	models.Stock = NewModelStock(db, models.Article)
	models.Order = NewModelOrder(db, models.Stock, models.User)
	return models
}
