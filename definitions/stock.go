package definitions

import (
	"errors"
	"strings"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

/*
Stock definition
*/
type Stock struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Article   primitive.ObjectID `bson:"article,omitempty" json:"article"`
	Size      string             `bson:"size" json:"size"`
	CreatedAt int32              `bson:"createdAt" json:"createdAt"`
}

func NewStock(article primitive.ObjectID, size string) (*Stock, error) {
	if size == "" {
		return nil, errors.New("Empty size in stock creation")
	}
	return &Stock{
		Article:   article,
		Size:      strings.ToUpper(size),
		CreatedAt: int32(time.Now().Unix()),
	}, nil
}

func (s *Stock) SetID(id primitive.ObjectID) {
	s.ID = id
}
