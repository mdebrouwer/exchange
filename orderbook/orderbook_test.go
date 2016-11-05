package orderbook_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	ob "github.com/mdebrouwer/exchange/orderbook"
	"github.com/mdebrouwer/glog"
)

var _ = Describe("OrderBook", func() {
	var (
		orderbook *ob.Orderbook
	)

	BeforeEach(func() {
		var log = glog.NewLogger("test.log")
		var cfg = ob.NewOrderBookConfig("TEST_INSTRUMENT", 10)
		orderbook = ob.NewOrderbook(log, cfg.Instrument, cfg.TickSize)
	})

	Describe("Inserting a new Order to empty Orderbook", func() {
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
