package orderbook

import (
	"time"
)

type Order struct {
	orderId      int64
	creationTime time.Time
	counterparty string
	side         Side
	price        float32
	volume       float32
}

func NewOrder(orderId int64, counterparty string, side Side, price float32, volume float32) *Order {
	order := new(Order)
	order.orderId = orderId
	order.creationTime = time.Now()
	order.counterparty = counterparty
	order.side = side
	order.price = price
	order.volume = volume
	return order
}
