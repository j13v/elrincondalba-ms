package utils

import (
	"reflect"
	"strings"

	"github.com/graphql-go/relay"
	"github.com/jal88/elrincondalba-ms/mongodb"
)

type Edge struct {
	Node   interface{} `json:"node"`
	Cursor string      `json:"cursor"`
}

type PageInfo struct {
	StartCursor     string `json:"startCursor"`
	EndCursor       string `json:"endCursor"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
	HasNextPage     bool   `json:"hasNextPage"`
}

type Connection struct {
	Edges    []*Edge  `json:"edges"`
	PageInfo PageInfo `json:"pageInfo"`
}

func NewConnection() *Connection {
	return &Connection{
		Edges:    []*Edge{},
		PageInfo: PageInfo{},
	}
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func ternaryMax(a, b, c int) int {
	return max(max(a, b), c)
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func ternaryMin(a, b, c int) int {
	return min(min(a, b), c)
}

func ConnectionFromArraySlice(
	arraySlice []interface{},
	args relay.ConnectionArguments,
	meta relay.ArraySliceMetaInfo,
) *Connection {
	conn := NewConnection()
	sliceEnd := meta.SliceStart + len(arraySlice)
	beforeOffset := relay.GetOffsetWithDefault(args.Before, meta.ArrayLength)
	afterOffset := relay.GetOffsetWithDefault(args.After, -1)

	startOffset := ternaryMax(meta.SliceStart-1, afterOffset, -1) + 1
	endOffset := ternaryMin(sliceEnd, beforeOffset, meta.ArrayLength)

	if args.First != -1 {
		endOffset = min(endOffset, startOffset+args.First)
	}

	if args.Last != -1 {
		startOffset = max(startOffset, endOffset-args.Last)
	}

	begin := max(startOffset-meta.SliceStart, 0)
	end := len(arraySlice) - (sliceEnd - endOffset)

	if begin > end {
		return conn
	}

	slice := arraySlice[begin:end]

	edges := []*Edge{}
	for _, value := range slice {
		edges = append(edges, &Edge{
			Cursor: mongodb.GetHexFromObjectID(mongodb.GetObjectIDFromValue(value)),
			Node:   value,
		})
	}

	var firstEdgeCursor, lastEdgeCursor string
	if len(edges) > 0 {
		firstEdgeCursor = edges[0].Cursor
		lastEdgeCursor = edges[len(edges)-1:][0].Cursor
	}

	lowerBound := 0
	if len(args.After) > 0 {
		lowerBound = afterOffset + 1
	}

	upperBound := meta.ArrayLength
	if len(args.Before) > 0 {
		upperBound = beforeOffset
	}

	hasPreviousPage := false
	if args.Last != -1 {
		hasPreviousPage = startOffset > lowerBound
	}

	hasNextPage := false
	if args.First != -1 {
		hasNextPage = endOffset < upperBound
	}

	conn.Edges = edges
	conn.PageInfo = PageInfo{
		StartCursor:     firstEdgeCursor,
		EndCursor:       lastEdgeCursor,
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
	}

	return conn
}

func GetValueByJSONTag(value interface{}, tagName string) (string, bool) {
	elem := reflect.ValueOf(value)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	for i := 0; i < elem.NumField(); i++ {
		elemField := elem.Field(i)
		typeField := elem.Type().Field(i)
		if tagArgs := strings.Split(typeField.Tag.Get("json"), ","); tagArgs[0] == tagName {
			if str, ok := elemField.Interface().(string); ok {
				return str, true
			}
		}

	}
	return "", false
}
