package chess

import (
	"fmt"
	"math"
)

func GetCheckingSquares(board *Board, isWhiteKing bool) *[]*Square {
	checkingSquares := make([]*Square, 0)
	kingSquare := board.GetKingSquare(isWhiteKing)
	knightCheckSquares := []*Square{
		{kingSquare.Rank + 2, kingSquare.File + 1},
		{kingSquare.Rank + 2, kingSquare.File - 1},
		{kingSquare.Rank + 1, kingSquare.File + 2},
		{kingSquare.Rank + 1, kingSquare.File - 2},
		{kingSquare.Rank - 1, kingSquare.File + 2},
		{kingSquare.Rank - 1, kingSquare.File - 2},
		{kingSquare.Rank - 2, kingSquare.File + 1},
		{kingSquare.Rank - 2, kingSquare.File - 1},
	}
	for _, knightCheckSquare := range knightCheckSquares {
		if !knightCheckSquare.IsValidBoardSquare() {
			continue
		}
		piece := board.GetPieceOnSquare(knightCheckSquare)
		if isWhiteKing && piece == BLACK_KNIGHT {
			checkingSquares = append(checkingSquares, knightCheckSquare)
		} else if !isWhiteKing && piece == WHITE_KNIGHT {
			checkingSquares = append(checkingSquares, knightCheckSquare)
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
		tempBoard := *board
		UpdateBoardPiecesFromMove(&tempBoard, move)
		checkingSquares := GetCheckingSquares(&tempBoard, board.IsWhiteTurn)
		if len(*checkingSquares) == 0 {
			filteredMoves = append(filteredMoves, move)
		}
	}
	return &filteredMoves
}

func addKingChecksToMoves(board *Board, moves *[]*Move) {
	for _, move := range *moves {
		tempBoard := *board
		UpdateBoardPiecesFromMove(&tempBoard, move)
		move.KingCheckingSquares = *GetCheckingSquares(&tempBoard, !board.IsWhiteTurn)
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

func GetLegalMovesForKnight(board *Board, square *Square) (*[]*Move, error) {
	knightMoves := make([]*Move, 0)
	piece := board.GetPieceOnSquare(square)
	if board.IsWhiteTurn && piece != WHITE_KNIGHT {
		return nil, fmt.Errorf("square contains unexpected piece %s, expected WHITE_KNIGHT", piece)
	} else if !board.IsWhiteTurn && piece != BLACK_KNIGHT {
		return nil, fmt.Errorf("square contains unexpected piece %s, expected BLACK_KNIGHT", piece)
	}
	landSquares := []*Square{
		{square.Rank + 2, square.File - 1},
		{square.Rank + 2, square.File + 1},
		{square.Rank + 1, square.File - 2},
		{square.Rank + 1, square.File + 2},
		{square.Rank - 1, square.File - 2},
		{square.Rank - 1, square.File + 2},
		{square.Rank - 2, square.File - 1},
		{square.Rank - 2, square.File + 1},
	}
	for _, landSquare := range landSquares {
		if !landSquare.IsValidBoardSquare() {
			continue
		}
		landPiece := board.GetPieceOnSquare(landSquare)
		if landPiece == EMPTY {
			move := Move{piece, square, landSquare, EMPTY, make([]*Square, 0), EMPTY}
			knightMoves = append(knightMoves, &move)
		} else if landPiece.IsWhite() != board.IsWhiteTurn {
			move := Move{piece, square, landSquare, landPiece, make([]*Square, 0), EMPTY}
			knightMoves = append(knightMoves, &move)
		}
	}
	knightMoves = *filterMovesByKingSafety(board, &knightMoves)
	addKingChecksToMoves(board, &knightMoves)
	return &knightMoves, nil
}

func GetLegalMovesForBishop(board *Board, square *Square) (*[]*Move, error) {
	bishopMoves := make([]*Move, 0)
	piece := board.GetPieceOnSquare(square)
	if board.IsWhiteTurn && piece != WHITE_BISHOP {
		return nil, fmt.Errorf("square contains unexpected piece %s, expected WHITE_BISHOP", piece)
	} else if !board.IsWhiteTurn && piece != BLACK_BISHOP {
		return nil, fmt.Errorf("square contains unexpected piece %s, expected BLACK_BISHOP", piece)
	}
	for _, dir := range [4][2]int{{1, 1}, {-1, 1}, {1, -1}, {-1, -1}} {
		dis := 0
		for {
			dis++
			landSquare := Square{square.Rank + uint8(dis*dir[0]), square.File + uint8(dis*dir[1])}
			if !landSquare.IsValidBoardSquare() {
				break
			}
			landPiece := board.GetPieceOnSquare(&landSquare)
			if landPiece == EMPTY {
				move := Move{piece, square, &landSquare, EMPTY, make([]*Square, 0), EMPTY}
				bishopMoves = append(bishopMoves, &move)
			} else {
				if piece.IsWhite() != landPiece.IsWhite() {
					move := Move{piece, square, &landSquare, landPiece, make([]*Square, 0), EMPTY}
					bishopMoves = append(bishopMoves, &move)
				}
				break
			}
		}
	}
	bishopMoves = *filterMovesByKingSafety(board, &bishopMoves)
	addKingChecksToMoves(board, &bishopMoves)
	return &bishopMoves, nil
}

func GetLegalMovesForRook(board *Board, square *Square) (*[]*Move, error) {
	rookMoves := make([]*Move, 0)
	piece := board.GetPieceOnSquare(square)
	if board.IsWhiteTurn && piece != WHITE_ROOK {
		return nil, fmt.Errorf("square contains unexpected piece %s, expected WHITE_ROOK", piece)
	} else if !board.IsWhiteTurn && piece != BLACK_ROOK {
		return nil, fmt.Errorf("square contains unexpected piece %s, expected BLACK_ROOK", piece)
	}
	for _, dir := range [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		dis := 0
		for {
			dis++
			landSquare := Square{square.Rank + uint8(dis*dir[0]), square.File + uint8(dis*dir[1])}
			if !landSquare.IsValidBoardSquare() {
				break
			}
			landPiece := board.GetPieceOnSquare(&landSquare)
			if landPiece == EMPTY {
				move := Move{piece, square, &landSquare, EMPTY, make([]*Square, 0), EMPTY}
				rookMoves = append(rookMoves, &move)
			} else {
				if piece.IsWhite() != landPiece.IsWhite() {
					move := Move{piece, square, &landSquare, landPiece, make([]*Square, 0), EMPTY}
					rookMoves = append(rookMoves, &move)
				}
				break
			}
		}
	}
	rookMoves = *filterMovesByKingSafety(board, &rookMoves)
	addKingChecksToMoves(board, &rookMoves)
	return &rookMoves, nil
}

func GetLegalMovesForQueen(board *Board, square *Square) (*[]*Move, error) {
	queenMoves := make([]*Move, 0)
	piece := board.GetPieceOnSquare(square)
	if board.IsWhiteTurn && piece != WHITE_QUEEN {
		return nil, fmt.Errorf("square contains unexpected piece %s, expected WHITE_ROOK", piece)
	} else if !board.IsWhiteTurn && piece != BLACK_QUEEN {
		return nil, fmt.Errorf("square contains unexpected piece %s, expected BLACK_ROOK", piece)
	}
	for _, dir := range [8][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}} {
		dis := 0
		for {
			dis++
			landSquare := Square{square.Rank + uint8(dis*dir[0]), square.File + uint8(dis*dir[1])}
			if !landSquare.IsValidBoardSquare() {
				break
			}
			landPiece := board.GetPieceOnSquare(&landSquare)
			if landPiece == EMPTY {
				move := Move{piece, square, &landSquare, EMPTY, make([]*Square, 0), EMPTY}
				queenMoves = append(queenMoves, &move)
			} else {
				if piece.IsWhite() != landPiece.IsWhite() {
					move := Move{piece, square, &landSquare, landPiece, make([]*Square, 0), EMPTY}
					queenMoves = append(queenMoves, &move)
				}
				break
			}
		}
	}
	queenMoves = *filterMovesByKingSafety(board, &queenMoves)
	addKingChecksToMoves(board, &queenMoves)
	return &queenMoves, nil
}

func GetLegalMovesForKing(board *Board, square *Square) (*[]*Move, error) {
	kingMoves := make([]*Move, 0)
	piece := board.GetPieceOnSquare(square)
	if board.IsWhiteTurn && piece != WHITE_KING {
		return nil, fmt.Errorf("unexpected piece on square %s, expected WHITE_KING", piece)
	} else if !board.IsWhiteTurn && piece != BLACK_KING {
		return nil, fmt.Errorf("unexpected piece on square %s, expected BLACK_KING", piece)
	}
	landSquares := []*Square{
		{square.Rank + 1, square.File - 1},
		{square.Rank + 1, square.File},
		{square.Rank + 1, square.File + 1},
		{square.Rank, square.File - 1},
		{square.Rank, square.File + 1},
		{square.Rank - 1, square.File - 1},
		{square.Rank - 1, square.File},
		{square.Rank - 1, square.File + 1},
	}
	enemyKingSquare := board.GetKingSquare(!board.IsWhiteTurn)
	for _, landSquare := range landSquares {
		if !landSquare.IsValidBoardSquare() {
			continue
		}
		enemyKingRankDis := math.Abs(float64(int(enemyKingSquare.Rank) - int(landSquare.Rank)))
		enemyKingFileDis := math.Abs(float64(int(enemyKingSquare.File) - int(landSquare.File)))
		if enemyKingRankDis < 2 && enemyKingFileDis < 2 {
			continue
		}
		landPiece := board.GetPieceOnSquare(landSquare)
		if landPiece == EMPTY || landPiece.IsWhite() != piece.IsWhite() {
			move := Move{piece, square, landSquare, landPiece, make([]*Square, 0), EMPTY}
			kingMoves = append(kingMoves, &move)
		}
	}
	if canCastleKingside(board, square) {
		kingDestSquare := Square{square.Rank, square.File + 2}
		kingMoves = append(kingMoves, &Move{piece, square, &kingDestSquare, EMPTY, make([]*Square, 0), EMPTY})
	}
	if canCastleQueenside(board, square) {
		kingDestSquare := Square{square.Rank, square.File - 2}
		kingMoves = append(kingMoves, &Move{piece, square, &kingDestSquare, EMPTY, make([]*Square, 0), EMPTY})
	}
	kingMoves = *filterMovesByKingSafety(board, &kingMoves)
	return &kingMoves, nil
}

func canCastleKingside(board *Board, square *Square) bool {
	if board.IsWhiteTurn && !board.CanWhiteCastleKingside {
		return false
	} else if !board.IsWhiteTurn && !board.CanBlackCastleKingside {
		return false
	}
	piece := board.GetPieceOnSquare(square)

	kingRightSquare := Square{square.Rank, square.File + 1}
	kingRightTwoSquare := Square{square.Rank, square.File + 2}
	kingRightPiece := board.GetPieceOnSquare(&kingRightSquare)
	kingTwoRightPiece := board.GetPieceOnSquare(&kingRightTwoSquare)
	if kingRightPiece != EMPTY || kingTwoRightPiece != EMPTY {
		return false
	}

	tempBoard := *board
	tempBoard.SetPieceOnSquare(piece, &kingRightSquare)
	tempBoard.SetPieceOnSquare(EMPTY, square)
	if len(*GetCheckingSquares(board, board.IsWhiteTurn)) > 0 {
		return false
	}
	tempBoard.SetPieceOnSquare(piece, &kingRightTwoSquare)
	tempBoard.SetPieceOnSquare(EMPTY, &kingRightSquare)
	if len(*GetCheckingSquares(board, board.IsWhiteTurn)) > 0 {
		return false
	}
	return true
}

func canCastleQueenside(board *Board, square *Square) bool {
	if board.IsWhiteTurn && !board.CanWhiteCastleQueenside {
		return false
	} else if !board.IsWhiteTurn && !board.CanBlackCastleQueenside {
		return false
	}
	piece := board.GetPieceOnSquare(square)

	kingLeftSquare := Square{square.Rank, square.File - 1}
	kingLeftTwoSquare := Square{square.Rank, square.File - 2}
	kingLeftPiece := board.GetPieceOnSquare(&kingLeftSquare)
	kingTwoLeftPiece := board.GetPieceOnSquare(&kingLeftTwoSquare)
	if kingLeftPiece != EMPTY || kingTwoLeftPiece != EMPTY {
		return false
	}

	tempBoard := *board
	tempBoard.SetPieceOnSquare(piece, &kingLeftSquare)
	tempBoard.SetPieceOnSquare(EMPTY, square)
	if len(*GetCheckingSquares(board, board.IsWhiteTurn)) > 0 {
		return false
	}
	tempBoard.SetPieceOnSquare(piece, &kingLeftTwoSquare)
	tempBoard.SetPieceOnSquare(EMPTY, &kingLeftSquare)
	if len(*GetCheckingSquares(board, board.IsWhiteTurn)) > 0 {
		return false
	}
	return true
}

func UpdateBoardFromMove(board *Board, move *Move) {
	UpdateBoardPiecesFromMove(board, move)
	if move.CapturedPiece != EMPTY || move.Piece.IsPawn() {
		board.HalfMoveClockCount = 0
	} else {
		board.HalfMoveClockCount++
	}

	if move.DoesAllowEnPassant() {
		board.OptEnPassantSquare = &Square{
			uint8(math.Min(float64(move.StartSquare.Rank), float64(move.EndSquare.Rank))) + 1,
			move.StartSquare.File,
		}
	}

	if !board.IsWhiteTurn {
		board.FullMoveCount++
	}
	board.IsWhiteTurn = !board.IsWhiteTurn

	UpdateCastleRightsFromMove(board, move)
}

func UpdateBoardPiecesFromMove(board *Board, move *Move) {
	movingPiece := board.GetPieceOnSquare(move.StartSquare)
	var landingPiece Piece
	if move.PawnUpgradedTo != EMPTY {
		landingPiece = move.PawnUpgradedTo
	} else {
		landingPiece = movingPiece
	}
	board.SetPieceOnSquare(landingPiece, move.EndSquare)
	board.SetPieceOnSquare(EMPTY, move.StartSquare)
	if board.OptEnPassantSquare != nil && movingPiece.IsPawn() && move.EndSquare.EqualTo(board.OptEnPassantSquare) {
		enPassantedPawnSquare := Square{
			move.StartSquare.Rank,
			board.OptEnPassantSquare.File,
		}
		board.SetPieceOnSquare(EMPTY, &enPassantedPawnSquare)
	}
	if move.Piece == WHITE_KING {
		board.optWhiteKingSquare = move.EndSquare
		if move.EndSquare.EqualTo(&Square{1, 7}) {
			board.SetPieceOnSquare(WHITE_ROOK, &Square{1, 6})
			board.SetPieceOnSquare(EMPTY, &Square{1, 8})
		} else if move.EndSquare.EqualTo(&Square{1, 3}) {
			board.SetPieceOnSquare(WHITE_ROOK, &Square{1, 4})
			board.SetPieceOnSquare(EMPTY, &Square{1, 1})
		}
	} else if move.Piece == BLACK_KING {
		board.optWhiteKingSquare = move.EndSquare
		if move.EndSquare.EqualTo(&Square{8, 7}) {
			board.SetPieceOnSquare(BLACK_ROOK, &Square{8, 6})
			board.SetPieceOnSquare(EMPTY, &Square{8, 8})
		} else if move.EndSquare.EqualTo(&Square{8, 3}) {
			board.SetPieceOnSquare(BLACK_ROOK, &Square{8, 4})
			board.SetPieceOnSquare(EMPTY, &Square{8, 1})
		}
	}
}

func UpdateCastleRightsFromMove(board *Board, move *Move) {
	if !move.Piece.IsKing() && !move.Piece.IsRook() {
		return
	}
	if move.Piece.IsRook() {
		if move.StartSquare.EqualTo(&Square{1, 1}) {
			board.CanWhiteCastleQueenside = false
		} else if move.StartSquare.EqualTo(&Square{1, 8}) {
			board.CanWhiteCastleKingside = false
		} else if move.StartSquare.EqualTo(&Square{8, 1}) {
			board.CanBlackCastleQueenside = false
		} else if move.StartSquare.EqualTo(&Square{8, 8}) {
			board.CanBlackCastleKingside = false
		}
	} else if move.Piece.IsKing() {
		if move.Piece.IsWhite() {
			board.CanWhiteCastleKingside = false
			board.CanWhiteCastleQueenside = false
		} else {
			board.CanBlackCastleKingside = false
			board.CanBlackCastleQueenside = false
		}
	}
}
