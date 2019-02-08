package definitions

import (
	"errors"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

/*
User definition
*/
type User struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name    string             `bson:"name" json:"name"`
	Surname string             `bson:"surname" json:"surname"`
	Email   string             `bson:"email" json:"email"`
	Phone   string             `bson:"phone" json:"phone"`
	Address string             `bson:"address" json:"address"`
}

func NewUser(name string, surname string, email string, phone string, address string) (*User, error) {
	user := &User{}
	if name == "" {
		return nil, errors.New("Empty name in user creation")
	}
	user.Name = name

	if surname == "" {
		return nil, errors.New("Empty surname in user creation")
	}
	user.Surname = surname

	if email == "" {
		return nil, errors.New("Empty email in user creation")
	}
	user.Email = email

	if phone == "" {
		return nil, errors.New("Empty phone in user creation")
	}
	user.Phone = phone

	if address == "" {
		return nil, errors.New("Empty address in user creation")
	}
	user.Address = address

	return user, nil
}
