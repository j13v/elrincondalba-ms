package mongodb

import (
	defs "github.com/jal88/elrincondalba-ms/definitions"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func InitData(db *mongo.Database) {
	model := CreateModel(db)
	insertedArticleId1, _ := model.Article.Create(&defs.Article{
		Name:        "MiniFalda",
		Description: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		Price:       7.99,
		Images: []string{
			"607f1f77bcf86cd799439011",
			"607f1f77bcf86cd799439012"},
		Category: "Faldas",
		Rating:   2})
	insertedUserId1, _ := model.User.Create(&defs.User{
		DNI:     "50333339K",
		Name:    "Jorge",
		Surname: "Lopez Alonso",
		Email:   "tuano@tuplacer.com",
		Phone:   "690876646",
		Address: "Calle de las delicias 69",
		Notes:   "Solo por las tardes, trabajo la noche",
	})
	model.Order.Create(&defs.Order{
		Article:  insertedArticleId1.(primitive.ObjectID).Hex(),
		User:     insertedUserId1.(primitive.ObjectID).Hex(),
		Size:     "L",
		CreateAt: 1548427228,
		UpdateAt: 1548427228,
		State:    1})

	insertedId2, _ := model.Article.Create(&defs.Article{
		Name:        "Sandalia",
		Description: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		Price:       17.99,
		Images: []string{
			"617f1f77bcf86cd799439011",
			"617f1f77bcf86cd799439012"},
		Category: "Zapatos",
		Rating:   4})
	insertedUserId2, _ := model.User.Create(&defs.User{
		DNI:     "34546653L",
		Name:    "Ruben",
		Surname: "Lopez",
		Email:   "pizarrin@gmial.com",
		Phone:   "690876646",
		Address: "Calle de las Mercedes 69",
		Notes:   "Solo por las mañanas, trabajo 24/7",
	})
	model.Order.Create(&defs.Order{
		Article:  insertedId2.(primitive.ObjectID).Hex(),
		User:     insertedUserId2.(primitive.ObjectID).Hex(),
		Size:     "XL",
		CreateAt: 1548427221,
		UpdateAt: 1548427221,
		State:    2})
	insertedId3, _ := model.Article.Create(&defs.Article{
		Name:        "Camiseta de tirante",
		Description: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		Price:       33.99,
		Images: []string{
			"637f1f77bcf86cd799439011",
			"637f1f77bcf86cd799439012",
			"637f1f77bcf86cd799439013"},
		Category: "Camisetas",
		Rating:   3})
	model.Order.Create(&defs.Order{
		Article:  insertedId3.(primitive.ObjectID).Hex(),
		User:     insertedUserId1.(primitive.ObjectID).Hex(),
		Size:     "S",
		CreateAt: 1548422228,
		UpdateAt: 1548427228,
		State:    3})
}
