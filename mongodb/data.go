package mongodb

import (
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/mongo"
)

func InitData(db *mongo.Database) {
	model := CreateRepo(db)

	user1, err := model.User.Create(
		"50333339K",
		"Jorge",
		"Lopez Alonso",
		"tuano@tuplacer.com",
		"690876646",
		"Calle de las delicias 69",
	)

	if err != nil {
		log.Fatal(err)
	}

	user2, err := model.User.Create(
		"34546653L",
		"Ruben",
		"Lopez",
		"pizarrin@gmial.com",
		"690876646",
		"Calle de las Mercedes 69",
	)

	if err != nil {
		log.Fatal(err)
	}

	article1, err := model.Article.Create(
		"MiniFalda",
		"Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		7.99,
		[]string{
			"607f1f77bcf86cd799439011",
			"607f1f77bcf86cd799439012"},
		"Faldas",
		2,
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" Article created succesfull : %v\n", article1)

	article2, err := model.Article.Create(
		"Sandalia",
		"Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		17.99,
		[]string{
			"617f1f77bcf86cd799439011",
			"617f1f77bcf86cd799439012"},
		"Zapatos",
		4)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" Article created succesfull : %v\n", article2)

	article3, err := model.Article.Create(
		"Camiseta de tirante",
		"Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		33.99,
		[]string{
			"637f1f77bcf86cd799439011",
			"637f1f77bcf86cd799439012",
			"637f1f77bcf86cd799439013"},
		"Camisetas",
		3)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" Article created succesfull : %v\n", article3)

	article4, err := model.Article.Create(
		"Camiseta de tirante",
		"Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		33.99,
		[]string{
			"637f1f77bcf86cd799439011",
			"637f1f77bcf86cd799439012",
			"637f1f77bcf86cd799439013"},
		"Camisetas",
		4)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" Article created succesfull : %v\n", article4)

	article5, err := model.Article.Create(
		"Camiseta de tirante",
		"Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		33.99,
		[]string{
			"637f1f77bcf86cd799439011",
			"637f1f77bcf86cd799439012",
			"637f1f77bcf86cd799439013"},
		"Camisetas",
		5)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" Article created succesfull : %v\n", article5)

	article6, err := model.Article.Create(
		"Camiseta de tirante",
		"Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		33.99,
		[]string{
			"637f1f77bcf86cd799439011",
			"637f1f77bcf86cd799439012",
			"637f1f77bcf86cd799439013"},
		"Camisetas",
		4)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" Article created succesfull : %v\n", article6)

	article7, err := model.Article.Create(
		"Camiseta de tirante",
		"Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		33.99,
		[]string{
			"637f1f77bcf86cd799439011",
			"637f1f77bcf86cd799439012",
			"637f1f77bcf86cd799439013"},
		"Camisetas",
		3)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" Article created succesfull : %v\n", article7)

	article8, err := model.Article.Create(
		"Camiseta de tirante",
		"Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		33.99,
		[]string{
			"637f1f77bcf86cd799439011",
			"637f1f77bcf86cd799439012",
			"637f1f77bcf86cd799439013"},
		"Camisetas",
		3)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" Article created succesfull : %v\n", article8)

	stock1, err := model.Stock.Create(article1.ID, "S")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" Stock created succesfull : %v\n", stock1)

	stock2, err := model.Stock.Create(article1.ID, "L")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" Stock created succesfull : %v\n", stock2)

	stock3, err := model.Stock.Create(article1.ID, "S")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" Stock created succesfull : %v\n", stock3)

	stock4, err := model.Stock.Create(article1.ID, "L")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" Stock created succesfull : %v\n", stock4)

	stock5, err := model.Stock.Create(article2.ID, "S")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" Stock created succesfull : %v\n", stock5)

	_, err = model.Stock.Create(article3.ID, "M")

	if err != nil {
		log.Fatal(err)
	}

	order1, err := model.Order.Create(
		stock1.ID,
		user1.ID,
		"Solo por las mañanas, trabajo 24/7",
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" Order created succesfull : %v\n", order1)

	_, err = model.Order.Create(
		stock2.ID,
		user2.ID,
		"Solo por las tardes, trabajo la noche",
	)

	if err != nil {
		log.Fatal(err)
	}

	_, err = model.Order.Create(
		stock3.ID,
		user1.ID,
		"Solo por las mañanas, trabajo la noche",
	)

	if err != nil {
		log.Fatal(err)
	}

	model.Order.UpdateState(order1.ID, 5)
}
