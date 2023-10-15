package chess

type Piece uint8

const (
	EMPTY Piece = iota
	WHITE_PAWN
	WHITE_KNIGHT
	WHITE_BISHOP
	WHITE_ROOK
	WHITE_QUEEN
	WHITE_KING
	BLACK_PAWN
	BLACK_KNIGHT
	BLACK_BISHOP
	BLACK_ROOK
	BLACK_QUEEN
	BLACK_KING
)

func (p Piece) String() string {
	return []string{"EMPTY", "WHITE_PAWN", "WHITE_KNIGHT", "WHITE_BISHOP", "WHITE_ROOK",
		"WHITE_QUEEN", "WHITE_KING", "BLACK_PAWN", "BLACK_KNIGHT", "BLACK_BISHOP",
		"BLACK_ROOK", "BLACK_QUEEN", "BLACK_KING",
	}[p]
}

func (p Piece) IsWhite() bool {
	return p == WHITE_PAWN || p == WHITE_KNIGHT || p == WHITE_BISHOP ||
		p == WHITE_ROOK || p == WHITE_QUEEN || p == WHITE_KING
}

func (p Piece) IsPawn() bool {
	return p == WHITE_PAWN || p == BLACK_PAWN
}

func (p Piece) IsKnight() bool {
	return p == WHITE_KNIGHT || p == BLACK_KNIGHT
}

func (p Piece) IsBishop() bool {
	return p == WHITE_BISHOP || p == BLACK_BISHOP
}

func (p Piece) IsRook() bool {
	return p == WHITE_ROOK || p == BLACK_ROOK
}

func (p Piece) IsQueen() bool {
	return p == WHITE_QUEEN || p == BLACK_QUEEN
}

func (p Piece) IsKing() bool {
	return p == WHITE_KING || p == BLACK_KING
}
