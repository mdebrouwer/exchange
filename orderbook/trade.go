package orderbook

import (
	"time"
)

type Trade struct {
	trade_id       int64
	creation_time  time.Time
	aggressor_side Side
	buy_cpty       string
	sell_cpty      string
	price          float32
	volume         float32
}

func NewTrade(aggressor_side Side, buy_cpty string, sell_cpty string, price float32, volume float32) *Trade {
	trade := new(Trade)
	trade.creation_time = time.Now()
	trade.aggressor_side = aggressor_side
	trade.buy_cpty = buy_cpty
	trade.sell_cpty = sell_cpty
	trade.price = price
	trade.volume = volume
	return trade
}
