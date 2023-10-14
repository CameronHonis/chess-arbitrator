package chess

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Board struct {
	Pieces                  [8][8]Piece
	OptEnPassantSquare      *Square
	IsWhiteTurn             bool
	CanWhiteCastleQueenside bool
	CanWhiteCastleKingside  bool
	CanBlackCastleQueenside bool
	CanBlackCastleKingside  bool
	HalfMoveClockCount      uint8
	FullMoveCount           uint16
	RepetitionsByMiniFEN    map[string]uint8
	IsTerminal              bool
	IsWhiteWinner           bool
	IsBlackWinner           bool
	// memoizers
	optMaterialCount   *MaterialCount
	optWhiteKingSquare *Square
	optBlackKingSquare *Square
}

func NewBoard(pieces *[8][8]Piece,
	enPassantSquare *Square,
	isWhiteTurn bool,
	canWhiteCastleKingside bool,
	canWhiteCastleQueenside bool,
	canBlackCastleKingside bool,
	canBlackCastleQueenside bool,
	halfMoveClockCount uint8,
	fullMoveCount uint16,
	repetitionsByMiniFEN map[string]uint8,
	isTerminal bool,
	isWhiteWinner bool,
	isBlackWinner bool) *Board {
	return &Board{
		*pieces, enPassantSquare, isWhiteTurn,
		canWhiteCastleQueenside, canWhiteCastleKingside,
		canBlackCastleQueenside, canBlackCastleKingside,
		halfMoveClockCount, fullMoveCount, repetitionsByMiniFEN, isTerminal, isWhiteWinner, isBlackWinner, nil, nil, nil,
	}
}

func BoardFromFEN(fen string) (*Board, error) {
	pieceByFENrune := map[rune]Piece{
		'p': BLACK_PAWN,
		'n': BLACK_KNIGHT,
		'b': BLACK_BISHOP,
		'r': BLACK_ROOK,
		'q': BLACK_QUEEN,
		'k': BLACK_KING,
		'P': WHITE_PAWN,
		'N': WHITE_KNIGHT,
		'B': WHITE_BISHOP,
		'R': WHITE_ROOK,
		'Q': WHITE_QUEEN,
		'K': WHITE_KING,
	}
	pieces := [8][8]Piece{
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
	}
	boardBuilder := NewBoardBuilder()
	fenSegs := strings.Split(fen, " ")
	if len(fenSegs) != 6 {
		return nil, fmt.Errorf("invalid FEN: wrong number of FEN segments. Expected 6 vs. actual %d", len(fenSegs))
	}
	for fenSegIdx, fenSeg := range fenSegs {
		if fenSegIdx == 0 {
			rank := uint8(8)
			file := uint8(1)
			for _, FENRune := range []rune(fenSeg) {
				if rank < 1 {
					return nil, fmt.Errorf("invalid FEN: too many rows")
				}
				if FENRune == '/' {
					if file < 9 {
						return nil, fmt.Errorf("invalid FEN: not enough files at rank: %d", rank)
					}
					rank--
					file = 1
					continue
				}
				if file > 8 {
					return nil, fmt.Errorf("invalid FEN: too many files at rank: %d", rank)
				}
				if unicode.IsDigit(FENRune) {
					file += uint8(FENRune) - 48
				} else {
					piece, exists := pieceByFENrune[FENRune]
					if !exists {
						coords := (&Square{rank, file}).ToAlgebraicCoords()
						return nil, fmt.Errorf("invalid FEN: unknown piece char '%c' at %s", FENRune, coords)
					}
					pieces[rank-1][file-1] = piece
					file++
				}
			}
			if rank > 1 {
				return nil, fmt.Errorf("invalid FEN: not enough rows")
			}
			boardBuilder.WithPieces(&pieces)
		} else if fenSegIdx == 1 {
			if fenSeg == "w" {
				boardBuilder.WithIsWhiteTurn(true)
			} else if fenSeg == "b" {
				boardBuilder.WithIsWhiteTurn(false)
			} else {
				return nil, fmt.Errorf("invalid FEN: unknown turn specifier %s", fenSeg)
			}
		} else if fenSegIdx == 2 {
			if fenSeg == "-" || fenSeg == "_" {
				continue
			}
			for _, castleRightsRune := range []rune(fenSeg) {
				if castleRightsRune == 'K' {
					boardBuilder.WithCanWhiteCastleKingside(true)
				} else if castleRightsRune == 'Q' {
					boardBuilder.WithCanWhiteCastleQueenside(true)
				} else if castleRightsRune == 'k' {
					boardBuilder.WithCanBlackCastleKingside(true)
				} else if castleRightsRune == 'q' {
					boardBuilder.WithCanBlackCastleQueenside(true)
				} else {
					return nil, fmt.Errorf("invalid FEN: unknown castle rights specifier, got %c", castleRightsRune)
				}
			}
		} else if fenSegIdx == 3 {
			if fenSeg == "-" || fenSeg == "_" {
				continue
			}
			enPassantSquare, err := SquareFromAlgebraicCoords(fenSeg)
			if err != nil {
				return nil, err
			}
			boardBuilder.WithEnPassantSquare(enPassantSquare)
		} else if fenSegIdx == 4 {
			halfMoveClockCount, intErr := strconv.Atoi(fenSeg)
			if intErr != nil {
				err := fmt.Errorf("invalid FEN: could not parse half move clock count, got error: %w", intErr)
				return nil, err
			}
			if halfMoveClockCount < 0 || halfMoveClockCount > 255 {
				err := fmt.Errorf("invalid FEN: half move clock count outside expected range [0, 255], got (%d)", halfMoveClockCount)
				return nil, err
			}
			boardBuilder.WithHalfMoveClockCount(uint8(halfMoveClockCount))
		} else if fenSegIdx == 5 {
			fullMoveCount, intErr := strconv.Atoi(fenSeg)
			if intErr != nil {
				err := fmt.Errorf("invalid FEN: could not parse full move count, got error: %w", intErr)
				return nil, err
			}
			if fullMoveCount < 0 || fullMoveCount > 65535 {
				err := fmt.Errorf("invalid FEN: full move count outside expected range [0, 65535], got (%d)", fullMoveCount)
				return nil, err
			}
			boardBuilder.WithFullMoveCount(uint16(fullMoveCount))
		}
	}
	return boardBuilder.Build(), nil
}

func GetInitBoard() *Board {
	miniFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -"
	repetitionsByMiniFEN := map[string]uint8{
		miniFEN: uint8(1),
	}
	return NewBoard(&[8][8]Piece{
		{WHITE_ROOK, WHITE_KNIGHT, WHITE_BISHOP, WHITE_QUEEN, WHITE_KING, WHITE_BISHOP, WHITE_KNIGHT, WHITE_ROOK},
		{WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN},
		{BLACK_ROOK, BLACK_KNIGHT, BLACK_BISHOP, BLACK_QUEEN, BLACK_KING, BLACK_BISHOP, BLACK_KNIGHT, BLACK_ROOK},
	}, nil, true, true, true, true,
		true, 0, 1, repetitionsByMiniFEN, false, false, false)
}

func (board *Board) SetPieceOnSquare(piece Piece, square *Square) *Board {
	board.Pieces[square.Rank-1][square.File-1] = piece
	return board
}

func (board *Board) GetPieceOnSquare(square *Square) Piece {
	return board.Pieces[square.Rank-1][square.File-1]
}

func (board *Board) IsForcedDraw() bool {
	if board.optMaterialCount == nil {
		board.ComputeMaterialCount(true)
	}
	mat := board.optMaterialCount
	if mat.WhitePawnCount > 0 || mat.BlackPawnCount > 0 {
		return false
	}
	if mat.WhiteQueenCount > 0 || mat.BlackQueenCount > 0 {
		return false
	}
	if mat.WhiteRookCount > 0 || mat.BlackRookCount > 0 {
		return false
	}
	if mat.WhiteKnightCount > 1 || mat.BlackKnightCount > 1 {
		return false
	}

	if mat.WhiteLightBishopCount > 0 && mat.WhiteDarkBishopCount > 0 {
		return false
	}
	whiteBishopCount := mat.WhiteLightBishopCount + mat.WhiteDarkBishopCount
	if whiteBishopCount > 0 && mat.WhiteKnightCount > 0 {
		return false
	}

	if mat.BlackLightBishopCount > 0 && mat.BlackDarkBishopCount > 0 {
		return false
	}
	blackBishopCount := mat.BlackLightBishopCount + mat.BlackDarkBishopCount
	if blackBishopCount > 0 && mat.BlackKnightCount > 0 {
		return false
	}
	return true
}

func (board *Board) ComputeMaterialCount(includeKings bool) *MaterialCount {
	if board.optMaterialCount != nil {
		return board.optMaterialCount
	}

	materialCount := MaterialCount{}
	for r := uint8(0); r < 8; r++ {
		for c := uint8(0); c < 8; c++ {
			piece := board.Pieces[r][c]
			if piece == WHITE_PAWN {
				materialCount.WhitePawnCount++
			} else if piece == WHITE_KNIGHT {
				materialCount.WhiteKnightCount++
			} else if piece == WHITE_BISHOP {
				square := Square{Rank: r + 1, File: c + 1}
				if square.IsLightSquare() {
					materialCount.WhiteLightBishopCount++
				} else {
					materialCount.WhiteDarkBishopCount++
				}
			} else if piece == WHITE_ROOK {
				materialCount.WhiteRookCount++
			} else if piece == WHITE_QUEEN {
				materialCount.WhiteQueenCount++
			} else if piece == WHITE_KING {
				board.optWhiteKingSquare = &Square{Rank: r + 1, File: c + 1}
			} else if piece == BLACK_PAWN {
				materialCount.BlackPawnCount++
			} else if piece == BLACK_KNIGHT {
				materialCount.BlackKnightCount++
			} else if piece == BLACK_BISHOP {
				square := Square{Rank: r + 1, File: c + 1}
				if square.IsLightSquare() {
					materialCount.BlackLightBishopCount++
				} else {
					materialCount.BlackDarkBishopCount++
				}
			} else if piece == BLACK_ROOK {
				materialCount.BlackRookCount++
			} else if piece == BLACK_QUEEN {
				materialCount.BlackQueenCount++
			} else if piece == BLACK_KING && includeKings {
				board.optBlackKingSquare = &Square{Rank: r + 1, File: c + 1}
			}
		}
	}
	board.optMaterialCount = &materialCount
	return &materialCount
}

func (board *Board) ComputeKingSquares() {
	for r := uint8(0); r < 8; r++ {
		for c := uint8(0); c < 8; c++ {
			piece := board.Pieces[r][c]
			if piece == WHITE_KING {
				square := Square{Rank: r + 1, File: c + 1}
				board.optWhiteKingSquare = &square
			}
			if piece == BLACK_KING {
				square := Square{Rank: r + 1, File: c + 1}
				board.optBlackKingSquare = &square
			}
		}
	}
}

func (board *Board) GetKingSquare(isWhiteKing bool) *Square {
	if board.optWhiteKingSquare == nil || board.optBlackKingSquare == nil {
		board.ComputeKingSquares()
	}
	if isWhiteKing {
		return board.optWhiteKingSquare
	} else {
		return board.optBlackKingSquare
	}
}

func (board *Board) ToFEN() string {
	var fenSegsBuilder strings.Builder

	pieceRuneByPiece := []rune{'x', 'P', 'N', 'B', 'R', 'Q', 'K', 'p', 'n', 'b', 'r', 'q', 'k'}
	fenPiecesSeg := make([]rune, 0, 64)
	for r := 7; r >= 0; r-- {
		consecEmptyCount := 0
		for c := 0; c < 8; c++ {
			piece := board.Pieces[r][c]
			if piece == EMPTY {
				consecEmptyCount++
				continue
			}
			if consecEmptyCount > 0 {
				fenPiecesSeg = append(fenPiecesSeg, rune(consecEmptyCount+48))
				consecEmptyCount = 0
			}
			fenPiecesSeg = append(fenPiecesSeg, pieceRuneByPiece[piece])
		}
		if consecEmptyCount > 0 {
			fenPiecesSeg = append(fenPiecesSeg, rune(consecEmptyCount+48))
			consecEmptyCount = 0
		}
		if r > 0 {
			fenPiecesSeg = append(fenPiecesSeg, '/')
		}
	}
	fenSegsBuilder.WriteString(string(fenPiecesSeg))
	fenSegsBuilder.WriteRune(' ')

	var turnRune rune
	if board.IsWhiteTurn {
		turnRune = 'w'
	} else {
		turnRune = 'b'
	}
	fenSegsBuilder.WriteRune(turnRune)
	fenSegsBuilder.WriteRune(' ')

	fenCastleSeg := make([]rune, 0, 4)
	if board.CanWhiteCastleKingside {
		fenCastleSeg = append(fenCastleSeg, 'K')
	}
	if board.CanWhiteCastleQueenside {
		fenCastleSeg = append(fenCastleSeg, 'Q')
	}
	if board.CanBlackCastleKingside {
		fenCastleSeg = append(fenCastleSeg, 'k')
	}
	if board.CanBlackCastleQueenside {
		fenCastleSeg = append(fenCastleSeg, 'q')
	}
	if len(fenCastleSeg) == 0 {
		fenCastleSeg = append(fenCastleSeg, '-')
	}
	fenSegsBuilder.WriteString(string(fenCastleSeg))
	fenSegsBuilder.WriteRune(' ')

	if board.OptEnPassantSquare != nil {
		fenSegsBuilder.WriteString(board.OptEnPassantSquare.ToAlgebraicCoords())
	} else {
		fenSegsBuilder.WriteRune('-')
	}
	fenSegsBuilder.WriteRune(' ')

	fenSegsBuilder.WriteString(strconv.Itoa(int(board.HalfMoveClockCount)))
	fenSegsBuilder.WriteRune(' ')

	fenSegsBuilder.WriteString(strconv.Itoa(int(board.FullMoveCount)))

	return fenSegsBuilder.String()
}
