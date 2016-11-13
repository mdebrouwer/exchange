package orderbook

import (
	"github.com/mdebrouwer/exchange/log"
)

type Orderbook interface {
	InsertOrder(order Order) ([]Trade, error)
	AmendOrder(order Order) error
	DeleteOrder(order Order) error
	GetTopLevel() (PriceLevel, PriceLevel)
	GetBestBid() PriceLevel
	GetBestAsk() PriceLevel
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
	return nil, nil
}

func (ob *orderbook) AmendOrder(order Order) error {
	return nil
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
