package board

import (
	"errors"
	"strings"
)

// Single byte to store 2D coordinates,
// as the max value we need to store in a given
// dimension is 8.
// We store the X coordinate in the upper 4 bits,
// and the Y coordinate in the lower 4 bits
type Coordinate byte

func (coord Coordinate) GetCoords() (byte, byte) {
	return byte(coord) >> 4, byte(coord) & 0b1111
}

func CreateCoordInt(x int, y int) Coordinate {
	return CreateCoordByte(byte(x), byte(y))
}

func CreateCoordByte(x byte, y byte) Coordinate {
	return Coordinate(x<<4 + y)
}

type Castlability struct {
	CanQueenSide bool
	CanKingSide  bool
}

// A board that represents a given game.
// Represents all data that is stored in a FEN string.
// i.e. given a Board, you can find the FEN, and vice-versa
type Board struct {
	// The game's pieces, in [row][col]
	// with the first row, first column being the bottom
	// left of the board from white's perspective (i.e. a1)
	Board         *[8][8]Piece
	Active        Piece
	WhiteCastling Castlability
	BlackCastling Castlability
	EnPassant     Coordinate
	HalfMoves     int
	FullMoves     int
}

func FromFEN(str string) (*Board, error) {
	result := Board{}

	records := strings.Split(str, " ")

	if len(records) != 6 {
		return nil, errors.New("malformed FEN string, expected 6 records")
	}

	board, err := GenerateBoard(records[0])

	if err != nil {
		return &result, err
	}

	result.Board = board
	return &result, nil
}
