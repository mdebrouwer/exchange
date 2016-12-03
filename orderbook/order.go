package orderbook

import (
	"errors"
	"time"
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
	price        Price
	volume       float64
}

func NewOrder(creationTime time.Time, counterparty string, side Side, price Price, volume float64) Order {
	return Order{
		orderId:      OrderId(creationTime.UnixNano()), //TODO: Create unique id
		creationTime: creationTime,
		version:      1,
		counterparty: counterparty,
		side:         side,
		price:        price,
		volume:       volume,
	}
}

func (o Order) GetOrderId() OrderId {
	return o.orderId
}

func (o Order) GetCreationTime() time.Time {
	return o.creationTime
}

func (o Order) GetVersion() OrderVersion {
	return o.version
}

func (o Order) GetCounterparty() string {
	return o.counterparty
}

func (o Order) GetSide() Side {
	return o.side
}

func (o Order) GetPrice() Price {
	return o.price
}

func (o Order) GetVolume() float64 {
	return o.volume
}

func (o Order) AmendVolume(volume float64) (Order, error) {
	if o.version >= MaxOrderVersion {
		return o, errors.New("Cannot amend volume for order: MaxOrderVersion exceeded!")
	} else {
		o.version++
	}
	o.volume = volume
	return o, nil
}
