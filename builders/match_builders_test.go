package builders_test

import (
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/builders"
	"github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("MatchBuilder", func() {
	var matchBuilder *builders.MatchBuilder
	BeforeEach(func() {
		now := time.Now()
		matchBuilder = builders.NewMatchBuilder().FromMatch(&models.Match{
			Uuid:                  "some-uuid",
			Board:                 chess.GetInitBoard(),
			WhiteClientKey:        "some-white-client-key",
			WhiteTimeRemainingSec: 600,
			BlackClientKey:        "some-black-client-key",
			BlackTimeRemainingSec: 600,
			TimeControl:           &models.TimeControl{},
			LastMoveTime:          &now,
			Result:                models.MATCH_RESULT_IN_PROGRESS,
		})
	})
	Describe("WithBoard", func() {
		var board *chess.Board
		When("the match is non-terminal", func() {
			BeforeEach(func() {
				board = &chess.Board{}
			})
			It("should return a match with the given board", func() {
				newMatch := matchBuilder.WithBoard(board).Build()
				Expect(newMatch.Board).To(Equal(board))
			})
			It("does not modify the match result", func() {
				newMatch := matchBuilder.WithBoard(board).Build()
				Expect(newMatch.Result).To(Equal(models.MATCH_RESULT_IN_PROGRESS))
			})
		})
		When("the match is terminal", func() {
			When("white wins by checkmate", func() {
				BeforeEach(func() {
					board, _ = chess.BoardFromFEN("k7/1QK5/8/8/8/8/8/8 b - - 0 1")
					Expect(board.Result).To(Equal(chess.BOARD_RESULT_WHITE_WINS_BY_CHECKMATE))
				})
				It("returns a match with the given result", func() {
					newMatch := matchBuilder.WithBoard(board).Build()
					Expect(newMatch.Result).To(Equal(models.MATCH_RESULT_WHITE_WINS_BY_CHECKMATE))
				})
			})
			When("black wins by checkmate", func() {
				BeforeEach(func() {
					board, _ = chess.BoardFromFEN("K7/1qk5/8/8/8/8/8/8 w - - 0 1")
					Expect(board.Result).To(Equal(chess.BOARD_RESULT_BLACK_WINS_BY_CHECKMATE))
				})
				It("returns a match with the given result", func() {
					newMatch := matchBuilder.WithBoard(board).Build()
					Expect(newMatch.Result).To(Equal(models.MATCH_RESULT_BLACK_WINS_BY_CHECKMATE))
				})
			})
			When("it is a draw by stalemate", func() {
				BeforeEach(func() {
					board, _ = chess.BoardFromFEN("k7/8/1Q6/8/4K3/8/8/8 b - - 0 1")
					Expect(board.Result).To(Equal(chess.BOARD_RESULT_DRAW_BY_STALEMATE))
				})
				It("returns a match with the given result", func() {
					newMatch := matchBuilder.WithBoard(board).Build()
					Expect(newMatch.Result).To(Equal(models.MATCH_RESULT_DRAW_BY_STALEMATE))
				})
			})
			When("it is a draw by insufficient material", func() {
				BeforeEach(func() {
					board, _ = chess.BoardFromFEN("k7/8/2K5/8/4N3/8/8/8 w - - 0 1")
					Expect(board.Result).To(Equal(chess.BOARD_RESULT_DRAW_BY_INSUFFICIENT_MATERIAL))
				})
				It("returns a match with the given result", func() {
					newMatch := matchBuilder.WithBoard(board).Build()
					Expect(newMatch.Result).To(Equal(models.MATCH_RESULT_DRAW_BY_INSUFFICIENT_MATERIAL))
				})
			})
			When("it is a draw by fifty move rule", func() {
				BeforeEach(func() {
					board, _ = chess.BoardFromFEN("k7/8/8/8/2QK5/8/8/8 w - - 50 1")
					Expect(board.Result).To(Equal(chess.BOARD_RESULT_DRAW_BY_FIFTY_MOVE_RULE))
				})
				It("returns a match with the given result", func() {
					newMatch := matchBuilder.WithBoard(board).Build()
					Expect(newMatch.Result).To(Equal(models.MATCH_RESULT_DRAW_BY_FIFTY_MOVE_RULE))
				})
			})
			When("it is a draw by threefold repetition", func() {
				BeforeEach(func() {
					board, _ = chess.BoardFromFEN("k7/8/2K5/8/8/1P6/8/8 w - - 0 1")
					board.RepetitionsByMiniFEN[board.ToMiniFEN()] = 3
					board.Result = chess.BOARD_RESULT_DRAW_BY_THREEFOLD_REPETITION
				})
				It("returns a match with the given result", func() {
					newMatch := matchBuilder.WithBoard(board).Build()
					Expect(newMatch.Result).To(Equal(models.MATCH_RESULT_DRAW_BY_THREEFOLD_REPETITION))
				})
			})
		})
	})
	Describe("WithWhiteTimeRemainingSec", func() {
		It("sets the whiteTimeRemainingSec", func() {
			newMatch := matchBuilder.WithWhiteTimeRemainingSec(300).Build()
			Expect(newMatch.WhiteTimeRemainingSec).To(Equal(float64(300)))
		})
		When("the time is 0", func() {
			It("sets the result to black wins by timeout", func() {
				newMatch := matchBuilder.WithWhiteTimeRemainingSec(0).Build()
				Expect(newMatch.Result).To(Equal(models.MATCH_RESULT_BLACK_WINS_BY_TIMEOUT))
			})
		})
	})
	Describe("WithBlackTimeRemainingSec", func() {
		It("sets the blackTimeRemainingSec", func() {
			newMatch := matchBuilder.WithBlackTimeRemainingSec(300).Build()
			Expect(newMatch.BlackTimeRemainingSec).To(Equal(float64(300)))
		})
		When("the time is 0", func() {
			It("sets the result to black wins by timeout", func() {
				newMatch := matchBuilder.WithBlackTimeRemainingSec(0).Build()
				Expect(newMatch.Result).To(Equal(models.MATCH_RESULT_WHITE_WINS_BY_TIMEOUT))
			})
		})
	})
	Describe("WithResult", func() {
		When("the result is not terminal", func() {
			BeforeEach(func() {
				matchBuilder.WithResult(models.MATCH_RESULT_WHITE_WINS_BY_CHECKMATE)
			})
			It("sets the result", func() {
				newMatch := matchBuilder.WithResult(models.MATCH_RESULT_IN_PROGRESS).Build()
				Expect(newMatch.Result).To(Equal(models.MATCH_RESULT_IN_PROGRESS))
			})
			It("propagates the result to the board result", func() {
				newMatch := matchBuilder.WithResult(models.MATCH_RESULT_IN_PROGRESS).Build()
				Expect(newMatch.Board.Result).To(Equal(chess.BOARD_RESULT_IN_PROGRESS))
			})
		})
		When("the result is a checkmate", func() {
			It("sets the result", func() {
				newMatch := matchBuilder.WithResult(models.MATCH_RESULT_WHITE_WINS_BY_CHECKMATE).Build()
				Expect(newMatch.Result).To(Equal(models.MATCH_RESULT_WHITE_WINS_BY_CHECKMATE))
			})
			It("propagates the result to the board result", func() {
				newMatch := matchBuilder.WithResult(models.MATCH_RESULT_WHITE_WINS_BY_CHECKMATE).Build()
				Expect(newMatch.Board.Result).To(Equal(chess.BOARD_RESULT_WHITE_WINS_BY_CHECKMATE))
			})
		})
		When("the result is a stalemate", func() {
			It("sets the result", func() {
				newMatch := matchBuilder.WithResult(models.MATCH_RESULT_DRAW_BY_STALEMATE).Build()
				Expect(newMatch.Result).To(Equal(models.MATCH_RESULT_DRAW_BY_STALEMATE))
			})
			It("propagates the result to the board result", func() {
				newMatch := matchBuilder.WithResult(models.MATCH_RESULT_DRAW_BY_STALEMATE).Build()
				Expect(newMatch.Board.Result).To(Equal(chess.BOARD_RESULT_DRAW_BY_STALEMATE))
			})
		})
		When("the result is a terminal result 'off the board'", func() {
			It("sets the result", func() {
				newMatch := matchBuilder.WithResult(models.MATCH_RESULT_WHITE_WINS_BY_TIMEOUT).Build()
				Expect(newMatch.Result).To(Equal(models.MATCH_RESULT_WHITE_WINS_BY_TIMEOUT))
			})
			It("does not propagate the result to the board result", func() {
				newMatch := matchBuilder.WithResult(models.MATCH_RESULT_WHITE_WINS_BY_TIMEOUT).Build()
				Expect(newMatch.Board.Result).To(Equal(chess.BOARD_RESULT_IN_PROGRESS))
			})
		})
	})
})
