package builders

import (
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/helpers"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/google/uuid"
	"time"
)

func NewMatch(whiteClientKey models.Key, blackClientKey models.Key, timeControl *models.TimeControl, result models.MatchResult) *models.Match {
	matchId := uuid.New().String()
	now := time.Now()
	return &models.Match{
		Uuid:                  matchId,
		Board:                 chess.GetInitBoard(),
		WhiteClientKey:        whiteClientKey,
		WhiteTimeRemainingSec: float64(timeControl.InitialTimeSec),
		BlackClientKey:        blackClientKey,
		BlackTimeRemainingSec: float64(timeControl.InitialTimeSec),
		TimeControl:           timeControl,
		LastMoveTime:          &now,
		Result:                result,
	}
}

type MatchBuilder struct {
	match *models.Match
}

func NewMatchBuilder() *MatchBuilder {
	now := time.Now()
	return &MatchBuilder{
		match: &models.Match{
			Uuid:         uuid.New().String(),
			Board:        chess.GetInitBoard(),
			LastMoveTime: &now,
		},
	}
}

func (mb *MatchBuilder) WithUuid(uuid string) *MatchBuilder {
	mb.match.Uuid = uuid
	return mb
}

func (mb *MatchBuilder) WithBoard(board *chess.Board) *MatchBuilder {
	mb.match.Board = board
	if board.Result == chess.BOARD_RESULT_WHITE_WINS_BY_CHECKMATE {
		mb.match.Result = models.MATCH_RESULT_WHITE_WINS_BY_CHECKMATE
	} else if board.Result == chess.BOARD_RESULT_BLACK_WINS_BY_CHECKMATE {
		mb.match.Result = models.MATCH_RESULT_BLACK_WINS_BY_CHECKMATE
	} else if board.Result == chess.BOARD_RESULT_DRAW_BY_STALEMATE {
		mb.match.Result = models.MATCH_RESULT_DRAW_BY_STALEMATE
	} else if board.Result == chess.BOARD_RESULT_DRAW_BY_INSUFFICIENT_MATERIAL {
		mb.match.Result = models.MATCH_RESULT_DRAW_BY_INSUFFICIENT_MATERIAL
	} else if board.Result == chess.BOARD_RESULT_DRAW_BY_THREEFOLD_REPETITION {
		mb.match.Result = models.MATCH_RESULT_DRAW_BY_THREEFOLD_REPETITION
	} else if board.Result == chess.BOARD_RESULT_DRAW_BY_FIFTY_MOVE_RULE {
		mb.match.Result = models.MATCH_RESULT_DRAW_BY_FIFTY_MOVE_RULE
	}
	return mb
}

func (mb *MatchBuilder) WithWhiteClientKey(clientKey models.Key) *MatchBuilder {
	mb.match.WhiteClientKey = clientKey
	return mb
}

func (mb *MatchBuilder) WithWhiteTimeRemainingSec(timeRemainingSec float64) *MatchBuilder {
	mb.match.WhiteTimeRemainingSec = timeRemainingSec
	if timeRemainingSec <= 0 {
		mb.match.Result = models.MATCH_RESULT_BLACK_WINS_BY_TIMEOUT
	}
	return mb
}

func (mb *MatchBuilder) WithBlackClientKey(clientKey models.Key) *MatchBuilder {
	mb.match.BlackClientKey = clientKey
	return mb
}

func (mb *MatchBuilder) WithBlackTimeRemainingSec(timeRemainingSec float64) *MatchBuilder {
	mb.match.BlackTimeRemainingSec = timeRemainingSec
	if timeRemainingSec <= 0 {
		mb.match.Result = models.MATCH_RESULT_WHITE_WINS_BY_TIMEOUT
	}
	return mb
}

func (mb *MatchBuilder) WithTimeRemainingSec(timeRemainingSec float64) *MatchBuilder {
	mb.match.WhiteTimeRemainingSec = timeRemainingSec
	mb.match.BlackTimeRemainingSec = timeRemainingSec
	return mb
}

func (mb *MatchBuilder) WithTimeControl(timeControl *models.TimeControl) *MatchBuilder {
	mb.match.TimeControl = timeControl
	return mb
}

func (mb *MatchBuilder) WithLastMoveTime(lastMoveTime *time.Time) *MatchBuilder {
	mb.match.LastMoveTime = lastMoveTime
	return mb
}

func (mb *MatchBuilder) WithBotName(botName string) *MatchBuilder {
	mb.match.BotName = botName
	return mb
}

func (mb *MatchBuilder) WithClientKeys(clientAKey models.Key, clientBKey models.Key) *MatchBuilder {
	clientAIsWhite := helpers.RandomBool()
	var whiteClientKey, blackClientKey models.Key
	if clientAIsWhite {
		whiteClientKey = clientAKey
		blackClientKey = clientBKey
	} else {
		whiteClientKey = clientBKey
		blackClientKey = clientAKey
	}
	mb.match.WhiteClientKey = whiteClientKey
	mb.match.BlackClientKey = blackClientKey
	return mb
}

func (mb *MatchBuilder) WithResult(result models.MatchResult) *MatchBuilder {
	mb.match.Result = result
	return mb
}

func (mb *MatchBuilder) FromChallenge(challenge *models.Challenge) *MatchBuilder {
	mb.match = NewMatch(challenge.ChallengerKey, challenge.ChallengedKey, challenge.TimeControl, models.MATCH_RESULT_IN_PROGRESS)
	if challenge.IsChallengerWhite {
		mb.WithWhiteClientKey(challenge.ChallengerKey)
		mb.WithBlackClientKey(challenge.ChallengedKey)
	} else if challenge.IsChallengerBlack {
		mb.WithWhiteClientKey(challenge.ChallengedKey)
		mb.WithBlackClientKey(challenge.ChallengerKey)
	} else {
		mb.WithClientKeys(challenge.ChallengerKey, challenge.ChallengedKey)
	}
	mb.WithBotName(challenge.BotName)
	return mb
}

func (mb *MatchBuilder) FromMatch(match *models.Match) *MatchBuilder {
	matchCopy := *match
	mb.match = &matchCopy
	return mb
}

func (mb *MatchBuilder) Build() *models.Match {
	return mb.match
}
