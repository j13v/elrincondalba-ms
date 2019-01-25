package models

import (
  "github.com/mongodb/mongo-go-driver/mongo"
)

type DbModel struct {
	Article ModelArticle
	Order ModelOrder
}

func CreateDbModel(db *mongo.Database) DbModel {
	dbmodel := DbModel{}
	dbmodel.Article = GetModelArticle(db)
	dbmodel.Order = GetModelOrder(db)
	return dbmodel
}
