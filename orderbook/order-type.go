package orderbook

type OrderType int32

const (
	GOOD_FOR_DAY OrderType = iota
	FILL_AND_KILL
)
