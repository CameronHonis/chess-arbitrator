package chess

type BoardBuilder struct {
	board *Board
}

func NewBoardBuilder() *BoardBuilder {
	board := Board{}
	board.RepetitionsByMiniFEN = make(map[string]uint8)
	return &BoardBuilder{
		board: &board,
	}
}

func (bb *BoardBuilder) WithPieces(pieces *[8][8]Piece) *BoardBuilder {
	bb.board.Pieces = *pieces
	return bb
}

func (bb *BoardBuilder) WithEnPassantSquare(enPassantSquare *Square) *BoardBuilder {
	bb.board.OptEnPassantSquare = enPassantSquare
	return bb
}

func (bb *BoardBuilder) WithIsWhiteTurn(isWhiteTurn bool) *BoardBuilder {
	bb.board.IsWhiteTurn = isWhiteTurn
	return bb
}

func (bb *BoardBuilder) WithCanWhiteCastleQueenside(canWhiteCastleQueenside bool) *BoardBuilder {
	bb.board.CanWhiteCastleQueenside = canWhiteCastleQueenside
	return bb
}

func (bb *BoardBuilder) WithCanWhiteCastleKingside(canWhiteCastleKingside bool) *BoardBuilder {
	bb.board.CanWhiteCastleKingside = canWhiteCastleKingside
	return bb
}
func (bb *BoardBuilder) WithCanBlackCastleQueenside(canBlackCastleQueenside bool) *BoardBuilder {
	bb.board.CanBlackCastleQueenside = canBlackCastleQueenside
	return bb
}

func (bb *BoardBuilder) WithCanBlackCastleKingside(canBlackCastleKingside bool) *BoardBuilder {
	bb.board.CanBlackCastleKingside = canBlackCastleKingside
	return bb
}

func (bb *BoardBuilder) WithHalfMoveClockCount(halfMoveClockCount uint8) *BoardBuilder {
	bb.board.HalfMoveClockCount = halfMoveClockCount
	return bb
}

func (bb *BoardBuilder) WithFullMoveCount(fullMoveCount uint16) *BoardBuilder {
	bb.board.FullMoveCount = fullMoveCount
	return bb
}

func (bb *BoardBuilder) WithRepetitionsByMiniFEN(repetitionsByMiniFEN map[string]uint8) *BoardBuilder {
	bb.board.RepetitionsByMiniFEN = repetitionsByMiniFEN
	return bb
}

func (bb *BoardBuilder) WithIsTerminal(isTerminal bool) *BoardBuilder {
	bb.board.IsTerminal = isTerminal
	return bb
}

func (bb *BoardBuilder) WithIsWhiteWinner(isWhiteWinner bool) *BoardBuilder {
	bb.board.IsTerminal = isWhiteWinner
	return bb
}
func (bb *BoardBuilder) WithIsBlackWinner(isBlackWinner bool) *BoardBuilder {
	bb.board.IsTerminal = isBlackWinner
	return bb
}
func (bb *BoardBuilder) Build() *Board {
	return bb.board
}
