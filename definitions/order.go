package definitions

import (
	"errors"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

const (
	ORDER_STATUS_PENDING   int8 = 0
	ORDER_STATUS_PURCHASED int8 = 1
	ORDER_STATUS_PREPARING int8 = 2
	ORDER_STATUS_SHIPPING  int8 = 3
	ORDER_STATUS_RECEIVED  int8 = 4
	ORDER_STATUS_CANCELLED int8 = 5
)

/*
Order definition
*/
type Order struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Stock     primitive.ObjectID `bson:"stock,omitempty" json:"stock"`
	User      primitive.ObjectID `bson:"user,omitempty" json:"user"`
	State     int8               `bson:"state" json:"state"`
	Notes     string             `bson:"notes" json:"notes"`
	CreatedAt int32              `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt int32              `bson:"updatedAt,omitempty" json:"updatedAt"`
}

func NewOrder(stock primitive.ObjectID, user primitive.ObjectID, notes string) (*Order, error) {
	order := &Order{}
	//TODO if validateInitialValueObjectID stock
	order.Stock = stock
	//TODO if validateInitialValueObjectID user
	order.User = user
	now := int32(time.Now().Unix())
	order.CreatedAt, order.UpdatedAt = now, now
	order.State = ORDER_STATUS_PENDING
	order.Notes = notes
	return order, nil
}

func ValidateNextOrderState(prveState int8, nextState int8) error {
	switch prveState {
	case 0:
		if !(nextState == 1 || nextState == 5) {
			return errors.New("You can not set this state.")
		}
		break
	case 1:
		if !(nextState == 2 || nextState == 5) {

			return errors.New("You can not set this state.")
		}
		break
	case 2:
		if !(nextState == 3 || nextState == 5) {

			return errors.New("You can not set this state.")
		}
		break
	case 3:
		if !(nextState == 4 || nextState == 5) {

			return errors.New("You can not set this state.")
		}
		break

	}
	return nil
}
