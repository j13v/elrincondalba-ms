package mongodb

import (
	"github.com/j13v/elrincondalba-ms/mongodb/models"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Repo struct {
	Article *models.ModelArticle
	Order   *models.ModelOrder
	User    *models.ModelUser
	Stock   *models.ModelStock
}

func CreateRepo(db *mongo.Database) Repo {
	repModels := Repo{}
	repModels.Article = models.NewModelArticle(db)
	repModels.User = models.NewModelUser(db)
	repModels.Stock = models.NewModelStock(db, repModels.Article)
	repModels.Order = models.NewModelOrder(db, repModels.Stock, repModels.User)
	return repModels
}
