package orderbook

import (
	"time"
)

type TradeId int64

type Trade struct {
	tradeId       TradeId
	creationTime  time.Time
	aggressorSide Side
	buyCpty       string
	sellCpty      string
	price         float32
	volume        float64
}

func NewTrade(aggressorSide Side, buyCpty string, sellCpty string, price float32, volume float64) *Trade {
	now := time.Now()
	trade := new(Trade)
	trade.tradeId = TradeId(now.UnixNano()) //TODO: Create unique id
	trade.creationTime = now
	trade.aggressorSide = aggressorSide
	trade.buyCpty = buyCpty
	trade.sellCpty = sellCpty
	trade.price = price
	trade.volume = volume
	return trade
}
