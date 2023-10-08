package chess

type BoardBuilder struct {
	board *Board
}

func NewBoardBuilder() *BoardBuilder {
	return &BoardBuilder{
		board: &Board{},
	}
}

func (bb *BoardBuilder) WithPieces(pieces *[8][8]Piece) *BoardBuilder {
	bb.board.Pieces = pieces
	return bb
}

func (bb *BoardBuilder) WithEnPassantSquare(enPassantSquare *Square) *BoardBuilder {
	bb.board.EnPassantSquare = enPassantSquare
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

func (bb *BoardBuilder) Build() *Board {
	return bb.board
}