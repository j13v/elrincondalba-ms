package mongodb

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/rand"

	"github.com/icrowley/fake"
	defs "github.com/j13v/elrincondalba-ms/definitions"
	"github.com/j13v/elrincondalba-ms/logger"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
)

func InitData(db *mongo.Database) {
	dummyImage := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x96, 0x00, 0x00, 0x00, 0x96, 0x04, 0x03, 0x00, 0x00, 0x00, 0xce, 0x2f, 0x6c,
		0xd1, 0x00, 0x00, 0x00, 0x1b, 0x50, 0x4c, 0x54, 0x45, 0xcc, 0xcc, 0xcc, 0x96, 0x96, 0x96, 0xaa,
		0xaa, 0xaa, 0xb7, 0xb7, 0xb7, 0xc5, 0xc5, 0xc5, 0xbe, 0xbe, 0xbe, 0xb1, 0xb1, 0xb1, 0xa3, 0xa3,
		0xa3, 0x9c, 0x9c, 0x9c, 0x8b, 0x2a, 0x70, 0xc6, 0x00, 0x00, 0x00, 0x09, 0x70, 0x48, 0x59, 0x73,
		0x00, 0x00, 0x0e, 0xc4, 0x00, 0x00, 0x0e, 0xc4, 0x01, 0x95, 0x2b, 0x0e, 0x1b, 0x00, 0x00, 0x01,
		0x00, 0x49, 0x44, 0x41, 0x54, 0x68, 0x81, 0xed, 0xd2, 0x31, 0x6f, 0x83, 0x30, 0x18, 0x84, 0xe1,
		0x8b, 0x31, 0x26, 0xa3, 0x09, 0xb4, 0x73, 0x93, 0x21, 0x33, 0xde, 0x3a, 0xc2, 0xd0, 0x9d, 0xaa,
		0x1d, 0x18, 0x89, 0x84, 0x50, 0xc6, 0x44, 0x91, 0x98, 0x51, 0x52, 0xa9, 0xfd, 0xd9, 0xfd, 0x0c,
		0x69, 0x3b, 0x9b, 0xb5, 0xf7, 0x8c, 0x87, 0xf4, 0x06, 0x1c, 0x03, 0x44, 0x44, 0x44, 0x44, 0x44,
		0x44, 0x44, 0x44, 0x44, 0x44, 0xb4, 0x82, 0x49, 0xd3, 0x16, 0xc7, 0xdb, 0xdf, 0x64, 0xa7, 0xc9,
		0x7c, 0x15, 0xa1, 0xad, 0x2d, 0xd4, 0xd0, 0xd7, 0xd1, 0xe3, 0xa1, 0xfc, 0x59, 0x94, 0x9d, 0x26,
		0xd7, 0xe7, 0x81, 0x29, 0x97, 0x22, 0x91, 0xdf, 0x4f, 0x5a, 0xd5, 0xc2, 0xe4, 0x3a, 0x03, 0xa2,
		0xd4, 0x4e, 0xd3, 0x80, 0x7d, 0x60, 0xeb, 0xc5, 0x21, 0x79, 0x97, 0x2f, 0x2d, 0xa3, 0x11, 0xb8,
		0xaa, 0x13, 0xa0, 0xdf, 0xec, 0x34, 0xe5, 0x78, 0x0e, 0x6c, 0xa1, 0xc2, 0xea, 0xbc, 0x41, 0x55,
		0xc6, 0xd2, 0x72, 0xae, 0x96, 0x45, 0x5b, 0x3f, 0xe9, 0x1c, 0xcd, 0x82, 0x56, 0xd7, 0xb4, 0x95,
		0x2f, 0x20, 0xd9, 0xe0, 0xde, 0xea, 0x9a, 0xce, 0xca, 0xa3, 0xe0, 0x96, 0x1c, 0xd1, 0x38, 0xbf,
		0x97, 0xc9, 0xee, 0x2d, 0x99, 0x3e, 0x16, 0xbd, 0x17, 0x10, 0xdb, 0xf9, 0xbc, 0xd6, 0x9f, 0xbf,
		0xad, 0xd8, 0x2e, 0x3a, 0x2f, 0x07, 0x73, 0x92, 0xff, 0xf1, 0x09, 0x38, 0xbe, 0x16, 0x73, 0xcb,
		0x4f, 0x03, 0x76, 0xe1, 0xad, 0xf5, 0xe5, 0x50, 0xc8, 0xfd, 0xaa, 0xa1, 0x33, 0x35, 0xce, 0x2d,
		0x3f, 0xb9, 0xfe, 0x21, 0xbc, 0x15, 0x9f, 0xe5, 0xce, 0xfb, 0x7b, 0xaf, 0x46, 0x7f, 0xbf, 0x7c,
		0xcb, 0x4f, 0x0b, 0xee, 0x3d, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0xfd, 0x3f,
		0xdf, 0x59, 0xa5, 0x25, 0x51, 0x8a, 0xf0, 0x7f, 0xae, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e,
		0x44, 0xae, 0x42, 0x60, 0x82,
	}
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
			logger.WithFields(logrus.Fields{
				"err": err.Error(),
			}).Warn("Failed to create user")
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
			[]defs.File{
				defs.File{
					Filename: "dummy",
					File:     bytes.NewReader(dummyImage),
					Size:     int64(len(dummyImage)),
				},
			},
			categories[rand.Intn(len(categories))],
			int8(rand.Intn(5)),
		)

		if err != nil {
			logger.WithFields(logrus.Fields{
				"err": err.Error(),
			}).Warn("Failed to create article")
		} else {
			logger.WithFields(logrus.Fields{
				"id":      article.ID.Hex(),
				"article": fmt.Sprintf("%s, %s", article.Name, article.Category),
			}).Info("Created article")

			articles = append(articles, *article)
		}
	}

	for i := 0; i <= len(articles)*3; i++ {
		article := articles[rand.Intn(len(articles))]
		stockItem, err := repo.Article.AddStock(article.ID, sizes[rand.Intn(len(sizes))])

		if err != nil {
			logger.WithFields(logrus.Fields{
				"err": err.Error(),
			}).Warn("Failed to create stock")
		} else {
			logger.WithFields(logrus.Fields{
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
			logger.WithFields(logrus.Fields{
				"err": err.Error(),
			}).Warn("Failed to create order")
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
				logger.WithFields(logrus.Fields{
					"err": err.Error(),
				}).Warn("Failed to update order")
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
