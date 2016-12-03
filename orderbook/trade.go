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
	price         Price
	volume        float64
}

func NewTrade(creationTime time.Time, aggressorSide Side, buyCpty string, sellCpty string, price Price, volume float64) Trade {
	return Trade{
		tradeId:       TradeId(creationTime.UnixNano()), //TODO: Create unique id
		creationTime:  creationTime,
		aggressorSide: aggressorSide,
		buyCpty:       buyCpty,
		sellCpty:      sellCpty,
		price:         price,
		volume:        volume,
	}
}

func (t Trade) GetTradeId() TradeId {
	return t.tradeId
}

func (t Trade) GetCreationTime() time.Time {
	return t.creationTime
}

func (t Trade) GetPrice() Price {
	return t.price
}

func (t Trade) GetVolume() float64 {
	return t.volume
}
