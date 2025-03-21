package board

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMakeMove(t *testing.T) {
	t.Run("e4", func(t *testing.T) {
		start := getStartGame()

		e4 := start.CreateMoveStr("e2", "e4")

		start.MakeMove(e4)

		residual := start.GetStr("e2")

		if residual != 0 {
			t.Errorf("board did not properly move pawn, expected 0, got %x",
				residual)
		}

		residual = start.GetStr("e4")

		if residual != White|Pawn {
			t.Errorf("board did not properly move pawn, expected %x, got %x",
				White|Pawn, residual)
		}
	})

	t.Run("Castling", func(t *testing.T) {
		t.Run("Updates", func(t *testing.T) {
			start, err := FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQkq - 0 1")
			if err != nil {
				t.Fatal(err)
			}
			move := start.CreateMoveStr("e1", "f1")
			start.MakeMove(move)

			castle := start.WhiteCastling

			if castle.CanKingSide || castle.CanQueenSide {
				t.Errorf("board did not mark white as longer able to castle (%v)", castle)
			}
		})
		t.Run("Succeeds", func(t *testing.T) {
			start, err := FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQkq - 0 1")
			if err != nil {
				t.Fatal(err)
			}

			castle := start.CreateMoveStr("e1", "h1")
			start.MakeMove(castle)
			castlability := start.WhiteCastling

			if castlability.CanKingSide || castlability.CanQueenSide {
				t.Errorf("board did not mark white as no longer able to castle (%v)", castlability)
			}

			if start.GetStr("f1") != White|Rook {
				t.Errorf("board did not properly place white rook at f1, got %v", start.GetStr("f1"))
			}

			if start.GetStr("g1") != White|King {
				t.Errorf("board did not properly place white king at f1, got %v", start.GetStr("g1"))
			}

			if start.GetStr("h1") != 0 {
				t.Errorf("board did not properly remove rook at h1, got %v", start.GetStr("h1"))
			}
		})
	})

	t.Run("Promotion", func(t *testing.T) {
		start, err := FromFEN("8/3P4/8/8/8/8/8/8 w - - 0 1")

		if err != nil {
			t.Error(err)
		}

		move := start.CreateMoveStr("d7", "d8")
		move.promotionTo = Queen

		start.MakeMove(move)

		if start.Get(move.to) != White|Queen {
			t.Errorf("board did not properly promote pawn to queen, got %c", start.Get(move.to).GetRune())
		}
	})
}

func TestUndoMove(t *testing.T) {
	t.Run("ChangesTurn", func(t *testing.T) {
		start := getStartGame()
		start.MakeMoveStr("e4")
		if start.Active != Black {
			t.Error("board failed to change active player after moving")
		}

		start.UndoMove()
		if start.Active != White {
			t.Error("board failed to change active player after undoing")
		}
	})

	t.Run("Reverts Pawn", func(t *testing.T) {
		start := getStartGame()
		start.MakeMoveStr("e4")
		start.UndoMove()

		startFen := start.ToFEN()
		if startFen != START_POSITION {
			t.Errorf("reverting to original position did not restore to %v (got %v)", START_POSITION, startFen)
		}
	})

	t.Run("Reverts Castling", func(t *testing.T) {
		start, err := FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQkq - 0 1")
		if err != nil {
			t.Error(err)
		}

		kingMove := start.CreateMoveStr("e1", "f1")
		start.MakeMove(kingMove)
		start.UndoMove()

		castlability := start.WhiteCastling

		if !castlability.CanKingSide || !castlability.CanQueenSide {
			t.Errorf("board did not update white castlability after undoing (%v)", castlability)
		}
	})
}

func TestMakeMoveStr(t *testing.T) {
	start := getStartGame()
	start.MakeMoveStr("e4")
	residual := start.GetStr("e2")

	if residual != 0 {
		t.Errorf("board did not properly move pawn, expected 0, got %x",
			residual)
	}

	residual = start.GetStr("e4")

	if residual != White|Pawn {
		t.Errorf("board did not properly move pawn, expected %x, got %x",
			White|Pawn, residual)
	}
}

func TestGetCoords(t *testing.T) {
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
			expected: CreateCoordByte(1, 6),
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

func TestCreateEquivalence(t *testing.T) {
	tests := []struct {
		input    [2]int
		expected Coordinate
	}{
		{
			input:    [2]int{0, 0},
			expected: 0<<4 + 0,
		},
		{
			input:    [2]int{1, 1},
			expected: CreateCoordInt(1, 1),
		},
		{
			input:    [2]int{6, 1},
			expected: CreateCoordByte(6, 1),
		},
		{
			input:    [2]int{3, 2},
			expected: CreateCoordByte(3, 2),
		},
		{
			input:    [2]int{7, 7},
			expected: 7<<4 + 7,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprint(test.input), func(t *testing.T) {
			intMade := CreateCoordInt(test.input[0], test.input[1])
			byteMade := CreateCoordByte(byte(test.input[0]), byte(test.input[1]))

			if intMade != byteMade {
				t.Errorf("int is not equivalent as byte (%v vs %v)", intMade, byteMade)
			}

			if intMade != test.expected {
				t.Errorf("int is not expected (got %v, expected %v)", intMade, test.expected)
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
			input:    CreateCoordByte(1, 6),
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
		startBoard := getStartGame()

		resultBoard, err := FromFEN(START_POSITION)

		if err != nil {
			t.Errorf("encountered error when parsing from fen: %v", err)
		}

		if !cmp.Equal(startBoard, *resultBoard) {
			t.Errorf("Boards are not equal, expected %v, got %v", startBoard, *resultBoard)
		}
	})

	t.Run("Simple Opening", func(t *testing.T) {
		pos := "rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2"

		coords := Coordinate(5<<4 + 2)

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
			expected: START_POSITION,
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
