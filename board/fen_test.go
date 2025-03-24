package board

import "testing"

func TestGenerateBoard(t *testing.T) {
	tests := map[string]struct {
		input  string
		result [8][8]Piece
	}{
		"Starting": {
			input:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
			result: *getStartBoard(),
		},
		"Empty Board": {
			input:  "8/8/8/8/8/8/8/8",
			result: [8][8]Piece{},
		},
		"Test Board": {
			input:  getTestPieceString(),
			result: *getTestBoard(),
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
			result: [8]Piece{Black | Pawn, Black | Pawn, Black | Pawn, Black | Pawn, Black | Pawn, Black | Pawn, Black | Pawn, Black | Pawn},
		},
		"Starting Row": {
			input:  "RNBQKBNR",
			result: getStartRow(),
		},
		"Rooks Only": {
			input:  "r6R",
			result: [8]Piece{Rook | Black, 0, 0, 0, 0, 0, 0, White | Rook},
		},
		"Alternating Pawns": {
			input:  "pP1Pp2P",
			result: [8]Piece{Pawn | Black, White | Pawn, 0, White | Pawn, Pawn | Black, 0, 0, White | Pawn},
		},
		"Many Kings": {
			input:  "KKKKkkkk",
			result: [8]Piece{White | King, White | King, White | King, White | King, Black | King, Black | King, King | Black, King | Black},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := generateBoardRow(test.input)

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

func getStartRow() [8]Piece {
	return [8]Piece{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
}

func getStartBoard() *[8][8]Piece {
	pawnRow := [8]Piece{Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn}

	board := [8][8]Piece{
		getStartRow(), // White starting row
		pawnRow,       // White pawn row
		{},
		{},
		{},
		{},
		pawnRow,       // Black pawn row
		getStartRow(), // Black starting row
	}

	markRowColor(&board[0], White)
	markRowColor(&board[1], White)
	markRowColor(&board[6], Black)
	markRowColor(&board[7], Black)

	return &board
}

func getTestBoard() *[8][8]Piece {
	board := [8][8]Piece{
		{0, 0, White | Pawn, 0, 0, 0, White | Knight},
		{0, Pawn | Black, 0, 0, 0, Bishop | Black},
		{White | Bishop, 0, 0, 0, Black | Rook},
		{0, 0, 0, White | Queen, 0, 0, 0, Black | Queen},
		{Pawn | Black, Pawn | Black, Pawn | Black},
		{Pawn | Black, Pawn | Black},
		{Pawn | Black},
		{0, 0, 0, 0, 0, 0, 0, King | Black},
	}

	return &board
}

func getTestPieceString() string {
	return "7k/p7/pp6/ppp5/3Q3q/B3r3/1p3b2/2P3N1"
}

func TestGeneratePieceString(t *testing.T) {
	tests := map[string]struct {
		input    [8][8]Piece
		expected string
	}{
		"Starting": {
			input:    *getStartBoard(),
			expected: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
		},
		"Empty Board": {
			input:    [8][8]Piece{},
			expected: "8/8/8/8/8/8/8/8",
		},
		"Test Board": {
			input:    *getTestBoard(),
			expected: getTestPieceString(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := GeneratePieceString(test.input)

			if test.expected != result {
				t.Errorf("got invalid board, expected %s, got %s",
					test.expected, result)
			}
		})
	}
}
