package orderbook

import (
	"errors"
	"fmt"
	"math"
)

type PriceLevel struct {
	price float32
	bids  []*Order
	asks  []*Order
}

func NewPriceLevel(price float32) *PriceLevel {
	pl := new(PriceLevel)
	pl.price = price
	pl.bids = make([]*Order, 0)
	pl.asks = make([]*Order, 0)
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

func (pl *PriceLevel) InsertOrder(order *Order) ([]*Trade, error) {
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

	return make([]*Trade, 0), nil
}

func (pl *PriceLevel) AmendOrder(amendOrder Order, volume float64) error {
	if amendOrder.Price() != pl.price {
		return errors.New(fmt.Sprintf("Cannot amend order, invalid price. Price=%v, Order=%+v.", pl.price, amendOrder))
	}

	if volume == amendOrder.Volume() {
		return errors.New(fmt.Sprintf("Cannot amend order, invalid volume. Price=%v, Order=%+v.", pl.price, amendOrder))
	}

	_, order, err := pl.findOrder(amendOrder.orderId)
	if err != nil {
		return err
	}

	if volume > order.Volume() {
		pl.DeleteOrder(order.orderId)
		err := order.AmendVolume(volume)
		pl.InsertOrder(order)
		return err
	} else {
		return order.AmendVolume(volume)
	}
}

func (pl *PriceLevel) DeleteOrder(orderId OrderId) error {
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

func (pl *PriceLevel) matchOrders(quotes []*Order, order *Order) []*Trade {
	trades := make([]*Trade, 0, 1)
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

func (pl *PriceLevel) findOrder(orderId OrderId) (int, *Order, error) {
	var index int
	var order *Order
	
	if len(pl.bids) > 0 {
		index, order = find(pl.bids, orderId)
	} else if len(pl.asks) > 0 {
		index, order = find(pl.asks, orderId)
	}

	if order != nil {
		return index, order, nil
	}

	return -1, nil, errors.New(fmt.Sprintf("OrderId doesn't exist at price level. OrderId=%+v.", orderId))
}

func find(orders []*Order, orderId OrderId) (int, *Order) {
	for index, order := range orders {
		if order.orderId == orderId {
			return index, order
		}
	}

	return -1, nil
}
