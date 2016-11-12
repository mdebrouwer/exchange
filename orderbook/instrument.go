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

func NewInstrument(name string, tickSize TickSize) *Instrument {
	inst := new(Instrument)
	inst.instrumentId = 0 //TODO: Create unique id
	inst.creationTime = time.Now()
	inst.version = 1
	inst.name = name
	inst.tickSize = tickSize
	return inst
}

func (i Instrument) Name() string {
	return i.name
}

func (i Instrument) TickSize() TickSize {
	return i.tickSize
}
