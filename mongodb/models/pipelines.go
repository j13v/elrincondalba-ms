package models

import (
	defs "github.com/j13v/elrincondalba-ms/definitions"
	"github.com/mongodb/mongo-go-driver/bson"
)

func combineBsonArrays(args ...bson.A) bson.A {
	res := bson.A{}
	for _, stage := range args {
		res = append(res, stage...)
	}
	return res
}

func combinePipelines(args ...bson.A) bson.A {
	return combineBsonArrays(args...)
}

func assertPipeline(assertion bool, pipeline bson.A) bson.A {
	if assertion {
		return pipeline
	}
	return bson.A{}
}

func assertDocument(assertion bool, doc bson.M) bson.M {
	if assertion {
		return doc
	}
	return bson.M{}
}

func combineDocuments(args ...bson.M) bson.A {
	res := bson.A{}
	for _, stage := range args {
		res = append(res, stage)
	}
	return res
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
