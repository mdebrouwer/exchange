package command

import (
	"time"

	ob "github.com/mdebrouwer/exchange/orderbook"
	"github.com/mdebrouwer/exchange/uuid"
)

type Command interface {
	Act(orderbook ob.Orderbook) error
}

type NewOrderCommand struct {
	CounterParty string
	Side         ob.Side
	Price        ob.Price
	Volume       ob.Volume
	Time         time.Time
	UUID         uuid.UUID
}

func (c NewOrderCommand) Act(orderbook ob.Orderbook) error {
	order := ob.NewOrder(c.UUID, c.Time, c.CounterParty, c.Side, c.Price, c.Volume)
	_, err := orderbook.InsertOrder(order)
	return err
}
