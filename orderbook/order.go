package orderbook

import (
	"time"
	"errors"
)

type OrderId int64
type OrderVersion uint16

const MaxOrderVersion = 65535

type Order struct {
	orderId      OrderId
	creationTime time.Time
	version      OrderVersion
	counterparty string
	side         Side
	price        float32
	volume       float64
}

func NewOrder(counterparty string, side Side, price float32, volume float64) *Order {
	now := time.Now()
	order := new(Order)
	order.orderId = OrderId(now.UnixNano()) //TODO: Create unique id
	order.creationTime = now
	order.version = 1
	order.counterparty = counterparty
	order.side = side
	order.price = price
	order.volume = volume
	return order
}

func (o Order) OrderId() OrderId {
	return o.orderId
}

func (o Order) Counterparty() string {
	return o.counterparty
}

func (o Order) Side() Side {
	return o.side
}

func (o Order) Price() float32 {
	return o.price
}

func (o Order) Volume() float64 {
	return o.volume
}

func (o *Order) AmendPrice(price float32) error {
	if o.version >= MaxOrderVersion {
		return errors.New("Cannot amend price for order: MaxOrderVersion exceeded!")
	} else {
		o.version++	
	}

	o.price = price
	return nil
}

func (o *Order) AmendVolume(volume float64) error {
	if o.version >= MaxOrderVersion {
		return errors.New("Cannot amend volume for order: MaxOrderVersion exceeded!")
	} else {
		o.version++	
	}
	
	o.volume = volume
	return nil
}
