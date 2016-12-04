package orderbook

import (
	"errors"
	"time"

	"github.com/mdebrouwer/exchange/uuid"
)

type OrderId uuid.UUID
type OrderVersion uint16

const MaxOrderVersion = 65535

type Order struct {
	orderId      OrderId
	creationTime time.Time
	version      OrderVersion
	counterparty string
	side         Side
	price        Price
	volume       Volume
}

func NewOrder(id uuid.UUID, creationTime time.Time, counterparty string, side Side, price Price, volume Volume) Order {
	return Order{
		orderId:      OrderId(id),
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

func (o Order) GetVolume() Volume {
	return o.volume
}

func (o Order) AmendVolume(volume Volume) (Order, error) {
	if o.version >= MaxOrderVersion {
		return o, errors.New("Cannot amend volume for order: MaxOrderVersion exceeded!")
	} else {
		o.version++
	}
	o.volume = volume
	return o, nil
}
