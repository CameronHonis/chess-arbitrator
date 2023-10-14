package geom_test

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/geom"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Geom", func() {
	Describe("Adding vectors", func() {
		Context("in 2d", func() {
			It("should add alike dimensions", func() {
				vecA := geom.Vector2d{1.23, 4.555}
				vecB := geom.Vector2d{1.337, 3.1415}
				vecC := geom.AddVector2ds(vecA, vecB)
				Expect(vecC.X).To(Equal(2.567))
				Expect(vecC.Y).To(Equal(7.6965))

				vecD := vecA.Add(vecB)
				Expect(vecD.X).To(Equal(2.567))
				Expect(vecD.Y).To(Equal(7.6965))
			})
		})
	})

	It("play around with array assignment", func() {
		a := [3]int{1, 2, 3}
		b := a
		b[0] = 100
		c := [2][2]int{{1, 2}, {3, 4}}
		d := c
		d[0][0] = 1337
		fmt.Sprintf("asdf %d", a[0])
	})
})
