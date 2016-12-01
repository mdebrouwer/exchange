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
		orderbook = ob.NewOrderbook(log.New(ioutil.Discard, "", 0), instrument)
	})

	Describe("Inserting a new Order to empty Orderbook", func() {
		Context("If side is Buy", func() {
			BeforeEach(func() {
				orderbook.InsertOrder(ob.NewOrder("CPTY1", ob.BUY, 100, 1))
			})
			It("should be added to the Orderbook and available from GetBestBid", func() {
				Expect(orderbook.GetBestBid()).ShouldNot(BeNil())
				Expect(orderbook.GetPriceLevels()).To(HaveLen(1))
				Expect(orderbook.GetBestBid().GetBids()).To(HaveLen(1))
				Expect(orderbook.GetBestBid().GetAsks()).To(HaveLen(0))
			})
			It("should not be available from GetBestAsk", func() {
				Expect(orderbook.GetBestAsk()).Should(BeNil())
			})
			It("should be added at the correct price level", func() {
				Expect(orderbook.GetBestBid().GetPrice() == 100)
			})
		})
		Context("If side is Sell", func() {
			BeforeEach(func() {
				orderbook.InsertOrder(ob.NewOrder("CPTY1", ob.SELL, 100, 1))
			})
			It("should be added to the Orderbook and available from GetBestAsk", func() {
				Expect(orderbook.GetBestAsk()).ShouldNot(BeNil())
				Expect(orderbook.GetPriceLevels()).To(HaveLen(1))
				Expect(orderbook.GetBestAsk().GetAsks()).To(HaveLen(1))
				Expect(orderbook.GetBestAsk().GetBids()).To(HaveLen(0))
			})
			It("should not be available from GetBestBid", func() {
				Expect(orderbook.GetBestBid()).Should(BeNil())
			})
			It("should be added at the correct price level", func() {
				Expect(orderbook.GetBestAsk().GetPrice() == 100)
			})
		})
	})
	Describe("Inserting a new Order to Orderbook at existing pricelevel", func() {
		Context("If side is Buy", func() {
			BeforeEach(func() {
				orderbook.InsertOrder(ob.NewOrder("CPTY1", ob.BUY, 100, 1))
				orderbook.InsertOrder(ob.NewOrder("CPTY2", ob.BUY, 100, 1))
			})
			It("should be added to the Orderbook and available from GetBestBid", func() {
				Expect(orderbook.GetBestBid().GetBids()).To(HaveLen(2))
				Expect(orderbook.GetBestBid().GetAsks()).To(HaveLen(0))
			})
		})
		Context("If side is Sell", func() {
			BeforeEach(func() {
				orderbook.InsertOrder(ob.NewOrder("CPTY1", ob.SELL, 100, 1))
				orderbook.InsertOrder(ob.NewOrder("CPTY2", ob.SELL, 100, 1))
			})
			It("should be added to the Orderbook and available from GetBestAsk", func() {
				Expect(orderbook.GetBestAsk().GetAsks()).To(HaveLen(2))
				Expect(orderbook.GetBestAsk().GetBids()).To(HaveLen(0))
			})
		})
	})
	Describe("Deleting an Order", func() {
		var sellOrder = ob.NewOrder("CPTY2", ob.SELL, 101, 1)
		var buyOrder = ob.NewOrder("CPTY1", ob.BUY, 99, 1)
		BeforeEach(func() {
			orderbook.InsertOrder(sellOrder)
			orderbook.InsertOrder(buyOrder)
		})
		Context("If side is Buy", func() {
			BeforeEach(func() {
				orderbook.DeleteOrder(buyOrder)
			})
			It("should be removed from the Orderbook and no longer available from GetBestBid", func() {
				Expect(orderbook.GetBestBid()).Should(BeNil())
			})
		})
		Context("If side is Sell", func() {
			BeforeEach(func() {
				orderbook.DeleteOrder(sellOrder)
			})
			It("should be added to the Orderbook and available from GetBestAsk", func() {
				Expect(orderbook.GetBestAsk()).Should(BeNil())
			})
		})
	})
	Describe("Complex set of events", func() {
		var sellOrder103 = ob.NewOrder("CPTY3", ob.SELL, 103, 1)
		var sellOrder102 = ob.NewOrder("CPTY2", ob.SELL, 102, 1)
		var sellOrder101 = ob.NewOrder("CPTY1", ob.SELL, 101, 1)
		var buyOrder99 = ob.NewOrder("CPTY4", ob.BUY, 99, 1)
		var buyOrder98 = ob.NewOrder("CPTY5", ob.BUY, 98, 1)
		var buyOrder97 = ob.NewOrder("CPTY6", ob.BUY, 97, 1)
		BeforeEach(func() {
			orderbook.InsertOrder(sellOrder103)
			orderbook.InsertOrder(sellOrder102)
			orderbook.InsertOrder(sellOrder101)
			orderbook.InsertOrder(buyOrder99)
			orderbook.InsertOrder(buyOrder98)
			orderbook.InsertOrder(buyOrder97)
		})
		Context("Check orderbook status", func() {
			It("should have the correct outstanding orders", func() {
				Expect(orderbook.GetBestBid()).ShouldNot(BeNil())
				Expect(orderbook.GetBestBid().GetPrice()).Should(Equal(99.0))
				Expect(orderbook.GetBestAsk()).ShouldNot(BeNil())
				Expect(orderbook.GetBestAsk().GetPrice()).Should(Equal(101.0))
			})
		})
		Context("Delete top level bid and ask", func() {
			BeforeEach(func() {
				orderbook.DeleteOrder(sellOrder101)
				orderbook.DeleteOrder(buyOrder99)
			})
			It("should have new top levels", func() {
				Expect(orderbook.GetBestBid()).ShouldNot(BeNil())
				Expect(orderbook.GetBestBid().GetPrice()).Should(Equal(98.0))
				Expect(orderbook.GetBestAsk()).ShouldNot(BeNil())
				Expect(orderbook.GetBestAsk().GetPrice()).Should(Equal(102.0))
			})
		})
		Context("Delete back level bid and ask", func() {
			BeforeEach(func() {
				orderbook.DeleteOrder(sellOrder103)
				orderbook.DeleteOrder(buyOrder97)
			})
			It("should not have new top levels", func() {
				Expect(orderbook.GetBestBid()).ShouldNot(BeNil())
				Expect(orderbook.GetBestBid().GetPrice()).Should(Equal(99.0))
				Expect(orderbook.GetBestAsk()).ShouldNot(BeNil())
				Expect(orderbook.GetBestAsk().GetPrice()).Should(Equal(101.0))
			})
		})
	})
})
