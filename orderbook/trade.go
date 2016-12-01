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
	price         float64
	volume        float64
}

func NewTrade(aggressorSide Side, buyCpty string, sellCpty string, price float64, volume float64) Trade {
	now := time.Now()
	return Trade{
		tradeId:       TradeId(now.UnixNano()), //TODO: Create unique id
		creationTime:  now,
		aggressorSide: aggressorSide,
		buyCpty:       buyCpty,
		sellCpty:      sellCpty,
		price:         price,
		volume:        volume,
	}
}
