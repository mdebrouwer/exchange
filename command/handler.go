package command

import (
	ob "github.com/mdebrouwer/exchange/orderbook"
)

type Handler interface {
	Handle(c Command) error
	Stop()
}

type request struct {
	command Command
	err     chan error
}

type handler struct {
	requests chan request
	stop     chan struct{}
}

func (h *handler) Handle(c Command) error {
	errChan := make(chan error)
	req := request{
		command: c,
		err:     errChan,
	}
	h.requests <- req
	return <-errChan
}

func (h *handler) Stop() {
	h.stop <- struct{}{}
}

func NewHandler(orderbook ob.Orderbook) Handler {
	requests := make(chan request)
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case req := <-requests:
				req.err <- req.command.Act(orderbook)
			case <-stop:
				return
			}
		}
	}()

	return &handler{
		requests: requests,
		stop:     stop,
	}
}
