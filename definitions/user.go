package definitions

import "github.com/mongodb/mongo-go-driver/bson/primitive"

/*
User definition
*/
type User struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	DNI     string             `bson:"dni" json:"dni"`
	Name    string             `bson:"name" json:"name"`
	Surname string             `bson:"surname" json:"surname"`
	Email   string             `bson:"email" json:"email"`
	Phone   string             `bson:"phone" json:"phone"`
	Address string             `bson:"address" json:"address"`
	Notes   string             `bson:"notes" json:"notes"`
}
