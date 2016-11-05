package orderbook

type OrderBookConfig struct {
	Instrument string
	TickSize   float32
}

func NewOrderBookConfig(instrument string, tickSize float32) *OrderBookConfig {
	cfg := new(OrderBookConfig)
	cfg.Instrument = instrument
	cfg.TickSize = tickSize
	return cfg
}
