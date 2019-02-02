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

const (
	ORDER_PAYMENT_METHOD_CASH          int8 = 0
	ORDER_PAYMENT_METHOD_BANK_TRANSFER int8 = 1
	ORDER_PAYMENT_METHOD_CREDIT_CARD   int8 = 2
	ORDER_PAYMENT_METHOD_PAYPAL        int8 = 3
)

/*
Order definition
*/
type Order struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Stock         primitive.ObjectID `bson:"stock,omitempty" json:"stock"`
	User          primitive.ObjectID `bson:"user,omitempty" json:"user"`
	State         int8               `bson:"state" json:"state"`
	Notes         string             `bson:"notes" json:"notes"`
	PaymentMethod int8               `bson:"paymentMethod" json:"paymentMethod"`
	PurchaseRef   string             `bson:"purchaseRef" json:"purchaseRef"`
	TrackingRef   string             `bson:"trackingRef" json:"trackingRef"`
	CreatedAt     int32              `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt     int32              `bson:"updatedAt,omitempty" json:"updatedAt"`
}

func (order *Order) Purchase(paymentMethod int8, purchaseRef string) error {
	_, err := PurchaseOrder(order, paymentMethod, purchaseRef)
	return err
}

func (order *Order) Prepare() error {
	_, err := PrepareOrder(order)
	return err
}

func (order *Order) Ship(trackingRef string) error {
	_, err := ShipOrder(order, trackingRef)
	return err
}

func (order *Order) ConfirmReceived() error {
	_, err := ConfirmReceivedOrder(order)
	return err
}

func (order *Order) Cancel() error {
	_, err := CancelOrder(order)
	return err
}

func (order *Order) UpdateState(state int8) error {
	_, err := UpdateStateOrder(order, state)
	return err
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

func PurchaseOrder(order *Order, paymentMethod int8, purchaseRef string) (*Order, error) {
	if _, err := UpdateStateOrder(order, ORDER_STATUS_PURCHASED); err != nil {
		return nil, err
	}
	// TODO Validate payment methods
	order.PaymentMethod = paymentMethod
	order.PurchaseRef = purchaseRef
	return order, nil
}

func PrepareOrder(order *Order) (*Order, error) {
	_, err := UpdateStateOrder(order, ORDER_STATUS_PREPARING)
	return order, err
}

func ShipOrder(order *Order, trackingRef string) (*Order, error) {
	if _, err := UpdateStateOrder(order, ORDER_STATUS_SHIPPING); err != nil {
		return nil, err
	}
	order.TrackingRef = trackingRef
	return order, nil
}

func ConfirmReceivedOrder(order *Order) (*Order, error) {
	_, err := UpdateStateOrder(order, ORDER_STATUS_RECEIVED)
	return order, err
}

func CancelOrder(order *Order) (*Order, error) {
	_, err := UpdateStateOrder(order, ORDER_STATUS_CANCELLED)
	return order, err
}

func UpdateStateOrder(order *Order, state int8) (*Order, error) {
	if err := ValidateNextOrderState(order.State, state); err != nil {
		return nil, err
	}

	order.State = state
	order.UpdatedAt = int32(time.Now().Unix())
	return order, nil
}

func ValidateNextOrderState(prveState int8, nextState int8) error {
	switch prveState {
	case ORDER_STATUS_PENDING:
		if !(nextState == ORDER_STATUS_PURCHASED || nextState == ORDER_STATUS_CANCELLED) {
			return errors.New("You can not set this state valid states (ORDER_STATUS_PURCHASED, ORDER_STATUS_CANCELLED).")
		}
		break
	case ORDER_STATUS_PURCHASED:
		if !(nextState == ORDER_STATUS_PREPARING || nextState == ORDER_STATUS_CANCELLED) {

			return errors.New("You can not set this state valid states (ORDER_STATUS_PREPARING, ORDER_STATUS_CANCELLED)")
		}
		break
	case ORDER_STATUS_PREPARING:
		if !(nextState == ORDER_STATUS_SHIPPING || nextState == ORDER_STATUS_CANCELLED) {

			return errors.New("You can not set this state valid states (ORDER_STATUS_SHIPPING, ORDER_STATUS_CANCELLED)")
		}
		break
	case ORDER_STATUS_SHIPPING:
		if !(nextState == ORDER_STATUS_RECEIVED || nextState == ORDER_STATUS_CANCELLED) {

			return errors.New("You can not set this state valid states (ORDER_STATUS_RECEIVED, ORDER_STATUS_CANCELLED)")
		}
		break

	}
	return nil
}
