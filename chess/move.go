package chess

type Move struct {
	piece               Piece
	StartSquare         *Square
	EndSquare           *Square
	CapturedPiece       *Piece
	KingCheckingSquares []Square
}
