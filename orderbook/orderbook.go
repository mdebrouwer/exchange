package orderbook

import (
	"github.com/mdebrouwer/exchange/log"
)

type Orderbook interface {
	InsertOrder(order Order) ([]Trade, error)
	DeleteOrder(order Order) error
	GetTopLevel() (PriceLevel, PriceLevel)
	GetBestBid() PriceLevel
	GetBestAsk() PriceLevel
	GetPriceLevels() []PriceLevel
}

type orderbook struct {
	logger     log.Logger
	instrument Instrument
	orderbook  map[TickSize]*priceLevel
	bestBid    *priceLevel
	bestAsk    *priceLevel
}

func NewOrderbook(logger log.Logger, instrument Instrument) Orderbook {
	ob := new(orderbook)
	ob.logger = logger
	ob.instrument = instrument
	ob.orderbook = make(map[TickSize]*priceLevel)
	return ob
}

func (ob *orderbook) InsertOrder(order Order) ([]Trade, error) {
	ob.logger.Printf("Inserting order for: %s\n", order.counterparty)
	return nil, nil
}

func (ob *orderbook) DeleteOrder(order Order) error {
	return nil
}

func (ob *orderbook) GetTopLevel() (PriceLevel, PriceLevel) {
	return ob.bestBid, ob.bestAsk
}

func (ob *orderbook) GetBestBid() PriceLevel {
	return ob.bestBid
}

func (ob *orderbook) GetBestAsk() PriceLevel {
	return ob.bestAsk
}

func (ob *orderbook) GetPriceLevels() []PriceLevel {
	return make([]PriceLevel, 0)
}
