package orderbook_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	ob "github.com/mdebrouwer/exchange/orderbook"
)

var _ = Describe("PriceLevel", func() {
	var (
		priceLevel *ob.PriceLevel
	)

	BeforeEach(func() {
		priceLevel = ob.NewPriceLevel(100)
	})

	Describe("Creating a new PriceLevel", func() {
		It("should have the correct price", func() {
			Expect(priceLevel.GetPrice()).Should(BeNumerically("==", 100))
		})
		It("should have empty Bids", func() {
			var bids = priceLevel.GetBids()
			Expect(bids).ShouldNot(BeNil())
			Expect(len(bids)).To(Equal(0))
		})
		It("should have empty Asks", func() {
			Expect(priceLevel.GetAsks()).ShouldNot(BeNil())
			Expect(len(priceLevel.GetAsks())).To(Equal(0))
		})
	})

	Describe("Inserting a new Order to empty PriceLevel", func() {
		Context("If side is BUY", func() {
			var err error
			var trades []*ob.Trade
			BeforeEach(func() {
				trades, err = priceLevel.InsertOrder(ob.NewOrder("CPTY1", ob.BUY, 100, 1))
			})
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should be added to the PriceLevel and available from GetBids", func() {
				var bids = priceLevel.GetBids()
				Expect(bids).ShouldNot(BeNil())
				Expect(len(bids)).To(Equal(1))
			})
			It("should be added to the PriceLevel and not available from GetAsks", func() {
				var asks = priceLevel.GetAsks()
				Expect(asks).ShouldNot(BeNil())
				Expect(len(asks)).To(Equal(0))
			})
		})

		Context("If side is SELL", func() {
			var err error
			var trades []*ob.Trade
			BeforeEach(func() {
				trades, err = priceLevel.InsertOrder(ob.NewOrder("CPTY1", ob.SELL, 100, 1))
			})
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should be added to the PriceLevel and available from GetAsks", func() {
				var asks = priceLevel.GetAsks()
				Expect(asks).ShouldNot(BeNil())
				Expect(len(asks)).To(Equal(1))
			})
			It("should be added to the PriceLevel and not available from GetBids", func() {
				var bids = priceLevel.GetBids()
				Expect(bids).ShouldNot(BeNil())
				Expect(len(bids)).To(Equal(0))
			})
		})

		Context("If the Order price does not match the level", func() {
			var err error
			var trades []*ob.Trade
			BeforeEach(func() {
				trades, err = priceLevel.InsertOrder(ob.NewOrder("CPTY1", ob.SELL, 101, 1))
			})
			It("should error", func() {
				Expect(err).To(HaveOccurred())
			})
			It("not be available from GetBids", func() {
				var bids = priceLevel.GetBids()
				Expect(bids).ShouldNot(BeNil())
				Expect(len(bids)).To(Equal(0))
			})
			It("not available from GetAsks", func() {
				var asks = priceLevel.GetAsks()
				Expect(asks).ShouldNot(BeNil())
				Expect(len(asks)).To(Equal(0))
			})
		})
	})
})
