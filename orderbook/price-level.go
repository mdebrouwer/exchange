package orderbook

import (
	"errors"
	"fmt"
	"math"
)

type PriceLevel interface {
	GetPrice() float32
	GetBids() []Order
	GetAsks() []Order
	InsertOrder(order Order) ([]Trade, error)
	AmendOrder(orderId OrderId, volume float64) error
	DeleteOrder(orderId OrderId) error
}

type priceLevel struct {
	price float32
	bids  []Order
	asks  []Order
}

func NewPriceLevel(price float32) *priceLevel {
	pl := new(priceLevel)
	pl.price = price
	pl.bids = make([]Order, 0)
	pl.asks = make([]Order, 0)
	return pl
}

func (pl *priceLevel) GetPrice() float32 {
	return pl.price
}

func (pl *priceLevel) GetBids() []Order {
	return pl.bids
}

func (pl *priceLevel) GetAsks() []Order {
	return pl.asks
}

func (pl *priceLevel) InsertOrder(order Order) ([]Trade, error) {
	if order.price != pl.price {
		return nil, errors.New(fmt.Sprintf("Cannot insert order, invalid price. Price=%v, Order=%+v.", pl.price, order))
	}

	if order.side == BUY {
		if len(pl.asks) > 0 {
			return pl.matchOrders(pl.asks, order), nil
		} else {
			pl.bids = append(pl.bids, order)
		}
	} else if order.side == SELL {
		if len(pl.bids) > 0 {
			return pl.matchOrders(pl.bids, order), nil
		} else {
			pl.asks = append(pl.asks, order)
		}
	}

	return make([]Trade, 0), nil
}

func (pl *priceLevel) AmendOrder(orderId OrderId, volume float64) error {
	index, order, err := pl.findOrder(orderId)
	if err != nil {
		return err
	}

	amendedOrder, err := order.AmendVolume(volume)

	if err != nil {
		return err
	}

	if volume > order.Volume() {
		pl.DeleteOrder(order.orderId)
		pl.InsertOrder(amendedOrder)
	} else {
		if order.side == BUY {
			pl.bids[index] = amendedOrder
		} else if order.side == SELL {
			pl.asks[index] = amendedOrder
		}
	}

	return nil
}

func (pl *priceLevel) DeleteOrder(orderId OrderId) error {
	index, order, err := pl.findOrder(orderId)
	if err != nil {
		return err
	}

	if order.side == BUY {
		pl.bids = append(pl.bids[:index], pl.bids[index+1:]...)
	} else if order.side == SELL {
		pl.asks = append(pl.asks[:index], pl.asks[index+1:]...)
	}

	return nil
}

func (pl *priceLevel) matchOrders(quotes []Order, order Order) []Trade {
	trades := make([]Trade, 0)
	for _, quote := range quotes {
		var buyCpty, sellCpty string
		if order.side == BUY {
			buyCpty = order.Counterparty()
			sellCpty = quote.Counterparty()
		} else if order.side == SELL {
			buyCpty = quote.Counterparty()
			sellCpty = order.Counterparty()
		}

		matchedVolume := math.Min(quote.Volume(), order.Volume())
		trade := NewTrade(order.side, buyCpty, sellCpty, pl.price, matchedVolume)
		trades = append(trades, trade)

		if matchedVolume >= quote.Volume() {
			pl.DeleteOrder(quote.orderId)
		} else {
			quote.AmendVolume(quote.Volume() - matchedVolume)
		}

		if matchedVolume < order.Volume() {
			break
		} else {
			order.AmendVolume(order.Volume() - matchedVolume)
		}
	}

	return trades
}

func (pl *priceLevel) findOrder(orderId OrderId) (index int, order Order, err error) {
	index, order = find(pl.bids, orderId)

	if index == -1 {
		index, order = find(pl.asks, orderId)
	}

	if index == -1 {
		err = errors.New(fmt.Sprintf("OrderId doesn't exist at price level. OrderId=%+v.", orderId))
	}

	return
}

func find(orders []Order, orderId OrderId) (int, Order) {
	for index, order := range orders {
		if order.orderId == orderId {
			return index, order
		}
	}

	return -1, Order{}
}
