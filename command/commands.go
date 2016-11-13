package command

import (
	ob "github.com/mdebrouwer/exchange/orderbook"
)

type Command interface {
	Act(orderbook ob.Orderbook) error
}

type NewOrderCommand struct {
	CounterParty string
	Side         ob.Side
	Price        float32
	Volume       float64
}

func (c NewOrderCommand) Act(orderbook ob.Orderbook) error {
	order := ob.NewOrder(c.CounterParty, c.Side, c.Price, c.Volume)
	_, err := orderbook.InsertOrder(order)
	return err
}
