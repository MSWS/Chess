package board

import "testing"

func TestGetMoves(t *testing.T) {
	t.Run("Starting Position", func(t *testing.T) {
		start := getStartGame()

		result := start.GetMoves()

		if len(result) != 20 {
			t.Errorf("invalid number of legal moves, expected %d, got %d", 20, len(result))
		}
	})
}
