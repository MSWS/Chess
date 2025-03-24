package board

import "testing"

func TestGetPiece(t *testing.T) {
	for _, test := range getTestData() {
		t.Run(string(test.rune), func(t *testing.T) {
			result, err := GetPiece(test.rune)

			if err != nil {
				t.Errorf("encountered unexpected error: %v", err)
			}

			if result != test.piece {
				t.Errorf("expected %x, got %x", test.piece, result)
			}
		})
	}
}

func TestGetRune(t *testing.T) {
	for _, test := range getTestData() {
		t.Run(string(test.rune), func(t *testing.T) {
			rune := test.piece.GetRune()
			if rune != test.rune {
				t.Errorf("expected %x to be %c, got %c",
					test.piece, test.rune, rune)
			}
		})
	}
}

func TestGetPieceColor(t *testing.T) {
	for _, color := range []Piece{White, Black} {
		name := "White"
		if color == Black {
			name = "Black"
		}
		t.Run(name, func(t *testing.T) {
			for _, piece := range []Piece{Rook, Knight, Bishop, Queen, King, Pawn} {
				input := piece | color
				t.Run(string(piece.GetRune()), func(t *testing.T) {
					if input.GetColor() != color {
						t.Errorf("%v returned %v instead of %v", piece, input.GetColor(), color)
					}
				})
			}
		})
	}
}

func TestGetPieceType(t *testing.T) {
	for _, color := range []Piece{White, Black} {
		name := "White"
		if color == Black {
			name = "Black"
		}
		t.Run(name, func(t *testing.T) {
			for _, piece := range []Piece{Rook, Knight, Bishop, Queen, King, Pawn} {
				input := piece | color
				t.Run(string(piece.GetRune()), func(t *testing.T) {
					if input.GetType() != piece {
						t.Errorf("%v returned %v instead of %v", piece, input.GetType(), piece)
					}
				})
			}
		})
	}
}

func getTestData() []struct {
	rune  rune
	piece Piece
} {
	return []struct {
		rune  rune
		piece Piece
	}{
		{rune: 'K', piece: White | King},
		{rune: 'Q', piece: White | Queen},
		{rune: 'R', piece: White | Rook},
		{rune: 'B', piece: White | Bishop},
		{rune: 'N', piece: White | Knight},
		{rune: 'P', piece: White | Pawn},
		{rune: 'k', piece: King | Black},
		{rune: 'q', piece: Queen | Black},
		{rune: 'r', piece: Rook | Black},
		{rune: 'b', piece: Bishop | Black},
		{rune: 'n', piece: Knight | Black},
		{rune: 'p', piece: Pawn | Black},
	}
}
