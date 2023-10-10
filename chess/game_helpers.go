package chess

func GetCheckingSquares(board *Board, isWhiteKing bool) *[]*Square {
	checkingSquares := make([]*Square, 0)
	kingSquare := board.GetKingSquare(isWhiteKing)
	knightCheckSquares := []Square{
		Square{kingSquare.Rank + 2, kingSquare.File + 1},
		Square{kingSquare.Rank + 2, kingSquare.File - 1},
		Square{kingSquare.Rank + 1, kingSquare.File + 2},
		Square{kingSquare.Rank + 1, kingSquare.File - 2},
		Square{kingSquare.Rank - 1, kingSquare.File + 2},
		Square{kingSquare.Rank - 1, kingSquare.File - 2},
		Square{kingSquare.Rank - 2, kingSquare.File + 1},
		Square{kingSquare.Rank - 2, kingSquare.File - 1},
	}
	for _, knightCheckSquare := range knightCheckSquares {
		if !knightCheckSquare.IsValidBoardSquare() {
			continue
		}
		piece := board.GetPieceOnSquare(&knightCheckSquare)
		if isWhiteKing && piece == BLACK_KNIGHT {
			checkingSquares = append(checkingSquares, &knightCheckSquare)
		} else if !isWhiteKing && piece == WHITE_KNIGHT {
			checkingSquares = append(checkingSquares, &knightCheckSquare)
		}
	}
	var pawnCheckSquares []Square
	if isWhiteKing {
		pawnCheckSquares = []Square{
			{kingSquare.Rank + 1, kingSquare.File - 1},
			{kingSquare.Rank + 1, kingSquare.File + 1},
		}
	} else {
		pawnCheckSquares = []Square{
			{kingSquare.Rank - 1, kingSquare.File - 1},
			{kingSquare.Rank - 1, kingSquare.File + 1},
		}
	}
	for _, pawnCheckSquare := range pawnCheckSquares {
		if !pawnCheckSquare.IsValidBoardSquare() {
			continue
		}
		piece := board.GetPieceOnSquare(&pawnCheckSquare)
		if isWhiteKing && piece == BLACK_PAWN {
			checkingSquares = append(checkingSquares, &pawnCheckSquare)
		} else if !isWhiteKing && piece == WHITE_PAWN {
			checkingSquares = append(checkingSquares, &pawnCheckSquare)
		}
	}
	for _, diagDir := range [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}} {
		for dis := 1; dis < 8; dis++ {
			diagSquare := Square{
				kingSquare.Rank + uint8(dis*diagDir[0]),
				kingSquare.File + uint8(dis*diagDir[1])}
			if !diagSquare.IsValidBoardSquare() {
				break
			}
			piece := board.GetPieceOnSquare(&diagSquare)
			if isWhiteKing {
				if piece == BLACK_BISHOP || piece == BLACK_QUEEN {
					checkingSquares = append(checkingSquares, &diagSquare)
				}
			} else {
				if piece == WHITE_BISHOP || piece == WHITE_QUEEN {
					checkingSquares = append(checkingSquares, &diagSquare)
				}
			}
			if piece != EMPTY {
				break
			}
		}
	}
	for _, straightDir := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		for dis := 1; dis < 8; dis++ {
			straightSquare := Square{
				kingSquare.Rank + uint8(dis*straightDir[0]),
				kingSquare.File + uint8(dis*straightDir[1]),
			}
			if !straightSquare.IsValidBoardSquare() {
				break
			}
			piece := board.GetPieceOnSquare(&straightSquare)
			if isWhiteKing {
				if piece == BLACK_ROOK || piece == BLACK_QUEEN {
					checkingSquares = append(checkingSquares, &straightSquare)
				}
			} else {
				if piece == WHITE_ROOK || piece == WHITE_QUEEN {
					checkingSquares = append(checkingSquares, &straightSquare)
				}
			}
			if piece != EMPTY {
				break
			}
		}
	}
	return &checkingSquares
}

func GetResultingBoard(startBoard Board, move Move) {

}
