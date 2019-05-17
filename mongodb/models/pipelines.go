package models

import (
	defs "github.com/j13v/elrincondalba-ms/definitions"
	"github.com/mongodb/mongo-go-driver/bson"
)

const (
	ARTICLE_STOCK_SORTING_POPULARITY int = 1
	ARTICLE_STOCK_SORTING_PURCHASES  int = 2
	ARTICLE_STOCK_SORTING_RECENTS    int = 3
)

func arrayBsonMToBsonA(arr []bson.M) (out bson.A) {
	for _, item := range arr {
		out = append(out, item)
	}
	return
}

func composePipeline(stages ...bson.M) bson.A {
	return arrayBsonMToBsonA(stages)
}

func combinePipelines(pipelines ...bson.A) (out bson.A) {
	for _, pipeline := range pipelines {
		for _, stage := range pipeline {
			out = append(out, stage)
		}
	}
	return
}

func assertPipeline(assertion bool, pipeline bson.A) bson.A {
	if assertion {
		return pipeline
	}
	return bson.A{}
}

func assertPipelineStage(assertion bool, doc bson.M) bson.M {
	if assertion {
		return doc
	}
	return bson.M{}
}

var pipelineStockOrder = bson.A{
	bson.M{
		"$lookup": bson.M{
			"from":         "order",
			"localField":   "_id",
			"foreignField": "stock",
			"as":           "order",
		},
	},
	bson.M{
		"$addFields": bson.M{
			"order": bson.M{"$slice": bson.A{"$order", -1}},
		},
	},
	bson.M{
		"$unwind": bson.M{
			"path": "$order",
			"preserveNullAndEmptyArrays": true,
		},
	},
	bson.M{
		"$project": bson.M{
			"item":      1,
			"article":   1,
			"size":      1,
			"createdAt": 1,
			"order":     "$order._id",
			"state": bson.M{
				"$ifNull": bson.A{"$order.state", -1},
			},
		},
	},
}

var pipelineStockOrderAvailable = bson.A{
	bson.M{
		"$match": bson.M{
			"$or": bson.A{
				bson.M{"state": bson.M{"$eq": -1}},
				bson.M{"state": bson.M{"$eq": defs.ORDER_STATUS_CANCELLED}},
			},
		},
	},
}

var pipelineStockArticleGroup = bson.A{
	bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"article": "$article",
				"size":    "$size",
			},
			"count": bson.M{"$sum": 1},
			"refs":  bson.M{"$push": "$_id"},
		},
	},
	bson.M{
		"$project": bson.M{
			"_id":     0,
			"article": "$_id.article",
			"size":    "$_id.size",
			"refs":    1,
			"count":   1,
		},
	},
}

var pipelineStockOrderArticleGroup = bson.A{
	bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"article": "$article",
				"size":    "$size",
			},
			"count": bson.M{"$sum": 1},
			"refs":  bson.M{"$push": bson.M{"id": "$_id", "order": "$order", "state": "$state"}},
		},
	},
	bson.M{
		"$project": bson.M{
			"_id":     0,
			"article": "$_id.article",
			"size":    "$_id.size",
			"refs":    1,
			"count":   1,
		},
	},
}

/*
pipelineArticleStockUnwinder pipeline fragment to unwind stock articles
*/
var pipelineArticleStockUnwinder = bson.A{
	bson.M{
		"$unwind": bson.M{
			"path": "$stock",
			"preserveNullAndEmptyArrays": true,
		},
	},
}

/*
pipelineArticleStockRewinder pipeline fragment to rewind previous unwinded stock articles
*/
var pipelineArticleStockRewinder = bson.A{
	bson.M{
		"$group": bson.M{
			"_id": "$_id",
			"article": bson.M{
				"$first": "$$ROOT",
			},
			"stock": bson.M{
				"$push": bson.M{
					"$cond": bson.A{
						bson.M{
							"$ne": bson.A{"$stock", bson.TypeUndefined},
						},
						bson.M{
							"$mergeObjects": bson.A{
								"$stock",
								bson.M{
									"orders":    "$orders",
									"available": "$available",
								},
							},
						},
						bson.TypeNull,
					},
				},
			},
			"purchases": bson.M{
				"$sum": bson.M{
					"$cond": bson.A{
						bson.M{
							// if last order status is received then stock is purchased
							"$in": bson.A{
								bson.M{
									"$arrayElemAt": bson.A{"$orders.state", -1},
								},
								bson.A{
									defs.ORDER_STATUS_RECEIVED,
								},
							},
						},
						1,
						0,
					},
				},
			},
			"available": bson.M{
				"$addToSet": "$available",
			},
		},
	},
	bson.M{
		"$project": bson.M{
			"article":   1,
			"purchases": "$purchases",
			"stock": bson.M{
				// Prevent invalid stock items of null value
				"$reduce": bson.M{
					"input":        "$stock",
					"initialValue": "$stock.0",
					"in": bson.M{
						"$cond": bson.A{
							bson.M{
								"$ne": bson.A{"$$value", bson.TypeNull},
							},
							bson.M{
								"$concatArrays": bson.A{"$$value", bson.A{"$$this"}},
							},
							bson.A{},
						},
					},
				},
			},
			// Check if any stock is available
			"available": bson.M{
				"$anyElementTrue": bson.A{"$available"},
			},
		},
	},
	bson.M{
		"$replaceRoot": bson.M{
			"newRoot": bson.M{
				"$mergeObjects": bson.A{
					"$article",
					bson.M{
						"stock":     "$stock",
						"purchases": "$purchases",
					},
				},
			},
		},
	},
	bson.M{
		// Hide partial article orders
		"$project": bson.M{
			"orders": 0,
		},
	},
}

/*
pipelineArticleStockOrder Add flag available stock item when order null or
valid state to notify when stock is ready to be purchased
*/
var pipelineArticleStockOrder = bson.A{
	bson.M{
		"$lookup": bson.M{
			"from":         "order",
			"localField":   "stock._id",
			"foreignField": "stock",
			"as":           "orders",
		},
	},
	bson.M{
		"$addFields": bson.M{
			"stock.available": bson.M{
				"$cond": bson.A{
					bson.M{
						"$and": bson.A{
							bson.M{
								"$ne": bson.A{
									"$stock._id", bson.TypeUndefined,
								},
							},
							bson.M{
								"$or": bson.A{
									bson.M{
										"$eq": bson.A{
											bson.M{
												"$size": "$orders",
											},
											0,
										},
									},
									bson.M{
										"$in": bson.A{
											bson.M{
												"$arrayElemAt": bson.A{"$orders.state", -1},
											},
											bson.A{
												defs.ORDER_STATUS_CANCELLED,
											},
										},
									},
								},
							},
						},
					},
					true, false,
				},
			},
		},
	},
}

func createArticleStockPipeline(filters *map[string]interface{}, sorting int) bson.A {
	// fmt.Println(createArticleStockFiltersPipeline(filters))
	return combinePipelines(
		pipelineArticleStockUnwinder,
		// Apply user filters or extra stages
		createArticleStockFiltersPipeline(filters),
		pipelineArticleStockOrder,
		pipelineArticleStockRewinder,
		//createArticleStockFlagsPipeline(),
		createArticleStockSortingPipeline(sorting),
	)
}

func createArticleStockFlagsPipeline() bson.A {
	return bson.A{
		bson.M{
			"$match": bson.M{
				// Show only available articles, this will also hide when stock filters are apply
				"available": bson.M{"$eq": true},
				// Do not show disabled articles
				"disabled": bson.M{"$eq": false},
			},
		},
	}
}

func createArticleStockSortingPipeline(sorting int) bson.A {
	switch sorting {
	case ARTICLE_STOCK_SORTING_POPULARITY:
		return bson.A{
			bson.M{
				"$sort": bson.M{
					// Order by popularity
					"rating": -1,
				},
			},
		}
	case ARTICLE_STOCK_SORTING_PURCHASES:
		return bson.A{
			bson.M{
				"$sort": bson.M{
					// Order by purchases
					"purchases": -1,
				},
			},
		}
	case ARTICLE_STOCK_SORTING_RECENTS:
		return bson.A{
			bson.M{
				"$sort": bson.M{
					// Order by recents
					"createdAt": -1,
				},
			},
		}
	default:
		return bson.A{}
	}
}

func createArticleStockFiltersPipeline(args *map[string]interface{}) bson.A {
	criterias := bson.M{}
	empty := true
	if args != nil {
		for argName, argValue := range *args {
			switch {
			case argName == "sizes":
				bsonArr := castArrayString(argValue)
				if len(bsonArr) == 0 {
					continue
				}
				criterias["stock.size"] = bson.M{
					"$in": argValue,
				}
				empty = false
			case argName == "categories":
				bsonArr := castArrayString(argValue)
				if len(bsonArr) == 0 {
					continue
				}
				criterias["category"] = bson.M{
					"$in": argValue,
				}
				empty = false
			case argName == "priceRange":
				bsonValue := bson.M{}
				priceRange := castArrayFloat(argValue)
				if len(priceRange) == 0 {
					continue
				}
				if len(priceRange) > 0 {
					bsonValue["$gte"] = priceRange[0]
				}
				if len(priceRange) > 1 {
					bsonValue["$lte"] = priceRange[1]
				}
				criterias["price"] = bsonValue
				empty = false
			}
		}
	}

	if empty {
		return bson.A{}
	}
	return bson.A{bson.M{
		"$match": criterias,
	}}
}

// func createArticleStockRewinderPipeline(stages ...bson.M) bson.A {
// 	return combinePipelines(composePipeline(
// 		bson.M{
// 			"$group": bson.M{
// 				"_id":     "$_id",
// 				"article": bson.M{"$first": "$$ROOT"},
// 				"stock":   bson.M{"$push": "$stock"},
// 			},
// 		},
// 		bson.M{
// 			"$replaceRoot": bson.M{
// 				"newRoot": bson.M{
// 					"$mergeObjects": bson.A{
// 						"$article",
// 						bson.M{
// 							"stock": "$stock",
// 						},
// 					},
// 				},
// 			},
// 		}, composePipeline(stages...),
// 	)
// }

func createStockEntriesPipeline(path interface{}, filter interface{}) bson.A {
	return bson.A{
		bson.M{
			"$unwind": bson.M{
				"path": "$stock",
				"preserveNullAndEmptyArrays": true,
			},
		},
		// Group entries by path
		bson.M{
			"$group": bson.M{
				"_id":     path,
				"entries": bson.M{"$push": "$$ROOT"},
			},
		},
		// Apply projection filter
		bson.M{
			"$project": bson.M{
				"_id":  0,
				"name": "$_id",
				"entries": bson.M{
					"$filter": bson.M{
						"input": "$entries",
						"cond":  filter,
					},
				},
			},
		},
		// Omit null
		bson.M{
			"$match": bson.M{
				"name": bson.M{"$ne": nil},
			},
		},
	}
}
