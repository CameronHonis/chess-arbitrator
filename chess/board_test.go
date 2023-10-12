package chess_test

import (
	"github.com/CameronHonis/chess-arbitrator/chess"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Board", func() {
	Describe("::IsForcedDraw", func() {
		When("the remaining material forces a draw", func() {
			When("only the kings are on the board", func() {
				It("returns true", func() {
					board, _ := chess.BoardFromFEN("8/8/8/8/7K/8/8/k7 w - - 0 1")
					Expect(board.IsForcedDraw()).To(BeTrue())
				})
			})
			When("only one bishop exists", func() {
				Context("and its a white bishop", func() {
					It("returns true", func() {
						board, _ := chess.BoardFromFEN("8/4B3/8/8/7K/8/8/k7 w - - 0 1")
						Expect(board.IsForcedDraw()).To(BeTrue())
					})
				})
				Context("and its a black bishop", func() {
					It("returns true", func() {
						board, _ := chess.BoardFromFEN("8/4b3/8/8/7K/8/8/k7 w - - 0 1")
						Expect(board.IsForcedDraw()).To(BeTrue())
					})
				})
			})
			When("both players have one bishop", func() {
				It("returns true", func() {
					board, _ := chess.BoardFromFEN("8/4B3/8/8/7K/8/3b4/k7 w - - 0 1")
					Expect(board.IsForcedDraw()).To(BeTrue())
				})
			})
			When("a player has two same colored bishops", func() {
				Context("and the bishops are white", func() {
					It("returns true", func() {
						board, _ := chess.BoardFromFEN("8/4B3/3B4/8/7K/8/8/k7 w - - 0 1")
						Expect(board.IsForcedDraw()).To(BeTrue())
					})
				})
				Context("and the bishops are black", func() {
					It("returns true", func() {
						board, _ := chess.BoardFromFEN("8/5b2/4b3/8/7K/8/8/k7 w - - 0 1")
						Expect(board.IsForcedDraw()).To(BeTrue())
					})
				})
			})
			When("only one knight exists", func() {
				Context("and the knight is white", func() {
					It("returns true", func() {
						board, _ := chess.BoardFromFEN("8/8/8/8/4N2K/8/8/k7 w - - 0 1")
						Expect(board.IsForcedDraw()).To(BeTrue())
					})
				})
				Context("and the knight is black", func() {
					It("returns true", func() {
						board, _ := chess.BoardFromFEN("8/8/4n3/8/7K/8/8/k7 w - - 0 1")
						Expect(board.IsForcedDraw()).To(BeTrue())
					})
				})
			})
			When("both players have only one knight", func() {
				It("returns true", func() {
					board, _ := chess.BoardFromFEN("5n2/8/8/8/4N2K/8/8/k7 w - - 0 1")
					Expect(board.IsForcedDraw()).To(BeTrue())
				})
			})
		})
		When("the remaining material does not force a draw", func() {
			When("only one rook exists", func() {
				It("returns false", func() {
					board, _ := chess.BoardFromFEN("8/2R5/8/8/7K/8/8/k7 w - - 0 1")
					Expect(board.IsForcedDraw()).To(BeFalse())
				})
			})
			When("only one queen exists", func() {
				It("returns false", func() {
					board, _ := chess.BoardFromFEN("8/8/4q3/8/7K/8/8/k7 w - - 0 1")
					Expect(board.IsForcedDraw()).To(BeFalse())
				})
			})
			When("only one pawn exists", func() {
				It("returns false", func() {
					board, _ := chess.BoardFromFEN("8/8/8/6P1/7K/8/8/k7 w - - 0 1")
					Expect(board.IsForcedDraw()).To(BeFalse())
				})
			})
			When("a player has two different colored bishops", func() {
				It("returns false", func() {
					board, _ := chess.BoardFromFEN("8/4B3/4B3/8/7K/8/8/k7 w - - 0 1")
					Expect(board.IsForcedDraw()).To(BeFalse())
				})
			})
			When("only one player has a bishop and knight", func() {
				It("returns false", func() {
					board, _ := chess.BoardFromFEN("8/4B3/4N3/8/7K/8/8/k7 w - - 0 1")
					Expect(board.IsForcedDraw()).To(BeFalse())
				})
			})
			When("only one player has two knights", func() {
				It("returns false", func() {
					board, _ := chess.BoardFromFEN("8/8/4n3/4n3/7K/8/8/k7 w - - 0 1")
					Expect(board.IsForcedDraw()).To(BeFalse())
				})
			})
		})
	})
	Describe("::getMaterialCount", func() {
		It("counts material of the initiate board", func() {
			board := chess.GetInitBoard()
			materialCount := board.ComputeMaterialCount(true)
			Expect(materialCount.WhitePawnCount).To(Equal(uint8(8)))
			Expect(materialCount.WhiteKnightCount).To(Equal(uint8(2)))
			Expect(materialCount.WhiteLightBishopCount).To(Equal(uint8(1)))
			Expect(materialCount.WhiteDarkBishopCount).To(Equal(uint8(1)))
			Expect(materialCount.WhiteRookCount).To(Equal(uint8(2)))
			Expect(materialCount.WhiteQueenCount).To(Equal(uint8(1)))
			Expect(materialCount.BlackPawnCount).To(Equal(uint8(8)))
			Expect(materialCount.BlackKnightCount).To(Equal(uint8(2)))
			Expect(materialCount.BlackLightBishopCount).To(Equal(uint8(1)))
			Expect(materialCount.BlackDarkBishopCount).To(Equal(uint8(1)))
			Expect(materialCount.BlackRookCount).To(Equal(uint8(2)))
			Expect(materialCount.BlackQueenCount).To(Equal(uint8(1)))
		})
	})
	Describe("#BoardFromFEN", func() {
		When("the FEN is valid", func() {
			When("the FEN is the initial board", func() {
				It("returns exactly the init board", func() {
					fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
					expBoard := chess.GetInitBoard()
					board, err := chess.BoardFromFEN(fen)
					Expect(err).ToNot(HaveOccurred())
					Expect(board).ToNot(BeNil())
					Expect(board.FullMoveCount).To(Equal(expBoard.FullMoveCount))
					Expect(board.HalfMoveClockCount).To(Equal(expBoard.HalfMoveClockCount))
					Expect(board.CanBlackCastleKingside).To(Equal(expBoard.CanBlackCastleKingside))
					Expect(board.CanBlackCastleQueenside).To(Equal(expBoard.CanBlackCastleQueenside))
					Expect(board.CanWhiteCastleKingside).To(Equal(expBoard.CanWhiteCastleKingside))
					Expect(board.CanWhiteCastleQueenside).To(Equal(expBoard.CanWhiteCastleQueenside))
					Expect(board.IsWhiteTurn).To(Equal(expBoard.IsWhiteTurn))
					Expect(board.OptEnPassantSquare).To(Equal(expBoard.OptEnPassantSquare))
					for i := 0; i < 8; i++ {
						for j := 0; j < 8; j++ {
							piece := board.Pieces[i][j]
							expPiece := board.Pieces[i][j]
							Expect(piece).To(Equal(expPiece))
						}
					}
				})
			})
			When("the FEN specifies that neither player has castle rights", func() {
				It("returns a board with all castle rights revoked", func() {
					fen := "3R2R1/8/2R5/2Rk2R1/4R3/2R5/R2R4/8 w - - 0 1"
					board, err := chess.BoardFromFEN(fen)
					Expect(err).ToNot(HaveOccurred())
					Expect(board.CanWhiteCastleQueenside).To(BeFalse())
					Expect(board.CanWhiteCastleKingside).To(BeFalse())
					Expect(board.CanBlackCastleQueenside).To(BeFalse())
					Expect(board.CanBlackCastleKingside).To(BeFalse())
				})
			})
			When("two white kings exist in the FEN pieces", func() {
				It("parses the board with no errors", func() {
					fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBKKBNR w KQkq - 0 1"
					_, err := chess.BoardFromFEN(fen)
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
		When("the FEN is not valid", func() {
			Context("the FEN does not have the correct amount of segments", func() {
				It("returns an error", func() {
					invalidFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0"
					_, err := chess.BoardFromFEN(invalidFEN)
					Expect(err).To(HaveOccurred())
				})
			})
			When("the issue is with the pieces", func() {
				Context("the FEN has too many pieces rows", func() {
					It("returns an error", func() {
						invalidFEN := "rnbqkbnr/pppppppp/8/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
						_, err := chess.BoardFromFEN(invalidFEN)
						Expect(err).To(HaveOccurred())
					})
				})
				Context("the FEN has one too few rows in pieces", func() {
					It("returns an error", func() {
						invalidFEN := "rnbqkbnr/pppppppp/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
						_, err := chess.BoardFromFEN(invalidFEN)
						Expect(err).To(HaveOccurred())
					})
				})
				Context("the FEN has too many pieces on the first row", func() {
					It("returns an error", func() {
						invalidFEN := "rnbqkbnrp/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
						_, err := chess.BoardFromFEN(invalidFEN)
						Expect(err).To(HaveOccurred())
					})
				})
				Context("the FEN has too few pieces on the first row", func() {
					It("returns an error", func() {
						invalidFEN := "rnbqkbn/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
						_, err := chess.BoardFromFEN(invalidFEN)
						Expect(err).To(HaveOccurred())
					})
				})
				Context("the FEN contains invalid piece chars", func() {
					It("returns an error", func() {
						invalidFEN := "xxxxxxxx/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
						_, err := chess.BoardFromFEN(invalidFEN)
						Expect(err).To(HaveOccurred())
					})
				})
			})
			Context("the FEN does not have a valid turn specifier character", func() {
				It("returns an error", func() {
					invalidFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR X KQkq - 0 1"
					_, err := chess.BoardFromFEN(invalidFEN)
					Expect(err).To(HaveOccurred())
				})
			})
			Context("the FEN contains an invalid castle rights specifier", func() {
				It("returns an error", func() {
					invalidFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w XQkq - 0 1"
					_, err := chess.BoardFromFEN(invalidFEN)
					Expect(err).To(HaveOccurred())
				})
			})
			Context("the FEN contains an invalid enPassant square", func() {
				It("returns an error", func() {
					invalidFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq jk4 0 1"
					_, err := chess.BoardFromFEN(invalidFEN)
					Expect(err).To(HaveOccurred())
				})
			})
			When("the issue is with the HalfMoveClockCount", func() {
				Context("the FEN contains a halfMoveClockCount greater than the range for a uint8", func() {
					It("returns an error", func() {
						invalidFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 277 1"
						_, err := chess.BoardFromFEN(invalidFEN)
						Expect(err).To(HaveOccurred())
					})
				})
				Context("the FEN contains a non-integer as the halfMoveClockCount", func() {
					It("returns an error", func() {
						invalidFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - X 1"
						_, err := chess.BoardFromFEN(invalidFEN)
						Expect(err).To(HaveOccurred())
					})
				})
			})
			When("the issue is with the FullMoveCount", func() {
				Context("the FEN contains a fullMoveCount greater than the range for a uint16", func() {
					It("returns an error", func() {
						invalidFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 70700"
						_, err := chess.BoardFromFEN(invalidFEN)
						Expect(err).To(HaveOccurred())
					})
				})
				Context("the FEN contains a non-integer as the fullMoveCount", func() {
					It("returns an error", func() {
						invalidFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 X"
						_, err := chess.BoardFromFEN(invalidFEN)
						Expect(err).To(HaveOccurred())
					})
				})
			})
		})
	})
})
