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
				t.Errorf("expected (%d, %d) to become %x, got %x",
					test.x, test.y, test.expected, coord)
			}
		})
	}
}

func TestCreateCoordAlgebra(t *testing.T) {
	tests := []struct {
		input    string
		expected Coordinate
	}{
		{
			input:    "a1",
			expected: 0<<4 + 0,
		},
		{
			input:    "b2",
			expected: CreateCoordInt(1, 1),
		},
		{
			input:    "g2",
			expected: CreateCoordByte(6, 1),
		},
		{
			input:    "h8",
			expected: 7<<4 + 7,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()

			result := CreateCoordAlgebra(test.input)

			if result != test.expected {
				t.Errorf("expected %s to become %x, got %x",
					test.input, test.expected, result)
			}
		})
	}
}

func TestGetAlgebra(t *testing.T) {
	tests := []struct {
		input    Coordinate
		expected string
	}{
		{
			input:    CreateCoordInt(0, 0),
			expected: "a1",
		},
		{
			input:    CreateCoordInt(1, 1),
			expected: "b2",
		},
		{
			input:    CreateCoordByte(6, 1),
			expected: "g2",
		},
		{
			input:    CreateCoordAlgebra("f5"),
			expected: "f5",
		},
		{
			input:    CreateCoordByte(7, 7),
			expected: "h8",
		},
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			t.Parallel()
			result := test.input.GetAlgebra()

			if result != test.expected {
				t.Errorf("got invalid algebra notation for %x, got %s, expected %s",
					test.input, result, test.expected)
			}
		})
	}
}

func TestFromFEN(t *testing.T) {
	startRow := [8]Piece{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	t.Run("Start Position", func(t *testing.T) {
		startPos := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

		startBoard := getStartGame()

		resultBoard, err := FromFEN(startPos)

		if err != nil {
			t.Errorf("encountered error when parsing from fen: %v", err)
		}

		if !cmp.Equal(startBoard, *resultBoard) {
			t.Errorf("Boards are not equal, expected %v, got %v", startBoard, *resultBoard)
		}
	})

	t.Run("Simple Opening", func(t *testing.T) {
		pos := "rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2"

		coords := Coordinate(2<<4 + 5)

		expectedBoard := &Board{
			Board: &[8][8]Piece{
				startRow,
				{Pawn, Pawn, Pawn, Pawn, 0, Pawn, Pawn, Pawn},
				{},
				{0, 0, 0, 0, Pawn | White},
				{0, 0, Pawn | Black},
				{},
				{Pawn, Pawn, 0, Pawn, Pawn, Pawn, Pawn, Pawn},
				startRow,
			},
			Active:        White,
			WhiteCastling: Castlability{true, true},
			BlackCastling: Castlability{true, true},
			EnPassant:     &coords,
			HalfMoves:     0,
			FullMoves:     2,
		}
		markRowColor(&expectedBoard.Board[0], White)
		markRowColor(&expectedBoard.Board[1], White)
		markRowColor(&expectedBoard.Board[6], Black)
		markRowColor(&expectedBoard.Board[7], Black)

		resultBoard, err := FromFEN(pos)

		if err != nil {
			t.Errorf("encountered error when parsing from fen: %v", err)
		}

		if !cmp.Equal(*expectedBoard, *resultBoard) {
			t.Errorf("boards are not equal, expected %v, got %v", *expectedBoard, *resultBoard)
		}
	})
}

func TestToFEN(t *testing.T) {
	tests := map[string]struct {
		input    Board
		expected string
	}{
		"Start Position": {
			input:    getStartGame(),
			expected: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result := test.input.ToFEN()

			if result != test.expected {
				t.Errorf("FEN strings are not equal, expected %s, got %s",
					test.expected, result)
			}
		})
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

func getStartGame() Board {
	return Board{
		Board:         getStartBoard(),
		Active:        White,
		WhiteCastling: Castlability{true, true},
		BlackCastling: Castlability{true, true},
		EnPassant:     nil,
		HalfMoves:     0,
		FullMoves:     1,
	}
}
