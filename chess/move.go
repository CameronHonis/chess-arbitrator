package chess

type Move struct {
	Piece               Piece
	StartSquare         *Square
	EndSquare           *Square
	CapturedPiece       Piece
	KingCheckingSquares []*Square
	PawnUpgradedTo      Piece
}

func (move *Move) EqualTo(otherMove *Move) bool {
	for squareIdx, square := range move.KingCheckingSquares {
		otherSquare := otherMove.KingCheckingSquares[squareIdx]
		if !square.EqualTo(otherSquare) {
			return false
		}
	}
	return move.Piece == otherMove.Piece &&
		move.StartSquare.EqualTo(otherMove.StartSquare) &&
		move.EndSquare.EqualTo(otherMove.EndSquare) &&
		move.CapturedPiece == otherMove.CapturedPiece &&
		move.PawnUpgradedTo == otherMove.PawnUpgradedTo
}
