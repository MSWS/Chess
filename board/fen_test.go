package board

import "testing"

func TestGenerateBoard(t *testing.T) {
	startRow := [8]Piece{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	pawnRow := [8]Piece{Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn}

	board := [8][8]Piece{
		startRow, // White starting row
		pawnRow,  // White pawn row
		{},
		{},
		{},
		{},
		pawnRow,  // Black pawn row
		startRow, // Black starting row
	}

	markRowColor(&board[0], White)
	markRowColor(&board[1], White)
	markRowColor(&board[6], Black)
	markRowColor(&board[7], Black)

	tests := map[string]struct {
		input  string
		result [8][8]Piece
	}{
		"Starting": {
			input:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
			result: board,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := GenerateBoard(test.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if *result != test.result {
				t.Errorf("boards are not equal, expected %v, got %v",
					test.result, *result)
			}
		})
	}
}

func TestGenerateRow(t *testing.T) {
	tests := map[string]struct {
		input  string
		result [8]Piece
	}{
		"Empty Row": {
			input:  "8",
			result: [8]Piece{},
		},
		"All Pawns": {
			input:  "pppppppp",
			result: [8]Piece{Pawn | Black, Pawn | Black, Pawn | Black, Pawn | Black, Pawn | Black, Pawn | Black, Pawn | Black, Pawn | Black},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := GenerateRow(test.input)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if *result != test.result {
				t.Errorf("rows do not match, expected %v, got %v",
					test.result, *result)
			}
		})
	}
}

func TestGetPiece(t *testing.T) {
	tests := []struct {
		input  rune
		result Piece
	}{
		{input: 'K', result: King | White},
		{input: 'Q', result: Queen | White},
		{input: 'R', result: Rook | White},
		{input: 'B', result: Bishop | White},
		{input: 'N', result: Knight | White},
		{input: 'P', result: Pawn | White},
		{input: 'k', result: King | Black},
		{input: 'q', result: Queen | Black},
		{input: 'r', result: Rook | Black},
		{input: 'b', result: Bishop | Black},
		{input: 'n', result: Knight | Black},
		{input: 'p', result: Pawn | Black},
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			result, err := GetPiece(test.input)

			if err != nil {
				t.Errorf("encountered unexpected error: %v", err)
			}

			if result != test.result {
				t.Errorf("expected %x, got %x", test.result, result)
			}
		})
	}
}
