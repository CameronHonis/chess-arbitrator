package chess_test

import (
	. "github.com/CameronHonis/chess-arbitrator/chess"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Move", func() {
	Describe("::DoesAllowEnPassant", func() {
		When("the move does allow en passant", func() {
			move := Move{WHITE_PAWN, &Square{2, 5}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY}
			It("returns true", func() {
				Expect(move.DoesAllowEnPassant()).To(BeTrue())
			})
		})
		When("the move does not allow en passant", func() {
			move := Move{WHITE_PAWN, &Square{2, 5}, &Square{3, 5}, EMPTY, make([]*Square, 0), EMPTY}
			It("returns false", func() {
				Expect(move.DoesAllowEnPassant()).To(BeFalse())
			})
		})
	})
})
