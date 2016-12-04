package orderbook_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"

	ob "github.com/mdebrouwer/exchange/orderbook"
	"github.com/mdebrouwer/exchange/uuid"
)

var _ = Describe("PriceLevel", func() {
	var priceLevel ob.PriceLevel
	BeforeEach(func() {
		priceLevel = ob.NewPriceLevel(100)
	})
	Describe("Creating a new PriceLevel", func() {
		It("should have the correct price", func() {
			Expect(priceLevel.GetPrice()).Should(BeNumerically("==", 100))
		})
		It("should have empty Bids", func() {
			Expect(priceLevel.GetBids()).To(BeEmpty())
		})
		It("should have empty Asks", func() {
			Expect(priceLevel.GetAsks()).To(BeEmpty())
		})
	})
	Describe("Inserting a new Order", func() {
		Context("If side is BUY", func() {
			var err error
			BeforeEach(func() {
				err = priceLevel.InsertOrder(ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY1", ob.BUY, 100, 1))
			})
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should be added to the PriceLevel and available from GetBids", func() {
				Expect(priceLevel.GetBids()).To(HaveLen(1))
			})
			It("should be added to the PriceLevel and not available from GetAsks", func() {
				Expect(priceLevel.GetAsks()).To(BeEmpty())
			})
		})
		Context("If side is SELL", func() {
			var err error
			BeforeEach(func() {
				err = priceLevel.InsertOrder(ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY1", ob.SELL, 100, 1))
			})
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should be added to the PriceLevel and available from GetAsks", func() {
				Expect(priceLevel.GetAsks()).To(HaveLen(1))
			})
			It("should be added to the PriceLevel and not available from GetBids", func() {
				Expect(priceLevel.GetBids()).To(BeEmpty())
			})
		})
		Context("If the Order price does not match the level", func() {
			var err error
			JustBeforeEach(func() {
				err = priceLevel.InsertOrder(ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY1", ob.SELL, 101, 1))
			})
			It("should error", func() {
				Expect(err).To(HaveOccurred())
			})
			It("not be available from GetBids", func() {
				Expect(priceLevel.GetBids()).To(BeEmpty())
			})
			It("not available from GetAsks", func() {
				Expect(priceLevel.GetAsks()).To(BeEmpty())
			})
		})
		Context("If a BUY Order is in cross", func() {
			var err error
			BeforeEach(func() {
				priceLevel.InsertOrder(ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY1", ob.SELL, 100, 1))
				err = priceLevel.InsertOrder(ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY2", ob.BUY, 100, 1))
			})
			It("should error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
		Context("If a SELL Order is in cross", func() {
			var err error
			BeforeEach(func() {
				priceLevel.InsertOrder(ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY1", ob.BUY, 100, 1))
				err = priceLevel.InsertOrder(ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY2", ob.SELL, 100, 1))
			})
			It("should error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Describe("Matching Orders", func() {
		var trades []ob.Trade
		var err error
		Context("If pricelevel empty and new order is BUY", func() {
			BeforeEach(func() {
				order := ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY2", ob.BUY, 100, 1)
				trades, err = priceLevel.MatchOrder(order)
			})
			It("should error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
		Context("If pricelevel empty and new order is SELL", func() {
			BeforeEach(func() {
				order := ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY2", ob.SELL, 100, 1)
				trades, err = priceLevel.MatchOrder(order)
			})
			It("should error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
		Context("If resting BUY quotes and new order is SELL", func() {
			BeforeEach(func() {
				priceLevel.InsertOrder(ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY1", ob.BUY, 100, 1))
				priceLevel.InsertOrder(ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY1", ob.BUY, 100, 1))
				priceLevel.InsertOrder(ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY1", ob.BUY, 100, 1))
				trades, err = priceLevel.MatchOrder(ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY2", ob.SELL, 100, 1))
			})
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should result in a trade", func() {
				trade := trades[0]
				Expect(trade).To(Equal(ob.NewTrade(trade.GetCreationTime(), ob.SELL, "CPTY1", "CPTY2", 100, 1)))
			})
			It("the BUY quote should be removed from the pricelevel", func() {
				Expect(priceLevel.GetBids()).To(HaveLen(2))
			})
		})
		Context("If resting SELL quotes and new order is BUY", func() {
			BeforeEach(func() {
				priceLevel.InsertOrder(ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY1", ob.SELL, 100, 1))
				priceLevel.InsertOrder(ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY1", ob.SELL, 100, 1))
				priceLevel.InsertOrder(ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY1", ob.SELL, 100, 1))
				trades, err = priceLevel.MatchOrder(ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY2", ob.BUY, 100, 1))
			})
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should result in a trade", func() {
				trade := trades[0]
				Expect(trade).To(Equal(ob.NewTrade(trade.GetCreationTime(), ob.BUY, "CPTY2", "CPTY1", 100, 1)))
			})
			It("the SELL quote should be removed from the pricelevel", func() {
				Expect(priceLevel.GetAsks()).To(HaveLen(2))
			})
		})
		Context("If resting BUY quote is partially filled by new SELL order", func() {
			var quote ob.Order
			BeforeEach(func() {
				quote = ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY1", ob.BUY, 100, 10)
				priceLevel.InsertOrder(quote)
				trades, err = priceLevel.MatchOrder(ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY2", ob.SELL, 100, 5))
			})
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should result in a trade", func() {
				trade := trades[0]
				Expect(trade).To(Equal(ob.NewTrade(trade.GetCreationTime(), ob.SELL, "CPTY1", "CPTY2", 100, 5)))
			})
			It("the BUY quote should still exist at the pricelevel", func() {
				Expect(priceLevel.GetBids()[0].GetOrderId()).To(Equal(quote.GetOrderId()))
				Expect(priceLevel.GetBids()[0].GetVersion()).To(Equal(ob.OrderVersion(2)))
				Expect(priceLevel.GetBids()[0].GetVolume()).To(Equal(ob.Volume(5.0)))
			})
		})
		Context("If resting SELL quote is partially filled by new BUY order", func() {
			var quote ob.Order
			BeforeEach(func() {
				quote = ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY1", ob.SELL, 100, 10)
				priceLevel.InsertOrder(quote)
				trades, err = priceLevel.MatchOrder(ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY2", ob.BUY, 100, 5))
			})
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should result in a trade", func() {
				trade := trades[0]
				Expect(trade).To(Equal(ob.NewTrade(trade.GetCreationTime(), ob.BUY, "CPTY2", "CPTY1", 100, 5)))
			})
			It("the SELL quote should still exist at the pricelevel", func() {
				Expect(priceLevel.GetAsks()[0].GetOrderId()).To(Equal(quote.GetOrderId()))
				Expect(priceLevel.GetAsks()[0].GetVersion()).To(Equal(ob.OrderVersion(2)))
				Expect(priceLevel.GetAsks()[0].GetVolume()).To(Equal(ob.Volume(5.0)))
			})
		})
	})
	Describe("Deleting an Order", func() {
		var err error
		var orderId ob.OrderId
		JustBeforeEach(func() {
			err = priceLevel.DeleteOrder(orderId)
		})
		Context("If side is BUY", func() {
			BeforeEach(func() {
				order := ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY2", ob.BUY, 100, 1)
				priceLevel.InsertOrder(order)
				orderId = order.GetOrderId()
			})
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should be removed from the PriceLevel and not available from GetBids", func() {
				Expect(priceLevel.GetBids()).To(BeEmpty())
			})
		})
		Context("If side is SELL", func() {
			BeforeEach(func() {
				order := ob.NewOrder(uuid.NewUUID(), time.Now(), "CPTY2", ob.SELL, 100, 1)
				priceLevel.InsertOrder(order)
				orderId = order.GetOrderId()
			})
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should be removed from the PriceLevel and not available from GetAsks", func() {
				Expect(priceLevel.GetAsks()).To(BeEmpty())
			})
		})
		Context("If the Order does not exist", func() {
			BeforeEach(func() {
				orderId = ob.OrderId(uuid.NewUUID())
			})
			It("should error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
