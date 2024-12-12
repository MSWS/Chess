package board

import (
	"errors"
	"fmt"
	"strconv"
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

func CreateCoordAlgebra(alg string) Coordinate {
	x := alg[0] - 'a'
	y := alg[1] - '1'

	return CreateCoordByte(x, y)
}

func (coord Coordinate) GetAlgebra() string {
	x, y := coord.GetCoords()
	return fmt.Sprintf("%c%c", x+'a', y+'1')
}

type Castlability struct {
	CanQueenSide bool
	CanKingSide  bool
}

func (board Board) Get(coord Coordinate) Piece {
	x, y := coord.GetCoords()
	return board.Board[x][y]
}

func (board Board) Move(from Coordinate, to Coordinate) {
	fromX, fromY := from.GetCoords()
	toX, toY := to.GetCoords()
	board.Board[toX][toY] = board.Board[fromX][fromY]
}

// A board that represents a given game.
// Represents all data that is stored in a FEN string.
// i.e. given a Board, you can find the FEN, and vice-versa
type Board struct {
	// The game's pieces, in [row][col]
	// with the first row, first column being the bottom
	// left of the board from white's perspective (i.e. a1)
	Board *[8][8]Piece

	// Active player's turn
	Active Piece

	WhiteCastling Castlability
	BlackCastling Castlability
	EnPassant     *Coordinate
	HalfMoves     int
	FullMoves     int
}

func (board Board) Equal(other Board) bool {
	return board.ToFEN() == other.ToFEN()
}

func (board Board) String() string {
	return fmt.Sprintf("Board{%s}", board.ToFEN())
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

	switch records[1][0] {
	case 'w':
		result.Active = White
	case 'b':
		result.Active = Black
	default:
		return &result, fmt.Errorf("invalid active color %c", records[1][0])
	}

	result.WhiteCastling = getCastlability(records[2], White)
	result.BlackCastling = getCastlability(records[2], Black)

	if records[3][0] != '-' {
		coords := CreateCoordAlgebra(records[3])
		result.EnPassant = &coords
	}

	halfMoves, err := strconv.Atoi(records[4])

	if err != nil {
		return &result, err
	}

	result.HalfMoves = halfMoves

	fullMoves, err := strconv.Atoi(records[5])

	if err != nil {
		return &result, err
	}

	result.FullMoves = fullMoves

	return &result, nil
}

func (board Board) ToFEN() string {
	var result strings.Builder
	result.WriteString(GeneratePieceString(*board.Board))
	result.WriteRune(' ')
	if board.Active == White {
		result.WriteRune('w')
	} else {
		result.WriteRune('b')
	}

	result.WriteRune(' ')
	oldLen := result.Len()

	if board.WhiteCastling.CanKingSide {
		result.WriteRune('K')
	}

	if board.WhiteCastling.CanQueenSide {
		result.WriteRune('Q')
	}

	if board.BlackCastling.CanKingSide {
		result.WriteRune('k')
	}

	if board.BlackCastling.CanQueenSide {
		result.WriteRune('q')
	}

	if oldLen == result.Len() {
		result.WriteRune('-')
	}

	result.WriteRune(' ')

	if board.EnPassant == nil {
		result.WriteRune('-')
	} else {
		result.WriteString(board.EnPassant.GetAlgebra())
	}

	result.WriteRune(' ')
	result.WriteString(fmt.Sprint(board.HalfMoves))

	result.WriteRune(' ')
	result.WriteString(fmt.Sprint(board.FullMoves))

	return result.String()
}

func getCastlability(str string, color Piece) Castlability {
	if color == White {
		return Castlability{
			CanQueenSide: strings.Contains(str, "Q"),
			CanKingSide:  strings.Contains(str, "K"),
		}
	} else if color == Black {
		return Castlability{
			CanQueenSide: strings.Contains(str, "q"),
			CanKingSide:  strings.Contains(str, "k"),
		}
	}

	return Castlability{}
}
