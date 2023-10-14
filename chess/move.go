package chess

import (
	"math"
	"sort"
)

type Move struct {
	Piece               Piece
	StartSquare         *Square
	EndSquare           *Square
	CapturedPiece       Piece
	KingCheckingSquares []*Square
	PawnUpgradedTo      Piece
}

func (move *Move) SortKingCheckingSquares() {
	sort.Slice(move.KingCheckingSquares, func(i, j int) bool {
		squareA := move.KingCheckingSquares[i]
		squareB := move.KingCheckingSquares[j]
		return squareA.Rank*8+squareA.File < squareB.Rank*8+squareB.File
	})
}

func (move *Move) EqualTo(otherMove *Move) bool {
	if len(move.KingCheckingSquares) != len(otherMove.KingCheckingSquares) {
		return false
	}
	move.SortKingCheckingSquares()
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

func (move *Move) DoesAllowEnPassant() bool {
	if !move.Piece.IsPawn() {
		return false
	}
	dis := math.Abs(float64(int(move.EndSquare.Rank) - int(move.StartSquare.Rank)))
	return dis > 1
}
