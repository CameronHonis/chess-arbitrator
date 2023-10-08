package chess

import (
	"errors"
	"fmt"
	"strings"
)

type Square struct {
	Rank uint8
	File uint8
}

func (s *Square) IsValidBoardSquare() bool {
	return s.Rank > 0 && s.Rank < 9 && s.File > 0 && s.File < 9
}

func (s *Square) ToAlgebraicCoords() string {
	return fmt.Sprintf("%c%d", rune(s.File+96), s.Rank)
}

func (s *Square) IsLightSquare() bool {
	if s.Rank%2 == 0 {
		return s.File%2 == 1
	} else {
		return s.File%2 == 0
	}
}

func (s *Square) IsDarkSquare() bool {
	return !s.IsLightSquare()
}

func SquareFromAlgebraicCoords(algCoords string) (*Square, error) {
	runeCoords := []rune(strings.ToLower(algCoords))
	if len(runeCoords) != 2 {
		return nil, errors.New(fmt.Sprintf("invalid algebraicCoords %s: expected char length (2), got (%d)", algCoords, len(runeCoords)))
	}
	file := uint8(runeCoords[0]) - 96
	rank := uint8(runeCoords[1]) - 48
	square := Square{Rank: rank, File: file}
	if !square.IsValidBoardSquare() {
		return nil, errors.New(fmt.Sprintf("invalid algebraicCoords %s: coords outside of board", algCoords))
	}
	return &square, nil
}
