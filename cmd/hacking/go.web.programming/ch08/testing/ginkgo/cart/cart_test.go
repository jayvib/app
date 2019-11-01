package cart

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestCart(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shopping Cart Suite")
}

var _ = Describe("Shopping cart", func() {
	Context("initially", func(){
		cart := Cart{}
		It("has 0 items", func(){
			Expect(cart.TotalUniqueItems()).Should(BeZero())
		})
		It("has 0 units", func(){
			Expect(cart.TotalUnits()).Should(BeZero())
		})
		Specify("the total amount is 0.00", func(){
			Expect(cart.TotalAmount()).Should(BeZero())
		})
	})

	Context("when a new item is added", func(){
		Context("the shopping cart", func(){
			It("has 1 more unique item than it had earlier", func(){})
			It("has 1 more units than it had earlier", func() {})
			Specify("total amount increases by item price", func(){})
		})
	})
})
