package orderbook

import (
	"github.com/mdebrouwer/glog"
)

type Orderbook struct {
	log        *glog.Log
	instrument string
	tickSize   float32
	orderbook  map[float32]*PriceLevel
	bestBid    *PriceLevel
	bestAsk    *PriceLevel
}

func NewOrderbook(log *glog.Log, instrument string, tickSize float32) *Orderbook {
	ob := new(Orderbook)
	ob.log = log
	ob.instrument = instrument
	ob.tickSize = tickSize
	ob.orderbook = make(map[float32]*PriceLevel)
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
