package set_test

import (
	. "github.com/CameronHonis/chess-arbitrator/set"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Set", func() {
	Describe("::Add", func() {
		It("adds the element to the set", func() {
			s := EmptySet[int]()
			s.Add(12)
			Expect(s.Size()).To(Equal(1))
		})
	})
	Describe("::Has", func() {
		When("the set does not have the item", func() {
			It("returns false", func() {
				s := EmptySet[string]()
				Expect(s.Has("asdf")).To(BeFalse())
			})
		})
		When("the set contains the item", func() {
			It("returns true", func() {
				s := EmptySet[string]()
				s.Add("asdf")
				Expect(s.Has("asdf")).To(BeTrue())
			})
		})
	})
	Describe("::Remove", func() {
		var s *Set[string]
		BeforeEach(func() {
			s = EmptySet[string]()
			s.Add("asdf")
			Expect(s.Size()).To(Equal(1))
			Expect(s.Has("asdf")).To(BeTrue())
		})
		It("removes the item from the list", func() {
			s.Remove("asdf")
			Expect(s.Has("asdf")).To(BeFalse())
			Expect(s.Size()).To(Equal(0))
		})
		When("the item does not exist in the set", func() {
			It("does not panic", func() {
				s.Remove("a")
			})
		})
	})
	Describe("::Flatten", func() {
		It("returns a slice of all the items in the set", func() {
			s := EmptySet[int]()
			s.Add(1)
			s.Add(2)
			s.Add(3)
			expFlattenedSet := []int{1, 2, 3}
			realFlattenedSet := s.Flatten()
			Expect(len(realFlattenedSet)).To(Equal(len(expFlattenedSet)))
			for _, item := range realFlattenedSet {
				foundMatch := false
				for _, expItem := range expFlattenedSet {
					if item == expItem {
						foundMatch = true
						break
					}
				}
				Expect(foundMatch).To(BeTrue())
			}
		})
	})
	Describe("#EmptySet", func() {
		It("returns a set of the specified generic", func() {
			s := EmptySet[uint16]()
			Expect(reflect.TypeOf(*s)).To(Equal(reflect.TypeOf(Set[uint16]{})))
		})
		It("returns an empty set", func() {
			s := EmptySet[string]()
			Expect(s.Size()).To(Equal(0))
		})
	})
})
