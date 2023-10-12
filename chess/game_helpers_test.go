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

func compareMoves(expMoves *[]Move, realMoves *[]*Move) {
	Expect(*realMoves).To(HaveLen(len(*expMoves)))
	for _, realMove := range *realMoves {
		foundMatch := false
		for _, expMove := range *expMoves {
			if expMove.EqualTo(realMove) {
				foundMatch = true
				break
			}
		}
		Expect(foundMatch).To(BeTrue(), "unexpected move %+v", realMove)
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
	Describe("#GetLegalMovesForPawn", func() {
		When("the pawn can capture in either direction", func() {
			Context("and the pawn is not blocked", func() {
				It("returns 2 capture moves and one non-capture move", func() {
					board, err := BoardFromFEN("k1K5/8/8/4p1n1/5P2/8/8/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{4, 6})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{4, 6}, &Square{5, 6}, EMPTY, make([]*Square, 0), EMPTY},
						{WHITE_PAWN, &Square{4, 6}, &Square{5, 5}, BLACK_PAWN, make([]*Square, 0), EMPTY},
						{WHITE_PAWN, &Square{4, 6}, &Square{5, 7}, BLACK_KNIGHT, make([]*Square, 0), EMPTY},
					}
					compareMoves(&expMoves, realMoves)
				})
			})
			Context("and the pawn is blocked", func() {
				It("returns 2 capture moves", func() {
					board, err := BoardFromFEN("k1K5/8/8/4ppn1/5P2/8/8/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{4, 6})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{4, 6}, &Square{5, 5}, BLACK_PAWN, make([]*Square, 0), EMPTY},
						{WHITE_PAWN, &Square{4, 6}, &Square{5, 7}, BLACK_KNIGHT, make([]*Square, 0), EMPTY},
					}
					compareMoves(&expMoves, realMoves)
				})
			})
		})
		When("the pawn can capture in one direction", func() {
			Context("and the pawn is not blocked", func() {
				It("returns 1 capture moves and a non-capturing move", func() {
					board, err := BoardFromFEN("k1K5/8/8/4p3/5P2/8/8/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{4, 6})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{4, 6}, &Square{5, 6}, EMPTY, make([]*Square, 0), EMPTY},
						{WHITE_PAWN, &Square{4, 6}, &Square{5, 5}, BLACK_PAWN, make([]*Square, 0), EMPTY},
					}
					compareMoves(&expMoves, realMoves)
				})
			})
			Context("and the pawn is blocked", func() {
				It("returns 1 capture moves", func() {
					board, err := BoardFromFEN("k7/8/8/4rq2/5P2/8/8/7K w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{4, 6})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{4, 6}, &Square{5, 5}, BLACK_ROOK, make([]*Square, 0), EMPTY},
					}
					compareMoves(&expMoves, realMoves)
				})
			})
		})
		When("the pawn cannot capture either direction", func() {
			Context("and the pawn is not blocked", func() {
				It("returns a non-capturing move", func() {
					board, err := BoardFromFEN("k1K5/8/8/8/5P2/8/8/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{4, 6})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{4, 6}, &Square{5, 6}, EMPTY, make([]*Square, 0), EMPTY},
					}
					compareMoves(&expMoves, realMoves)
				})
			})
			Context("and the pawn is blocked", func() {
				It("returns no moves", func() {
					board, err := BoardFromFEN("k1K5/8/8/5n2/5P2/8/8/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{4, 6})
					Expect(err).ToNot(HaveOccurred())
					Expect(*realMoves).To(HaveLen(0))
				})
			})
			Context("and the 'attacked' squares are occupied by friendly pieces", func() {
				It("returns only non-capturing moves", func() {
					board, err := BoardFromFEN("8/8/k7/8/3B1R2/4P3/7K/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{3, 5})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{3, 5}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY},
					}
					compareMoves(&expMoves, realMoves)
				})
			})
		})
		When("the pawn can capture en passant to the left", func() {
			It("includes the en passant capture move", func() {
				board, err := BoardFromFEN("k1K5/8/8/4pP2/8/8/8/8 w - e6 0 1")
				Expect(err).ToNot(HaveOccurred())
				realMoves, err := GetLegalMovesForPawn(board, &Square{5, 6})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_PAWN, &Square{5, 6}, &Square{6, 5}, BLACK_PAWN, make([]*Square, 0), EMPTY},
					{WHITE_PAWN, &Square{5, 6}, &Square{6, 6}, EMPTY, make([]*Square, 0), EMPTY},
				}
				compareMoves(&expMoves, realMoves)
			})
		})
		When("the pawn can capture en passant to the right", func() {
			It("includes the en passant capture move", func() {
				board, err := BoardFromFEN("k1K5/8/8/5Pp1/8/8/8/8 w - g6 0 1")
				Expect(err).ToNot(HaveOccurred())
				realMoves, err := GetLegalMovesForPawn(board, &Square{5, 6})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_PAWN, &Square{5, 6}, &Square{6, 7}, BLACK_PAWN, make([]*Square, 0), EMPTY},
					{WHITE_PAWN, &Square{5, 6}, &Square{6, 6}, EMPTY, make([]*Square, 0), EMPTY},
				}
				compareMoves(&expMoves, realMoves)
			})
		})
		When("the pawn is on the starting row", func() {
			Context("and both squares directly in front are not occupied", func() {
				It("returns all possible moves including a double jump", func() {
					board, err := BoardFromFEN("k1K5/8/8/8/8/8/1P6/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{2, 2})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{2, 2}, &Square{3, 2}, EMPTY, make([]*Square, 0), EMPTY},
						{WHITE_PAWN, &Square{2, 2}, &Square{4, 2}, EMPTY, make([]*Square, 0), EMPTY},
					}
					compareMoves(&expMoves, realMoves)
				})
			})
			Context("and the double jump square is blocked", func() {
				It("returns only the single jump non-capturing move", func() {
					board, err := BoardFromFEN("k1K5/8/8/8/1B6/8/1P6/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{2, 2})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{2, 2}, &Square{3, 2}, EMPTY, make([]*Square, 0), EMPTY},
					}
					compareMoves(&expMoves, realMoves)
				})
			})
			Context("and the square directly in front is occupied", func() {
				It("returns no moves", func() {
					board, err := BoardFromFEN("k1K5/8/8/8/8/1R6/1P6/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{2, 2})
					Expect(err).ToNot(HaveOccurred())
					Expect(*realMoves).To(HaveLen(0))
				})
			})
		})
		When("the pawn can be promoted", func() {
			When("and the pawn can capture", func() {
				Context("and the square in front is not occupied", func() {
					It("returns both capturing and non-capturing promotion moves", func() {
						board, err := BoardFromFEN("r7/1P6/8/8/8/8/8/k1K5 w - - 0 1")
						Expect(err).ToNot(HaveOccurred())
						realMoves, err := GetLegalMovesForPawn(board, &Square{7, 2})
						Expect(err).ToNot(HaveOccurred())
						expMoves := []Move{
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, make([]*Square, 0), WHITE_KNIGHT},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, make([]*Square, 0), WHITE_BISHOP},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, make([]*Square, 0), WHITE_ROOK},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, make([]*Square, 0), WHITE_QUEEN},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 1}, BLACK_ROOK, make([]*Square, 0), WHITE_KNIGHT},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 1}, BLACK_ROOK, make([]*Square, 0), WHITE_BISHOP},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 1}, BLACK_ROOK, make([]*Square, 0), WHITE_ROOK},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 1}, BLACK_ROOK, make([]*Square, 0), WHITE_QUEEN},
						}
						compareMoves(&expMoves, realMoves)
					})
				})
				Context("and the square in front is occupied", func() {
					It("only returns capture promotion moves", func() {
						board, err := BoardFromFEN("rn6/1P6/8/8/8/8/8/k1K5 w - - 0 1")
						Expect(err).ToNot(HaveOccurred())
						realMoves, err := GetLegalMovesForPawn(board, &Square{7, 2})
						Expect(err).ToNot(HaveOccurred())
						expMoves := []Move{
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 1}, BLACK_ROOK, make([]*Square, 0), WHITE_KNIGHT},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 1}, BLACK_ROOK, make([]*Square, 0), WHITE_BISHOP},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 1}, BLACK_ROOK, make([]*Square, 0), WHITE_ROOK},
							{WHITE_PAWN, &Square{7, 2}, &Square{8, 1}, BLACK_ROOK, make([]*Square, 0), WHITE_QUEEN},
						}
						compareMoves(&expMoves, realMoves)
					})
				})
			})
			Context("and the pawn cannot capture", func() {
				It("returns non-capturing promotion moves", func() {
					board, err := BoardFromFEN("8/1P6/8/8/8/8/8/k1K5 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{7, 2})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, make([]*Square, 0), WHITE_KNIGHT},
						{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, make([]*Square, 0), WHITE_BISHOP},
						{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, make([]*Square, 0), WHITE_ROOK},
						{WHITE_PAWN, &Square{7, 2}, &Square{8, 2}, EMPTY, make([]*Square, 0), WHITE_QUEEN},
					}
					compareMoves(&expMoves, realMoves)
				})
			})
		})
		When("the pawn is pinned to its king", func() {
			Context("and the pin is coming from a piece on the same file", func() {
				It("returns only moves that block the pin", func() {
					board, err := BoardFromFEN("k7/4r3/8/8/8/4P3/4K3/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{3, 5})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{3, 5}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY},
					}
					compareMoves(&expMoves, realMoves)
				})
			})
			Context("and the pin is coming from a piece on the same diagonal", func() {
				It("returns no moves", func() {
					board, err := BoardFromFEN("k7/8/1q6/8/8/4P3/5K2/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{3, 5})
					Expect(err).ToNot(HaveOccurred())
					Expect(*realMoves).To(HaveLen(0))
				})
			})
			Context("and the pin is coming from the same rank", func() {
				It("returns no moves", func() {
					board, err := BoardFromFEN("k7/8/8/8/8/1r2PK2/8/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{3, 5})
					Expect(err).ToNot(HaveOccurred())
					Expect(*realMoves).To(HaveLen(0))
				})
			})
		})
		When("the pawn can block a check", func() {
			It("returns the move that would block the check", func() {
				board, err := BoardFromFEN("k7/1b1B4/8/8/8/4PK2/8/8 w - - 0 1")
				Expect(err).ToNot(HaveOccurred())
				realMoves, err := GetLegalMovesForPawn(board, &Square{3, 5})
				Expect(err).ToNot(HaveOccurred())
				expMoves := []Move{
					{WHITE_PAWN, &Square{3, 5}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY},
				}
				compareMoves(&expMoves, realMoves)
			})
			Context("and the pawn can capture", func() {
				It("only returns the non-capturing move that blocks the check", func() {
					board, err := BoardFromFEN("k7/1b1B4/8/8/3q4/4PK2/8/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{3, 5})
					Expect(err).ToNot(HaveOccurred())
					expMoves := []Move{
						{WHITE_PAWN, &Square{3, 5}, &Square{4, 5}, EMPTY, make([]*Square, 0), EMPTY},
					}
					compareMoves(&expMoves, realMoves)
				})
			})
			Context("and the friendly king is under double check", func() {
				It("returns no moves", func() {
					board, err := BoardFromFEN("k7/1b1B4/8/4n3/8/4PK2/8/8 w - - 0 1")
					Expect(err).ToNot(HaveOccurred())
					realMoves, err := GetLegalMovesForPawn(board, &Square{3, 5})
					Expect(err).ToNot(HaveOccurred())
					Expect(*realMoves).To(HaveLen(0))
				})
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
