package definitions

import "errors"

/*
Stock definition
*/
type Stock struct {
	Size     string `bson:"size" json:"size"`
	Count    int    `bson:"count" json:"count"`
	CreateAt int32  `bson:"createAt" json:"createAt"`
}

// TODO Improve error definition
func NewStock(size string, count int) (*Stock, error) {
	stock := &Stock{}
	if size == "" {
		return nil, errors.New("Empty size")
	}
	if count == 0 {
		return nil, errors.New("Count must be greater or lower than 0")
	}
	return stock, nil
}
