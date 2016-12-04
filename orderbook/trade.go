package orderbook

import (
	"time"
)

type Trade struct {
	creationTime  time.Time
	aggressorSide Side
	buyCpty       string
	sellCpty      string
	price         Price
	volume        Volume
}

func NewTrade(creationTime time.Time, aggressorSide Side, buyCpty string, sellCpty string, price Price, volume Volume) Trade {
	return Trade{
		creationTime:  creationTime,
		aggressorSide: aggressorSide,
		buyCpty:       buyCpty,
		sellCpty:      sellCpty,
		price:         price,
		volume:        volume,
	}
}

func (t Trade) GetCreationTime() time.Time {
	return t.creationTime
}

func (t Trade) GetPrice() Price {
	return t.price
}

func (t Trade) GetVolume() Volume {
	return t.volume
}
