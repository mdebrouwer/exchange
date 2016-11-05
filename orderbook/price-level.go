package orderbook

import (
	"errors"
	"fmt"
)

type PriceLevel struct {
	price float32
	bids  []*Order
	asks  []*Order
}

func NewPriceLevel(price float32) *PriceLevel {
	pl := new(PriceLevel)
	pl.price = price
	pl.bids = make([]*Order, 0, 1)
	pl.asks = make([]*Order, 0, 1)
	return pl
}

func (pl *PriceLevel) GetPrice() float32 {
	return pl.price
}

func (pl *PriceLevel) GetBids() []*Order {
	return pl.bids
}

func (pl *PriceLevel) GetAsks() []*Order {
	return pl.asks
}

func (pl *PriceLevel) InsertOrder(order *Order) error {
	if order.price != pl.price {
		return errors.New(fmt.Sprintf("Cannot insert order, invalid price. Price=%v, Order=%+v.", pl.price, order))
	}

	if order.side == BUY {
		pl.bids = append(pl.bids, order)
	} else if order.side == SELL {
		pl.asks = append(pl.asks, order)
	}

	return nil
}

func (pl *PriceLevel) AmendOrder(order *Order, volume float32) error {
	if order.price != pl.price {
		return errors.New(fmt.Sprintf("Cannot amend order, invalid price. Price=%v, Order=%+v.", pl.price, order))
	}

	_, order, err := pl.findOrder(order.side, order.orderId)
	if err != nil {
		return err
	}

	order.volume = volume
	return nil
}

func (pl *PriceLevel) DeleteOrder(order *Order) error {
	if order.price != pl.price {
		return errors.New(fmt.Sprintf("Cannot delete order, invalid price. Price=%v, Order=%+v.", pl.price, order))
	}

	index, order, err := pl.findOrder(order.side, order.orderId)
	if err != nil {
		return err
	}

	if order.side == BUY {
		pl.bids = append(pl.bids[:index], pl.bids[index+1:]...)
	} else if order.side == SELL {
		pl.bids = append(pl.asks[:index], pl.asks[index+1:]...)
	}

	return nil
}

func (pl *PriceLevel) findOrder(side Side, orderId int64) (int, *Order, error) {
	var order *Order
	var index int
	if side == BUY {
		index, order = find(pl.bids, orderId)
	} else if order.side == SELL {
		index, order = find(pl.asks, orderId)
	}

	if order != nil {
		return index, order, nil
	}

	return index, order, errors.New(fmt.Sprintf("Order doesn't exist at price level. Order=%+v.", order))
}

func find(orders []*Order, orderId int64) (int, *Order) {
	for index, order := range orders {
		if order.orderId == orderId {
			return index, order
		}
	}

	return -1, nil
}
