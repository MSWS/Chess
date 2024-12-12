package board

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetCoords(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		input byte
		x, y  byte
	}{
		"Zero": {
			input: 0b0000_0000,
			x:     0,
			y:     0,
		},
		"Middle": {
			input: 3<<4 + 3,
			x:     3,
			y:     3,
		},
		"Last": {
			input: 7<<4 + 7,
			x:     7,
			y:     7,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			coord := Coordinate(test.input)

			x, y := coord.GetCoords()

			if x != test.x {
				t.Errorf("Expected x to be %x, got %x", test.x, x)
			}
			if y != test.y {
				t.Errorf("Expected y to be %x, got %x", test.y, y)
			}
		})
	}
}

func TestCreateCoordInt(t *testing.T) {
	tests := map[string]struct {
		x        int
		y        int
		expected byte
	}{
		"Zero": {
			x:        0,
			y:        0,
			expected: 0b0000_0000,
		},
		"Middle": {
			x:        4,
			y:        4,
			expected: 4<<4 + 4,
		},
		"Last": {
			x:        7,
			y:        7,
			expected: 7<<4 + 7,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			coord := CreateCoordInt(test.x, test.y)

			if byte(coord) != test.expected {
				t.Errorf("Expected (%d, %d) to become %x, got %x",
					test.x, test.y, test.expected, coord)
			}
		})
	}
}

func TestFromFEN(t *testing.T) {
	startPos := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	startRow := [8]Piece{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	pawnRow := [8]Piece{Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn}

	startBoard := &Board{
		Board: &[8][8]Piece{
			startRow, // White starting row
			pawnRow,  // White pawn row
			{},
			{},
			{},
			{},
			pawnRow,  // Black pawn row
			startRow, // Black starting row
		},
	}

	markRowColor(&startBoard.Board[0], White)
	markRowColor(&startBoard.Board[1], White)
	markRowColor(&startBoard.Board[6], Black)
	markRowColor(&startBoard.Board[7], Black)

	resultBoard, err := FromFEN(startPos)

	if err != nil {
		t.Errorf("encountered error when parsing from fen: %v", err)
	}

	if !cmp.Equal(*startBoard, *resultBoard) {
		t.Errorf("Boards are not equal, expected %v, got %v", *startBoard, *resultBoard)
	}
}

func markRowColor(row *[8]Piece, color Piece) {
	for index, piece := range row {
		if piece == 0 {
			continue
		}

		row[index] = row[index] | color
	}
}
