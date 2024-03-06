package models

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"time"
)

type MatchResult string

const (
	MATCH_RESULT_IN_PROGRESS                   MatchResult = "in_progress"
	MATCH_RESULT_WHITE_WINS_BY_CHECKMATE       MatchResult = "white_wins_by_checkmate"
	MATCH_RESULT_BLACK_WINS_BY_CHECKMATE       MatchResult = "black_wins_by_checkmate"
	MATCH_RESULT_WHITE_WINS_BY_RESIGNATION     MatchResult = "white_wins_by_resignation"
	MATCH_RESULT_BLACK_WINS_BY_RESIGNATION     MatchResult = "black_wins_by_resignation"
	MATCH_RESULT_WHITE_WINS_BY_TIMEOUT         MatchResult = "white_wins_by_timeout"
	MATCH_RESULT_BLACK_WINS_BY_TIMEOUT         MatchResult = "black_wins_by_timeout"
	MATCH_RESULT_DRAW_BY_STALEMATE             MatchResult = "draw_by_stalemate"
	MATCH_RESULT_DRAW_BY_INSUFFICIENT_MATERIAL MatchResult = "draw_by_insufficient_material"
	MATCH_RESULT_DRAW_BY_THREEFOLD_REPETITION  MatchResult = "draw_by_threefold_repetition"
	MATCH_RESULT_DRAW_BY_FIFTY_MOVE_RULE       MatchResult = "draw_by_fifty_move_rule"
)

type Match struct {
	Uuid                  string       `json:"uuid"`
	Board                 *chess.Board `json:"board"`
	WhiteClientKey        Key          `json:"whiteClientKey"`
	WhiteTimeRemainingSec float64      `json:"whiteTimeRemainingSec"`
	BlackClientKey        Key          `json:"blackClientKey"`
	BlackTimeRemainingSec float64      `json:"blackTimeRemainingSec"`
	TimeControl           *TimeControl `json:"timeControl"`
	BotName               string       `json:"botName"`
	LastMove              *chess.Move  `json:"lastMove"`
	LastMoveTime          *time.Time   `json:"-"`
	Result                MatchResult  `json:"result"`
}

func (m *Match) Topic() MessageTopic {
	return MessageTopic(fmt.Sprintf("match-%s", m.Uuid))
}
