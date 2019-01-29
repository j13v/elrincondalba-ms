package definitions

import "github.com/mongodb/mongo-go-driver/bson/primitive"

/*
Order definition
*/
type Order struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Stock     primitive.ObjectID `bson:"stock,omitempty" json:"stock"`
	User      primitive.ObjectID `bson:"user,omitempty" json:"user"`
	State     int8               `bson:"state,omitempty" json:"state"`
	Notes     string             `bson:"notes" json:"notes"`
	CreatedAt int32              `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt int32              `bson:"updatedAt,omitempty" json:"updatedAt"`
}
