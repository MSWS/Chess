package board

import (
	"fmt"
	"testing"
)

func getStartMovesForType(piece Piece) []Move {
	start := getStartGame()

	if piece.GetColor()&Black == Black {
		fmt.Printf("setting color to black")
	}

	result := []Move{}

	for _, move := range start.GetMoves() {
		if start.Get(move.from).GetType() != piece.GetType() {
			continue
		}

		result = append(result, move)
	}

	return result
}

func TestGetMoves(t *testing.T) {
	t.Run("Starting Position", func(t *testing.T) {
		start := getStartGame()

		result := start.GetMoves()

		if len(result) != 20 {
			t.Errorf("invalid number of legal moves, expected %d, got %d\n%v", 20, len(result), result)
		}
	})

	t.Run("Pawn Moves", func(t *testing.T) {
		pawnMoves := getStartMovesForType(Pawn)
		if len(pawnMoves) != 16 {
			t.Errorf("invalid number of pawn moves, expected %d, got %d\n%v", 16, len(pawnMoves), pawnMoves)
		}
	})

	t.Run("Knight Moves", func(t *testing.T) {
		knightMoves := getStartMovesForType(Knight)
		if len(knightMoves) != 4 {
			t.Errorf("invalid number of knight moves, expected %d, got %d\n%v", 4, len(knightMoves), knightMoves)
		}
	})
}
