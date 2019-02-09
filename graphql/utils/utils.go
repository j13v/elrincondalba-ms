package utils

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"github.com/j13v/elrincondalba-ms/mongodb/helpers"
)

func HexToBase64(strHex string) (string, error) {
	src := []byte(strHex)

	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, src)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(dst), nil
}

type ConnectionSliceMetadata struct {
	Total      int `json:"total"`
	SliceCount int `json:"sliceCount"`
	SliceStart int `json:"sliceStart"`
}

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
	TotalCount int64    `json:"totalCount"`
	Edges      []*Edge  `json:"edges"`
	PageInfo   PageInfo `json:"pageInfo"`
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
	meta interface{},
) *Connection {
	byteData, _ := json.Marshal(meta)
	metaSlice := ConnectionSliceMetadata{}
	json.Unmarshal(byteData, &metaSlice)
	conn := NewConnection()
	sliceEnd := metaSlice.SliceStart + len(arraySlice)
	beforeOffset := relay.GetOffsetWithDefault(args.Before, metaSlice.Total)
	afterOffset := relay.GetOffsetWithDefault(args.After, -1)

	startOffset := ternaryMax(metaSlice.SliceStart-1, afterOffset, -1) + 1
	endOffset := ternaryMin(sliceEnd, beforeOffset, metaSlice.Total)

	if args.First != -1 {
		endOffset = min(endOffset, startOffset+args.First)
	}

	if args.Last != -1 {
		startOffset = max(startOffset, endOffset-args.Last)
	}

	begin := max(startOffset-metaSlice.SliceStart, 0)
	end := len(arraySlice) - (sliceEnd - endOffset)

	if begin > end {
		return conn
	}

	slice := arraySlice[begin:end]

	edges := []*Edge{}
	for _, value := range slice {
		oid, _ := helpers.GetObjectIDFromValue(value)
		hid, _ := HexToBase64(helpers.GetHexFromObjectID(oid))
		edges = append(edges, &Edge{
			Cursor: hid,
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

	upperBound := metaSlice.Total
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

	conn.TotalCount = int64(metaSlice.Total)
	conn.Edges = edges
	conn.PageInfo = PageInfo{
		StartCursor:     firstEdgeCursor,
		EndCursor:       lastEdgeCursor,
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
	}

	return conn
}

func GetValueByJSONTag(value interface{}, tagName string) (interface{}, bool) {
	elem := reflect.ValueOf(value)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	for i := 0; i < elem.NumField(); i++ {
		elemField := elem.Field(i)
		typeField := elem.Type().Field(i)
		if tagArgs := strings.Split(typeField.Tag.Get("json"), ","); tagArgs[0] == tagName {
			return elemField.Interface(), true
		}

	}
	return nil, false
}

type ConnectionConfig struct {
	Name             string          `json:"name"`
	NodeType         *graphql.Object `json:"nodeType"`
	EdgeFields       graphql.Fields  `json:"edgeFields"`
	ConnectionFields graphql.Fields  `json:"connectionFields"`
}

var pageInfoType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "PageInfo",
	Description: "Information about pagination in a connection.",
	Fields: graphql.Fields{
		"hasNextPage": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "When paginating forwards, are there more items?",
		},
		"hasPreviousPage": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "When paginating backwards, are there more items?",
		},
		"startCursor": &graphql.Field{
			Type:        graphql.String,
			Description: "When paginating backwards, the cursor to continue.",
		},
		"endCursor": &graphql.Field{
			Type:        graphql.String,
			Description: "When paginating forwards, the cursor to continue.",
		},
	},
})

func ConnectionDefinitions(config ConnectionConfig) *graphql.Object {

	edgeType := graphql.NewObject(graphql.ObjectConfig{
		Name:        config.Name + "Edge",
		Description: "An edge in a connection",
		Fields: graphql.Fields{
			"node": &graphql.Field{
				Type:        config.NodeType,
				Description: "The item at the end of the edge",
			},
			"cursor": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: " cursor for use in pagination",
			},
		},
	})
	for fieldName, fieldConfig := range config.EdgeFields {
		edgeType.AddFieldConfig(fieldName, fieldConfig)
	}

	connectionType := graphql.NewObject(graphql.ObjectConfig{
		Name:        config.Name + "Connection",
		Description: "A connection to a list of items.",

		Fields: graphql.Fields{
			"totalCount": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Total records.",
			},
			"pageInfo": &graphql.Field{
				Type:        graphql.NewNonNull(pageInfoType),
				Description: "Information to aid in pagination.",
			},
			"edges": &graphql.Field{
				Type:        graphql.NewList(edgeType),
				Description: "Information to aid in pagination.",
			},
		},
	})
	for fieldName, fieldConfig := range config.ConnectionFields {
		connectionType.AddFieldConfig(fieldName, fieldConfig)
	}

	return connectionType
}

func NewIdArgs(id interface{}) *map[string]interface{} {
	return &map[string]interface{}{"id": id}
}
