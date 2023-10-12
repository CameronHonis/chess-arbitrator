package chess

import "fmt"

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

func filterMovesByKingSafety(board *Board, moves *[]*Move) *[]*Move {
	filteredMoves := make([]*Move, 0, len(*moves))
	for _, move := range *moves {
		tempBoard := board.CopyWithPieces()
		UpdateBoardPiecesFromMove(tempBoard, move)
		checkingSquares := GetCheckingSquares(tempBoard, board.IsWhiteTurn)
		if len(*checkingSquares) == 0 {
			filteredMoves = append(filteredMoves, move)
		}
	}
	return &filteredMoves
}

func addKingChecksToMoves(board *Board, moves *[]*Move) {
	for _, move := range *moves {
		move.KingCheckingSquares = *GetCheckingSquares(board, !board.IsWhiteTurn)
	}
}
func GetLegalMovesForPawn(board *Board, square *Square) (*[]*Move, error) {
	pawnMoves := make([]*Move, 0)
	piece := board.GetPieceOnSquare(square)
	var upgradePieces [4]Piece
	var squareInFront Square
	var leftCaptureSquare Square
	var rightCaptureSquare Square
	if board.IsWhiteTurn {
		if piece != WHITE_PAWN {
			return nil, fmt.Errorf("square contains unexpected piece %s, expected WHITE_PAWN", piece)
		}
		upgradePieces = [4]Piece{WHITE_KNIGHT, WHITE_BISHOP, WHITE_ROOK, WHITE_QUEEN}
		squareInFront = Square{square.Rank + 1, square.File}
		leftCaptureSquare = Square{square.Rank + 1, square.File - 1}
		rightCaptureSquare = Square{square.Rank + 1, square.File + 1}
	} else {
		if piece != BLACK_PAWN {
			return nil, fmt.Errorf("square contains unexpected piece %s, expected BLACK_PAWN", piece)
		}
		upgradePieces = [4]Piece{BLACK_KNIGHT, BLACK_BISHOP, BLACK_ROOK, BLACK_QUEEN}
		squareInFront = Square{square.Rank - 1, square.File}
		leftCaptureSquare = Square{square.Rank - 1, square.File - 1}
		rightCaptureSquare = Square{square.Rank - 1, square.File + 1}
	}
	pieceInFront := board.GetPieceOnSquare(&squareInFront)
	if pieceInFront == EMPTY {
		if squareInFront.Rank == 8 || squareInFront.Rank == 1 {
			move0 := Move{piece, square, &squareInFront, EMPTY, make([]*Square, 0), upgradePieces[0]}
			move1 := Move{piece, square, &squareInFront, EMPTY, make([]*Square, 0), upgradePieces[1]}
			move2 := Move{piece, square, &squareInFront, EMPTY, make([]*Square, 0), upgradePieces[2]}
			move3 := Move{piece, square, &squareInFront, EMPTY, make([]*Square, 0), upgradePieces[3]}
			pawnMoves = append(pawnMoves, &move0, &move1, &move2, &move3)
		} else {
			move := Move{piece, square, &squareInFront, EMPTY, make([]*Square, 0), EMPTY}
			pawnMoves = append(pawnMoves, &move)
		}
		if (board.IsWhiteTurn && square.Rank == 2) || (!board.IsWhiteTurn && square.Rank == 7) {
			var squareTwoInFront Square
			if board.IsWhiteTurn {
				squareTwoInFront = Square{square.Rank + 2, square.File}
			} else {
				squareTwoInFront = Square{square.Rank - 2, square.File}
			}
			pieceTwoInFront := board.GetPieceOnSquare(&squareTwoInFront)
			if pieceTwoInFront == EMPTY {
				move := Move{piece, square, &squareTwoInFront, EMPTY, make([]*Square, 0), EMPTY}
				pawnMoves = append(pawnMoves, &move)
			}
		}
	}
	if leftCaptureSquare.IsValidBoardSquare() {
		var leftCapturePiece Piece
		if board.OptEnPassantSquare != nil && leftCaptureSquare.EqualTo(board.OptEnPassantSquare) {
			leftCapturePiece = board.GetPieceOnSquare(&Square{square.Rank, square.File - 1})
		} else {
			leftCapturePiece = board.GetPieceOnSquare(&leftCaptureSquare)
		}
		if leftCapturePiece != EMPTY && leftCapturePiece.IsWhite() != piece.IsWhite() {
			if (piece.IsWhite() && square.Rank == 7) || (!piece.IsWhite() && square.Rank == 2) {
				move0 := Move{piece, square, &leftCaptureSquare, leftCapturePiece, make([]*Square, 0), upgradePieces[0]}
				move1 := Move{piece, square, &leftCaptureSquare, leftCapturePiece, make([]*Square, 0), upgradePieces[1]}
				move2 := Move{piece, square, &leftCaptureSquare, leftCapturePiece, make([]*Square, 0), upgradePieces[2]}
				move3 := Move{piece, square, &leftCaptureSquare, leftCapturePiece, make([]*Square, 0), upgradePieces[3]}
				pawnMoves = append(pawnMoves, &move0, &move1, &move2, &move3)
			} else {
				move := Move{piece, square, &leftCaptureSquare, leftCapturePiece, make([]*Square, 0), EMPTY}
				pawnMoves = append(pawnMoves, &move)
			}
		}
	}
	if rightCaptureSquare.IsValidBoardSquare() {
		var rightCapturePiece Piece
		if board.OptEnPassantSquare != nil && rightCaptureSquare.EqualTo(board.OptEnPassantSquare) {
			rightCapturePiece = board.GetPieceOnSquare(&Square{square.Rank, square.File + 1})
		} else {
			rightCapturePiece = board.GetPieceOnSquare(&rightCaptureSquare)
		}
		if rightCapturePiece != EMPTY && rightCapturePiece.IsWhite() != piece.IsWhite() {
			if (piece.IsWhite() && square.Rank == 7) || (!piece.IsWhite() && square.Rank == 2) {
				move0 := Move{piece, square, &rightCaptureSquare, rightCapturePiece, make([]*Square, 0), upgradePieces[0]}
				move1 := Move{piece, square, &rightCaptureSquare, rightCapturePiece, make([]*Square, 0), upgradePieces[1]}
				move2 := Move{piece, square, &rightCaptureSquare, rightCapturePiece, make([]*Square, 0), upgradePieces[2]}
				move3 := Move{piece, square, &rightCaptureSquare, rightCapturePiece, make([]*Square, 0), upgradePieces[3]}
				pawnMoves = append(pawnMoves, &move0, &move1, &move2, &move3)
			} else {
				move := Move{piece, square, &rightCaptureSquare, rightCapturePiece, make([]*Square, 0), EMPTY}
				pawnMoves = append(pawnMoves, &move)
			}
		}
	}
	pawnMoves = *filterMovesByKingSafety(board, &pawnMoves)
	addKingChecksToMoves(board, &pawnMoves)
	return &pawnMoves, nil
}

func UpdateBoardFromMove(startBoard *Board, move *Move) {
}

func UpdateBoardPiecesFromMove(board *Board, move *Move) {
	movingPiece := board.GetPieceOnSquare(move.StartSquare)
	board.SetPieceOnSquare(movingPiece, move.EndSquare)
	board.SetPieceOnSquare(EMPTY, move.StartSquare)
	if board.OptEnPassantSquare != nil && movingPiece.IsPawn() && move.EndSquare.EqualTo(board.OptEnPassantSquare) {
		enPassantedPawnSquare := Square{
			move.StartSquare.Rank,
			board.OptEnPassantSquare.File,
		}
		board.SetPieceOnSquare(EMPTY, &enPassantedPawnSquare)
	}
}
