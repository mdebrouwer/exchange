package orderbook

import (
	"time"
)

type TickSize float32

type Instrument struct {
	instrumentId int64
	creationTime time.Time
	version      int32
	name         string
	tickSize     TickSize
}

func NewInstrument(name string, tickSize TickSize) Instrument {
	return Instrument{
		instrumentId: 0, //TODO: Create unique id
		creationTime: time.Now(),
		version: 1,
		name: name,
		tickSize: tickSize,
	}
}

func (i Instrument) Name() string {
	return i.name
}

func (i Instrument) TickSize() TickSize {
	return i.tickSize
}
