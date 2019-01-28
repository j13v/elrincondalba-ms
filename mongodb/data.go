package mongodb

import (
	defs "github.com/jal88/elrincondalba-ms/definitions"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func InitArticlesData(db *mongo.Database) {
	article1 := &defs.Article{
		Name:        "MiniFalda",
		Description: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		Price:       7.99,
		Images: []string{
			"607f1f77bcf86cd799439011",
			"607f1f77bcf86cd799439012"},
		Category: "Faldas",
		Rating:   2}
	article2 := &defs.Article{
		Name:        "Sandalia",
		Description: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		Price:       17.99,
		Images: []string{
			"617f1f77bcf86cd799439011",
			"617f1f77bcf86cd799439012"},
		Category: "Zapatos",
		Rating:   4}
	article3 := &defs.Article{
		Name:        "Camiseta de tirante",
		Description: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		Price:       33.99,
		Images: []string{
			"637f1f77bcf86cd799439011",
			"637f1f77bcf86cd799439012",
			"637f1f77bcf86cd799439013"},
		Category: "Camisetas",
		Rating:   3}

	model := CreateModel(db)
	model.Article.Create(article1)
	model.Article.Create(article2)
	model.Article.Create(article3)
}

func InitOrdersData(db *mongo.Database) {
	order1 := &defs.Order{
		Article:  "5c4b58c67cbc327aa78383fd",
		User:     "707f1f77bcf86cd799439011",
		Size:     "L",
		CreateAt: 1548427228,
		UpdateAt: 1548427228,
		State:    1}
	order2 := &defs.Order{
		Article:  "5c4b58c67cbc327aa78383fd",
		User:     "707f1f77bcf86cd799439012",
		Size:     "XL",
		CreateAt: 1548427221,
		UpdateAt: 1548427221,
		State:    2}
	order3 := &defs.Order{
		Article:  "5c4b58c67cbc327aa78383fd",
		User:     "707f1f77bcf86cd799439013",
		Size:     "S",
		CreateAt: 1548422228,
		UpdateAt: 1548427228,
		State:    3}

	model := CreateModel(db)
	model.Order.Create(order1)
	model.Order.Create(order2)
	model.Order.Create(order3)
}

func InitData(db *mongo.Database) {
	InitArticlesData(db)
	InitOrdersData(db)
}
