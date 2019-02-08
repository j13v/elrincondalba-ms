package mongodb

import (
	"context"
	"fmt"
	"math"
	"math/rand"

	"github.com/icrowley/fake"
	defs "github.com/jal88/elrincondalba-ms/definitions"
	"github.com/jal88/elrincondalba-ms/logger"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
)

func InitData(db *mongo.Database) {
	categories := []string{
		"VESTIDOS",
		"MONOS",
		"FALDAS",
		"TRAJES",
		"BOLSOS",
		"PANTALONES",
		"CAMISAS",
		"ACCESORIOS",
		"ZAPATOS",
		"CHAQUETAS",
		"CARDIGANS",
		"BLUSAS",
		"CAMISETAS",
		"JERSEYS",
		"CINTURONES",
		"ABANICOS",
		"SOMBREROS",
		"RELOJES",
		"BUFANDAS",
	}
	sizes := []string{
		"XS",
		"S",
		"M",
		"L",
		"XL",
		"XXL",
	}
	logger := logger.NewLogger("initData")
	repo := CreateRepo(db)
	users := []defs.User{}
	articles := []defs.Article{}
	stock := []defs.Stock{}
	db.Drop(context.Background())
	for i := 0; i <= 10; i++ {
		user, err := repo.User.Create(
			fake.FirstName(),
			fake.LastName(),
			fake.EmailAddress(),
			fake.Phone(),
			fake.StreetAddress(),
		)

		if err != nil {
			logger.Warn("Failed to create user %s", fmt.Sprintf("%s, %s", user.Name, user.Surname), err)
		} else {
			logger.WithFields(logrus.Fields{
				"id":   user.ID.Hex(),
				"user": fmt.Sprintf("%s, %s", user.Name, user.Surname),
			}).Info("Created user")

			users = append(users, *user)
		}
	}

	for i := 0; i <= 150; i++ {
		article, err := repo.Article.Create(
			fake.ProductName(),
			fake.Sentences(),
			math.Floor(rand.Float64()*10000)/100,
			[]string{"adsffasdasd"},
			categories[rand.Intn(len(categories))],
			int8(rand.Intn(5)),
		)
		if err != nil {
			logger.Warn("Failed to create article %s", err)
		} else {
			logger.WithFields(logrus.Fields{
				"id":      article.ID.Hex(),
				"article": fmt.Sprintf("%s, %s", article.Name, article.Category),
			}).Info("Created article")

			articles = append(articles, *article)
		}
	}

	for i := 0; i <= len(articles)*30; i++ {
		article := articles[rand.Intn(len(articles))]
		stockItem, err := repo.Stock.Create(article.ID, sizes[rand.Intn(len(sizes))])

		if err != nil {
			logger.Warn("Failed to create stock %s", err)
		} else {
			logger.WithFields(logrus.Fields{
				"id":      stockItem.ID.Hex(),
				"size":    stockItem.Size,
				"article": article.ID,
			}).Info("Created stock item")

			stock = append(stock, *stockItem)
		}
	}

	for i := 0; i <= len(articles)*10; i++ {
		article := articles[rand.Intn(len(articles))]
		stockItem, err := repo.Stock.Create(article.ID, sizes[rand.Intn(len(sizes))])

		if err != nil {
			logger.Warn("Failed to create stock %s", err)
		} else {
			logger.WithFields(logrus.Fields{
				"id":      stockItem.ID.Hex(),
				"size":    stockItem.Size,
				"article": article.ID,
			}).Info("Created stock item")

			stock = append(stock, *stockItem)
		}
	}

	for i := 0; i <= len(articles); i++ {
		stockItem := stock[rand.Intn(len(stock))]
		user := users[rand.Intn(len(users))]
		order, err := repo.Order.Create(
			stockItem.ID,
			user.ID,
			"Solo por las tardes, trabajo la noche",
		)

		if err != nil {
			logger.Warn("Failed to create order %s", err)
		} else {
			logger.WithFields(logrus.Fields{
				"id":   order.ID.Hex(),
				"size": stockItem.Size,
				"user": user.Name,
			}).Info("Created order item")
		}

		if rand.Intn(10) > 8 {
			err := order.UpdateState(5)
			if err != nil {
				logger.Warn("Failed to update order %s", err)
			} else {
				logger.WithFields(logrus.Fields{
					"id":   order.ID.Hex(),
					"size": stockItem.Size,
					"user": user.Name,
				}).Info("Updated order to cancelled")
			}
		}
	}

	// _, err = repo.Order.Create(
	// 	stock2.ID,
	// 	user2.ID,
	// 	"Solo por las tardes, trabajo la noche",
	// )

	// article1, err := repo.Article.Create(
	// 	"MiniFalda",
	// 	"Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
	// 	7.99,
	// 	[]string{
	// 		"607f1f77bcf86cd799439011",
	// 		"607f1f77bcf86cd799439012"},
	// 	"Faldas",
	// 	2,
	// )
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Printf(" Article created succesfull : %v\n", article1)
	//
	// article2, err := repo.Article.Create(
	// 	"Sandalia",
	// 	"Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
	// 	17.99,
	// 	[]string{
	// 		"617f1f77bcf86cd799439011",
	// 		"617f1f77bcf86cd799439012"},
	// 	"Zapatos",
	// 	4)
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Printf(" Article created succesfull : %v\n", article2)
	//
	// article3, err := repo.Article.Create(
	// 	"Camiseta de tirante",
	// 	"Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
	// 	33.99,
	// 	[]string{
	// 		"637f1f77bcf86cd799439011",
	// 		"637f1f77bcf86cd799439012",
	// 		"637f1f77bcf86cd799439013"},
	// 	"Camisetas",
	// 	3)
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Printf(" Article created succesfull : %v\n", article3)
	//
	// article4, err := repo.Article.Create(
	// 	"Camiseta de tirante",
	// 	"Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
	// 	33.99,
	// 	[]string{
	// 		"637f1f77bcf86cd799439011",
	// 		"637f1f77bcf86cd799439012",
	// 		"637f1f77bcf86cd799439013"},
	// 	"Camisetas",
	// 	4)
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Printf(" Article created succesfull : %v\n", article4)
	//
	// article5, err := repo.Article.Create(
	// 	"Camiseta de tirante",
	// 	"Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
	// 	33.99,
	// 	[]string{
	// 		"637f1f77bcf86cd799439011",
	// 		"637f1f77bcf86cd799439012",
	// 		"637f1f77bcf86cd799439013"},
	// 	"Camisetas",
	// 	5)
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Printf(" Article created succesfull : %v\n", article5)
	//
	// article6, err := repo.Article.Create(
	// 	"Camiseta de tirante",
	// 	"Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
	// 	33.99,
	// 	[]string{
	// 		"637f1f77bcf86cd799439011",
	// 		"637f1f77bcf86cd799439012",
	// 		"637f1f77bcf86cd799439013"},
	// 	"Camisetas",
	// 	4)
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Printf(" Article created succesfull : %v\n", article6)
	//
	// article7, err := repo.Article.Create(
	// 	"Camiseta de tirante",
	// 	"Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
	// 	33.99,
	// 	[]string{
	// 		"637f1f77bcf86cd799439011",
	// 		"637f1f77bcf86cd799439012",
	// 		"637f1f77bcf86cd799439013"},
	// 	"Camisetas",
	// 	3)
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Printf(" Article created succesfull : %v\n", article7)
	//
	// article8, err := repo.Article.Create(
	// 	"Camiseta de tirante",
	// 	"Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
	// 	33.99,
	// 	[]string{
	// 		"637f1f77bcf86cd799439011",
	// 		"637f1f77bcf86cd799439012",
	// 		"637f1f77bcf86cd799439013"},
	// 	"Camisetas",
	// 	3)
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Printf(" Article created succesfull : %v\n", article8)
	//
	// stock1, err := repo.Stock.Create(article1.ID, "S")
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Printf(" Stock created succesfull : %v\n", stock1)
	//
	// stock2, err := repo.Stock.Create(article1.ID, "L")
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Printf(" Stock created succesfull : %v\n", stock2)
	//
	// stock3, err := repo.Stock.Create(article1.ID, "S")
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Printf(" Stock created succesfull : %v\n", stock3)
	//
	// stock4, err := repo.Stock.Create(article1.ID, "L")
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Printf(" Stock created succesfull : %v\n", stock4)
	//
	// stock5, err := repo.Stock.Create(article2.ID, "S")
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Printf(" Stock created succesfull : %v\n", stock5)
	//
	// _, err = repo.Stock.Create(article3.ID, "M")
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// order1, err := repo.Order.Create(
	// 	stock1.ID,
	// 	user1.ID,
	// 	"Solo por las mañanas, trabajo 24/7",
	// )
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Printf(" Order created succesfull : %v\n", order1)
	//
	// _, err = repo.Order.Create(
	// 	stock2.ID,
	// 	user2.ID,
	// 	"Solo por las tardes, trabajo la noche",
	// )
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// _, err = repo.Order.Create(
	// 	stock3.ID,
	// 	user1.ID,
	// 	"Solo por las mañanas, trabajo la noche",
	// )
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// repo.Order.UpdateState(order1.ID, 5)
}
