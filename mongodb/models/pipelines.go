package models

import (
	defs "github.com/j13v/elrincondalba-ms/definitions"
	"github.com/mongodb/mongo-go-driver/bson"
)

func combinePipeline(args ...bson.A) bson.A {
	pipeline := bson.A{}

	for _, stage := range args {
		pipeline = append(pipeline, stage...)
	}

	return pipeline
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
