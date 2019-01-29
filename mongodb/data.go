package mongodb

import (
	defs "github.com/jal88/elrincondalba-ms/definitions"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func InitData(db *mongo.Database) {
	model := CreateModel(db)

	user1, _ := model.User.Create(&defs.User{
		DNI:     "50333339K",
		Name:    "Jorge",
		Surname: "Lopez Alonso",
		Email:   "tuano@tuplacer.com",
		Phone:   "690876646",
		Address: "Calle de las delicias 69",
	})

	user2, _ := model.User.Create(&defs.User{
		DNI:     "34546653L",
		Name:    "Ruben",
		Surname: "Lopez",
		Email:   "pizarrin@gmial.com",
		Phone:   "690876646",
		Address: "Calle de las Mercedes 69",
	})

	article1, _ := model.Article.Create(&defs.Article{
		Name:        "MiniFalda",
		Description: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		Price:       7.99,
		Images: []string{
			"607f1f77bcf86cd799439011",
			"607f1f77bcf86cd799439012"},
		Category: "Faldas",
		Rating:   2,
	})

	article2, _ := model.Article.Create(&defs.Article{
		Name:        "Sandalia",
		Description: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		Price:       17.99,
		Images: []string{
			"617f1f77bcf86cd799439011",
			"617f1f77bcf86cd799439012"},
		Category: "Zapatos",
		Rating:   4})

	article3, _ := model.Article.Create(&defs.Article{
		Name:        "Camiseta de tirante",
		Description: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		Price:       33.99,
		Images: []string{
			"637f1f77bcf86cd799439011",
			"637f1f77bcf86cd799439012",
			"637f1f77bcf86cd799439013"},
		Category: "Camisetas",
		Rating:   3})

	stock1, _ := model.Stock.Create(article1.ID, "S")

	stock2, _ := model.Stock.Create(article1.ID, "L")

	stock3, _ := model.Stock.Create(article2.ID, "S")

	model.Stock.Create(article3.ID, "M")

	model.Order.Create(&defs.Order{
		User:  user1.ID,
		Stock: stock1.ID,
		Notes: "Solo por las mañanas, trabajo 24/7",
	})

	model.Order.Create(&defs.Order{
		User:  user2.ID,
		Stock: stock2.ID,
		Notes: "Solo por las tardes, trabajo la noche",
	})

	model.Order.Create(&defs.Order{
		User:  user1.ID,
		Stock: stock3.ID,
		Notes: "Solo por las mañanas, trabajo la noche",
	})
}
