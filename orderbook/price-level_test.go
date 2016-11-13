package orderbook_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	ob "github.com/mdebrouwer/exchange/orderbook"
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

	Describe("Inserting a new Order to empty PriceLevel", func() {
		Context("If side is BUY", func() {
			var err error
			var trades []ob.Trade
			BeforeEach(func() {
				trades, err = priceLevel.InsertOrder(ob.NewOrder("CPTY1", ob.BUY, 100, 1))
			})
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should not cause trades", func() {
				Expect(trades).Should(BeEmpty())
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
			var trades []ob.Trade
			BeforeEach(func() {
				trades, err = priceLevel.InsertOrder(ob.NewOrder("CPTY1", ob.SELL, 100, 1))
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should not cause trades", func() {
				Expect(trades).To(BeEmpty())
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
			var trades []ob.Trade
			JustBeforeEach(func() {
				trades, err = priceLevel.InsertOrder(ob.NewOrder("CPTY1", ob.SELL, 101, 1))
			})

			It("should error", func() {
				Expect(err).To(HaveOccurred())
			})

			It("should not cause trades", func() {
				Expect(trades).To(BeEmpty())
			})

			It("not be available from GetBids", func() {
				Expect(priceLevel.GetBids()).To(BeEmpty())
			})

			It("not available from GetAsks", func() {
				Expect(priceLevel.GetAsks()).To(BeEmpty())
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
				order := ob.NewOrder("CPTY2", ob.BUY, 100, 1)
				priceLevel.InsertOrder(order)
				orderId = order.OrderId()
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
				order := ob.NewOrder("CPTY2", ob.SELL, 100, 1)
				priceLevel.InsertOrder(order)
				orderId = order.OrderId()
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
				orderId = 7
			})

			It("should error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
