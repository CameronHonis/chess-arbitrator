package chess_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/CameronHonis/chess-arbitrator/chess"
)

func compareSquares(expSquares *[]Square, realSquares *[]*Square) {
	Expect(*realSquares).To(HaveLen(len(*expSquares)))
	for _, realSquare := range *realSquares {
		foundMatch := false
		for _, expSquare := range *expSquares {
			if realSquare.Rank == expSquare.Rank && realSquare.File == expSquare.File {
				foundMatch = true
				break
			}
		}
		Expect(foundMatch).To(Equal(true), "unexpected square %+v", realSquare)
	}
}

var _ = Describe("GameHelpers", func() {
	Describe("#GetCheckingKingSquares", func() {
		When("the board is the initial board", func() {
			It("returns an empty list", func() {
				board := GetInitBoard()
				checkingSquares := GetCheckingSquares(board, board.IsWhiteTurn)
				expSquares := make([]Square, 0)
				compareSquares(&expSquares, checkingSquares)
			})
		})
		When("there are multiple rooks on the board", func() {
			It("returns the square of each rook that is checking the king", func() {
				board, err := BoardFromFEN("3R2R1/8/2R5/2Rk2R1/4R3/2R5/R2R4/8 w - - 0 1")
				Expect(err).ToNot(HaveOccurred())
				rookSquares := GetCheckingSquares(board, false)
				expSquares := []Square{
					{2, 4},
					{8, 4},
					{5, 3},
					{5, 7},
				}
				compareSquares(&expSquares, rookSquares)

			})
		})
		When("there are multiple bishops on the board", func() {
			It("returns the square of each bishop that is checking the king", func() {
				board, err := BoardFromFEN("3BB2B/5B2/B2B1k1B/8/4BB1B/8/8/B7 w - - 0 1")
				Expect(err).ToNot(HaveOccurred())
				bishopSquares := GetCheckingSquares(board, false)
				expSquares := []Square{
					{1, 1},
					{8, 8},
					{8, 4},
					{4, 8},
				}
				compareSquares(&expSquares, bishopSquares)
			})
		})
		When("there are multiple pawns on the board", func() {
			Context("when the pawns are white", func() {
				It("returns the square of each pawn that is checking the king", func() {
					board, err := BoardFromFEN("3PP2P/4PPP1/P2PPkPP/4PPP1/4PP1P/8/8/P7 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					pawnSquares := GetCheckingSquares(board, false)
					expSquares := []Square{
						{5, 5},
						{5, 7},
					}
					compareSquares(&expSquares, pawnSquares)
				})
			})
			Context("when the pawns are black", func() {
				It("returns the square of each pawn that is checking the king", func() {
					board, err := BoardFromFEN("3pp2p/4ppp1/p2ppKpp/4ppp1/4pp1p/8/8/p7 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					pawnSquares := GetCheckingSquares(board, true)
					expSquares := []Square{
						{7, 5},
						{7, 7},
					}
					compareSquares(&expSquares, pawnSquares)
				})
			})
		})
		When("there are multiple knights on the board", func() {
			It("returns the square of each knight checking the king", func() {
				board, err := BoardFromFEN("4NNN1/3NNNNN/2N2k1N/3N1N1N/4N1N1/8/8/5N2 w - - 0 1")
				Expect(err).ToNot(HaveOccurred())
				knightSquares := GetCheckingSquares(board, false)
				expSquares := []Square{
					{4, 5},
					{4, 7},
					{5, 4},
					{5, 8},
					{7, 4},
					{7, 8},
					{8, 5},
					{8, 7},
				}
				compareSquares(&expSquares, knightSquares)
			})
		})
		When("there are multiple queens on the board", func() {
			It("returns the square of each knight checking the king", func() {
				board, err := BoardFromFEN("3QQ1Q1/3Q1QQQ/Q4k1Q/3Q1QQQ/4Q1Q1/8/8/Q7 w - - 0 1")
				Expect(err).ToNot(HaveOccurred())
				queenSquares := GetCheckingSquares(board, false)
				expSquares := []Square{
					{1, 1},
					{6, 1},
					{8, 4},
					{7, 6},
					{7, 7},
					{6, 8},
					{5, 7},
					{5, 6},
				}
				compareSquares(&expSquares, queenSquares)
			})
		})
		When("only same color pieces as king exist", func() {
			It("returns no checking squares", func() {
				board, err := BoardFromFEN("3q4/4p3/2q2k2/4p3/3bn1n1/5r2/8/8 w - - 0 1")
				Expect(err).ToNot(HaveOccurred())
				checkingSquares := GetCheckingSquares(board, false)
				Expect(*checkingSquares).To(HaveLen(0))
			})
		})
		When("all 'blockable' pieces are blocked", func() {
			It("returns no checking squares", func() {
				board, err := BoardFromFEN("3Q4/4n3/1Qq2kpR/6N1/7B/2N2N2/1B3R2/8 w - - 0 1")
				Expect(err).ToNot(HaveOccurred())
				checkingSquares := GetCheckingSquares(board, false)
				Expect(*checkingSquares).To(HaveLen(0))
			})
		})
	})
	Describe("#GetResultingBoard", func() {
		When("the move is a capture", func() {

		})
		When("the move is not a capture", func() {

		})
		When("its black's move", func() {

		})
		When("a rook moves, revoking a castling right", func() {

		})
		When("a king moves", func() {

		})
		When("the moves is castles", func() {

		})
		When("the resulting board is a stalemate", func() {

		})
		When("the move violates the 50-move rule", func() {

		})
		When("the remaining material forces a draw", func() {

		})
		When("a winner emerges from the resulting board", func() {
			When("the resulting move results in stalemate", func() {

			})
		})
	})
})
