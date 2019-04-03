package set_test

import (
	Set "algorithm/set"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Set", func() {

	var (
		mySet *Set.HashSet
	)
	BeforeEach(func() {
		mySet = Set.New()
	})

	Describe("Size", func() {
		Context("New set", func() {
			It("Should have 0 in size", func() {
				Expect(mySet.Size()).To(Equal(0))
			})
		})
	})

	Describe("Add item", func() {
		BeforeEach(func() {
			//mySet = Set.New()
			mySet.Add(5)
		})
		Context("New set", func() {

			It("Should have 1 in size", func() {
				Expect(mySet.Size()).To(Equal(1))
			})

			It("Should contain added value", func() {
				Expect(mySet.Contains(5)).To(Equal(true))
			})
		})

		Context("Add exist item", func() {

			It("Should have 1 in size", func() {
				mySet.Add(5)
				Expect(mySet.Size()).To(Equal(1))
			})

			It("Should contain added value", func() {
				mySet.Add(5)
				Expect(mySet.Contains(5)).To(Equal(true))
			})
		})
	})

	Describe("Contains", func() {
		BeforeEach(func() {
			mySet.Add(1)
			mySet.Add(2)
			mySet.Add(3)
		})
		It("Should should contain 1", func() {
			Expect(mySet.Contains(1)).To(Equal(true))
		})
		It("Should should not contain 4", func() {
			Expect(mySet.Contains(4)).To(Equal(false))
		})
	})

	Describe("Remove item", func() {
		var result bool
		BeforeEach(func() {
			//mySet = Set.New()
			mySet.Add(5)
			mySet.Add(4)
			mySet.Add(3)
		})
		Context("existing item", func() {
			BeforeEach(func() {
				result = mySet.Remove(5)
			})
			It("return type is true", func() {
				Expect(result).To(Equal(true))
			})
			It("Should not contain value", func() {
				Expect(mySet.Contains(5)).To(Equal(false))
			})
		})

		Context("Not existed item", func() {
			BeforeEach(func() {
				result = mySet.Remove(6)
			})

			It("return false", func() {
				Expect(result).To(Equal(false))
			})

			It("Should not contain value", func() {
				Expect(mySet.Contains(6)).To(Equal(false))
			})
		})
	})

	Describe("Get value", func() {
		Context("New Set", func() {
			It("Should return empty slide", func() {
				Expect(len(mySet.Values())).To(Equal(0))
			})
		})
		Context("Set have values", func() {
			BeforeEach(func() {
				//mySet = Set.New()
				mySet.Add(5)
				mySet.Add(4)
				mySet.Add(3)
			})

			It("Should return values slide of size 3", func() {
				Expect(len(mySet.Values())).To(Equal(3))
			})
		})

		Context("Set have some removed values", func() {
			BeforeEach(func() {
				//mySet = Set.New()
				mySet.Add(5)
				mySet.Add(4)
				mySet.Add(3)
				mySet.Remove(4)
			})

			It("Should return values slide of size 2", func() {
				Expect(len(mySet.Values())).To(Equal(2))
			})
		})


	})

})
