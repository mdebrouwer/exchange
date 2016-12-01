package orderbook

import (
	"errors"
	"fmt"
	"math"
)

type PriceLevel interface {
	GetPrice() float64
	GetBids() []Order
	GetAsks() []Order
	InsertOrder(order Order) error
	DeleteOrder(orderId OrderId) error
	MatchOrders(order Order) ([]Trade, error)
}

type priceLevel struct {
	price float64
	bids  []Order
	asks  []Order
}

func NewPriceLevel(price float64) *priceLevel {
	pl := new(priceLevel)
	pl.price = price
	pl.bids = make([]Order, 0)
	pl.asks = make([]Order, 0)
	return pl
}

func (pl *priceLevel) GetPrice() float64 {
	return pl.price
}

func (pl *priceLevel) GetBids() []Order {
	return pl.bids
}

func (pl *priceLevel) GetAsks() []Order {
	return pl.asks
}

func (pl *priceLevel) InsertOrder(order Order) error {
	if order.GetPrice() != pl.price {
		return errors.New(fmt.Sprintf("Cannot insert order, invalid price. Price=%v, Order=%+v.", pl.price, order))
	}

	if order.GetSide() == BUY {
		if len(pl.asks) > 0 {
			return errors.New(fmt.Sprintf("Cannot insert order, in cross. Order=%+v.", order))
		} else {
			pl.bids = append(pl.bids, order)
		}
	} else if order.GetSide() == SELL {
		if len(pl.bids) > 0 {
			return errors.New(fmt.Sprintf("Cannot insert order, in cross. Order=%+v.", order))
		} else {
			pl.asks = append(pl.asks, order)
		}
	}

	return nil
}

func (pl *priceLevel) MatchOrders(order Order) ([]Trade, error) {
	var quotes []Order
	if order.GetSide() == BUY && len(pl.asks) > 0 {
			return nil, errors.New(fmt.Sprintf("Cannot match buy order. Order=%+v.", order))
	} else {
		quotes = pl.bids
	}
	
	if order.GetSide() == SELL && len(pl.bids) > 0 {
		return nil, errors.New(fmt.Sprintf("Cannot match sell order. Order=%+v.", order))
	} else {
		quotes = pl.asks
	}
	
	trades := make([]Trade, 0)
	for _, quote := range quotes {
		var buyCpty, sellCpty string
		if order.side == BUY {
			buyCpty = order.GetCounterparty()
			sellCpty = quote.GetCounterparty()
		} else if order.GetSide() == SELL {
			buyCpty = quote.GetCounterparty()
			sellCpty = order.GetCounterparty()
		}

		matchedVolume := math.Min(quote.GetVolume(), order.GetVolume())
		trade := NewTrade(order.side, buyCpty, sellCpty, pl.price, matchedVolume)
		trades = append(trades, trade)

		if matchedVolume >= quote.GetVolume() {
			pl.DeleteOrder(quote.orderId)
		} else {
			quote.AmendVolume(quote.GetVolume() - matchedVolume)
		}

		if matchedVolume < order.GetVolume() {
			break
		} else {
			order.AmendVolume(order.GetVolume() - matchedVolume)
		}
	}

	return trades, nil
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
