package chess_test

import (
	"github.com/CameronHonis/chess-arbitrator/chess"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Square", func() {
	Describe("::IsValidBoardSquare", func() {
		When("The square lies outside of the board", func() {
			square := &chess.Square{Rank: 0, File: 0}
			Expect(square.IsValidBoardSquare()).To(BeFalse())
			square = &chess.Square{Rank: 0, File: 4}
			Expect(square.IsValidBoardSquare()).To(BeFalse())
			square = &chess.Square{Rank: 0, File: 9}
			Expect(square.IsValidBoardSquare()).To(BeFalse())
			square = &chess.Square{Rank: 4, File: 0}
			Expect(square.IsValidBoardSquare()).To(BeFalse())
			square = &chess.Square{Rank: 4, File: 9}
			Expect(square.IsValidBoardSquare()).To(BeFalse())
			square = &chess.Square{Rank: 9, File: 0}
			Expect(square.IsValidBoardSquare()).To(BeFalse())
			square = &chess.Square{Rank: 9, File: 5}
			Expect(square.IsValidBoardSquare()).To(BeFalse())
			square = &chess.Square{Rank: 9, File: 9}
			Expect(square.IsValidBoardSquare()).To(BeFalse())
		})
		When("The square is contained by the board", func() {
			square := &chess.Square{Rank: 4, File: 4}
			Expect(square.IsValidBoardSquare()).To(BeTrue())
			square = &chess.Square{Rank: 1, File: 1}
			Expect(square.IsValidBoardSquare()).To(BeTrue())
			square = &chess.Square{Rank: 8, File: 1}
			Expect(square.IsValidBoardSquare()).To(BeTrue())
			square = &chess.Square{Rank: 1, File: 8}
			Expect(square.IsValidBoardSquare()).To(BeTrue())
			square = &chess.Square{Rank: 8, File: 8}
			Expect(square.IsValidBoardSquare()).To(BeTrue())
		})
	})
	Describe("::ToAlgebraicCoords", func() {
		When("The square has rank=7, file=2", func() {
			It("returns 'b7'", func() {
				square := &chess.Square{Rank: 7, File: 2}
				Expect(square.ToAlgebraicCoords()).To(Equal("b7"))
			})
		})
		When("The square has rank=3, file=6", func() {
			It("returns 'f3'", func() {
				square := &chess.Square{Rank: 3, File: 6}
				Expect(square.ToAlgebraicCoords()).To(Equal("f3"))
			})
		})
		When("The square has rank=1, file=1", func() {
			It("returns 'a1'", func() {
				square := &chess.Square{Rank: 1, File: 1}
				Expect(square.ToAlgebraicCoords()).To(Equal("a1"))
			})
		})
		When("The square has rank=8, file=8", func() {
			It("returns 'h8'", func() {
				square := &chess.Square{Rank: 8, File: 8}
				Expect(square.ToAlgebraicCoords()).To(Equal("h8"))
			})
		})
	})
	Describe("::IsLightSquare", func() {
		When("The square is a light square", func() {
			It("returns true", func() {
				square := &chess.Square{Rank: 1, File: 2}
				Expect(square.IsLightSquare()).To(BeTrue())
				square = &chess.Square{Rank: 1, File: 4}
				Expect(square.IsLightSquare()).To(BeTrue())
				square = &chess.Square{Rank: 1, File: 6}
				Expect(square.IsLightSquare()).To(BeTrue())
				square = &chess.Square{Rank: 1, File: 8}
				Expect(square.IsLightSquare()).To(BeTrue())
				square = &chess.Square{Rank: 2, File: 1}
				Expect(square.IsLightSquare()).To(BeTrue())
				square = &chess.Square{Rank: 2, File: 3}
				Expect(square.IsLightSquare()).To(BeTrue())
				square = &chess.Square{Rank: 2, File: 5}
				Expect(square.IsLightSquare()).To(BeTrue())
				square = &chess.Square{Rank: 2, File: 7}
				Expect(square.IsLightSquare()).To(BeTrue())
			})
		})
		When("The square is a dark square", func() {
			It("returns false", func() {
				square := &chess.Square{Rank: 1, File: 1}
				Expect(square.IsLightSquare()).To(BeFalse())
				square = &chess.Square{Rank: 1, File: 3}
				Expect(square.IsLightSquare()).To(BeFalse())
				square = &chess.Square{Rank: 1, File: 5}
				Expect(square.IsLightSquare()).To(BeFalse())
				square = &chess.Square{Rank: 1, File: 7}
				Expect(square.IsLightSquare()).To(BeFalse())
				square = &chess.Square{Rank: 2, File: 2}
				Expect(square.IsLightSquare()).To(BeFalse())
				square = &chess.Square{Rank: 2, File: 4}
				Expect(square.IsLightSquare()).To(BeFalse())
				square = &chess.Square{Rank: 2, File: 6}
				Expect(square.IsLightSquare()).To(BeFalse())
				square = &chess.Square{Rank: 2, File: 8}
				Expect(square.IsLightSquare()).To(BeFalse())
			})
		})
	})
	Describe("#SquareFromAlgebraicCoords", func() {
		When("the algebraic coords are valid", func() {
			It("returns the corresponding Square", func() {
				square, err := chess.SquareFromAlgebraicCoords("a1")
				Expect(err).To(BeNil())
				Expect(square.Rank).To(Equal(uint8(1)))
				Expect(square.File).To(Equal(uint8(1)))
				square, err = chess.SquareFromAlgebraicCoords("c7")
				Expect(err).To(BeNil())
				Expect(square.Rank).To(Equal(uint8(7)))
				Expect(square.File).To(Equal(uint8(3)))
			})
			Context("And the coords contain capital letters", func() {
				It("returns the corresponding Square", func() {
					square, err := chess.SquareFromAlgebraicCoords("A1")
					Expect(err).To(BeNil())
					Expect(square.Rank).To(Equal(uint8(1)))
					Expect(square.File).To(Equal(uint8(1)))
				})
			})
		})
		When("the algebraic coords are outside of the board", func() {
			It("returns an error", func() {
				square, err := chess.SquareFromAlgebraicCoords("i5")
				Expect(square).To(BeNil())
				Expect(err).To(HaveOccurred())
				square, err = chess.SquareFromAlgebraicCoords("b9")
				Expect(square).To(BeNil())
				Expect(err).To(HaveOccurred())
				square, err = chess.SquareFromAlgebraicCoords("b0")
				Expect(square).To(BeNil())
				Expect(err).To(HaveOccurred())
			})
		})
		When("the algebraic coords isn't 2 chars long", func() {
			Context("the coords are too short", func() {
				It("returns an error", func() {
					square, err := chess.SquareFromAlgebraicCoords("a")
					Expect(square).To(BeNil())
					Expect(err).To(HaveOccurred())
					square, err = chess.SquareFromAlgebraicCoords("f")
					Expect(square).To(BeNil())
					Expect(err).To(HaveOccurred())
				})
			})
			Context("the coords are too long", func() {
				It("returns an error", func() {
					square, err := chess.SquareFromAlgebraicCoords("a12")
					Expect(square).To(BeNil())
					Expect(err).To(HaveOccurred())
					square, err = chess.SquareFromAlgebraicCoords("hh2")
					Expect(square).To(BeNil())
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
