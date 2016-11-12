package orderbook

import (
	"github.com/mdebrouwer/exchange/log"
)

type Orderbook struct {
	logger     log.Logger
	instrument Instrument
	orderbook  map[TickSize]*PriceLevel
	bestBid    *PriceLevel
	bestAsk    *PriceLevel
}

func NewOrderbook(logger log.Logger, instrument Instrument) *Orderbook {
	ob := new(Orderbook)
	ob.logger = logger
	ob.instrument = instrument
	ob.orderbook = make(map[TickSize]*PriceLevel)
	return ob
}

func (ob *Orderbook) InsertOrder(order Order) ([]*Trade, error) {
	return nil, nil
}

func (ob *Orderbook) AmendOrder(order Order) error {
	return nil
}

func (ob *Orderbook) DeleteOrder(order Order) error {
	return nil
}

func (ob *Orderbook) GetTopLevel() (*PriceLevel, *PriceLevel) {
	return ob.bestBid, ob.bestAsk
}

func (ob *Orderbook) GetBestBid() *PriceLevel {
	return ob.bestBid
}

func (ob *Orderbook) GetBestAsk() *PriceLevel {
	return ob.bestAsk
}
