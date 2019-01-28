package definitions

import "github.com/mongodb/mongo-go-driver/bson/primitive"

/*
Order definition
*/
type Order struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Article  string             `bson:"article" json:"article"`
	User     string             `bson:"user" json:"user"`
	Size     string             `bson:"size" json:"size"`
	CreateAt int32              `bson:"createAt" json:"createAt"`
	UpdateAt int32              `bson:"updateAt" json:"updateAt"`
	State    int8               `bson:"state" json:"state"`
}
