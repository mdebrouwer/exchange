package orderbook_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"log"

	ob "github.com/mdebrouwer/exchange/orderbook"
)

var _ = Describe("OrderBook", func() {
	var orderbook ob.Orderbook

	BeforeEach(func() {
		instrument := ob.NewInstrument("TEST_INSTRUMENT", 10)
		orderbook = ob.NewOrderbook(log.New(ioutil.Discard, "", 0), *instrument)
	})

	Describe("Inserting a new Order to empty Orderbook", func() {
		BeforeEach(func() {
			orderbook.InsertOrder(ob.NewOrder("CPTY1", ob.BUY, 100, 1))
		})
		Context("If side is Bid", func() {
			It("should be added to the Orderbook and available from GetBestBid", func() {
				Expect(orderbook.GetBestBid() != nil)
			})
			// It("should be added at the correct price level", func() {
			// 	Expect(orderbook.GetBestBid().GetPrice() == 100)
			// })
		})
	})
})
