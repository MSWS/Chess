package board

import (
	"strconv"
	"testing"
)

func getStartMovesForType(piece Piece) []Move {
	start := getStartGame()

	if piece.GetColor()&Black == Black {
		start.Active = Black
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
		t.Run("Total Moves", func(t *testing.T) {
			testForBothColors(t, func(color Piece) {
				start := getStartGame()
				start.Active = color

				result := start.GetMoves()

				if len(result) != 20 {
					t.Errorf("invalid number of legal moves, expected %d, got %d\n%v", 20, len(result), result)
				}
			})
		})

		t.Run("Pawn Moves", func(t *testing.T) {
			testForBothColors(t, func(color Piece) {
				pawnMoves := getStartMovesForType(Pawn | color)
				if len(pawnMoves) != 16 {
					t.Errorf("invalid number of pawn moves, expected %d, got %d\n%v", 16, len(pawnMoves), pawnMoves)
				}
			})
		})

		t.Run("Knight Moves", func(t *testing.T) {
			testForBothColors(t, func(color Piece) {
				knightMoves := getStartMovesForType(Knight | color)
				if len(knightMoves) != 4 {
					t.Errorf("invalid number of knight moves, expected %d, got %d\n%v", 4, len(knightMoves), knightMoves)
				}
			})
		})

		for _, piece := range []Piece{Rook, Bishop, Queen, King} {
			t.Run(string(piece.GetRune()), func(t *testing.T) {
				testForBothColors(t, func(color Piece) {
					moves := getStartMovesForType(piece | color)
					if len(moves) != 0 {
						t.Errorf("invalid number of knight moves, expected %d, got %d\n%v", 0, len(moves), moves)
					}
				})
			})
		}
	})

	t.Run("e4 e5", func(t *testing.T) {
		start := getStartGame()
		start.MakeMoveStr("e4")
		start.MakeMoveStr("e5")

		t.Cleanup(func() {
			start := getStartGame()
			start.MakeMoveStr("e4")
			start.MakeMoveStr("e5")
		})
		t.Run("Total Moves", func(t *testing.T) {
			total := start.GetMoves()

			if len(total) != 29 {
				t.Errorf("invalid number of legal moves, expected %d, got %d\n%v", 29, len(total), total)
			}
		})
	})

	t.Run("Nf3 g5", func(t *testing.T) {
		start, err := FromFEN("rnbqkbnr/pppp1ppp/8/4p3/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 0 2")

		if err != nil {
			t.Error(err)
		}

		total := start.GetMoves()

		if len(total) != 22 {
			t.Errorf("invalid number of legal moves, expected %d, got %d\n%v", 22, len(total), total)
		}
	})
}

func TestMovePly(t *testing.T) {
	knownPerts := []int{1, 20, 400, 8092, 197281, 4865609 /*, 119060324, 3195901860*/}

	start := getStartGame()
	for ply := 1; ply < len(knownPerts); ply++ {
		t.Run(strconv.Itoa(ply), func(t *testing.T) {
			calculated := start.perft(ply)
			if calculated != knownPerts[ply] {
				t.Errorf("expected %v moves, got %v", knownPerts[ply], calculated)
			}
		})
	}
}

func (game Board) perft(depth int) int {
	var nodes int

	moves := game.GetMoves()

	if depth == 1 {
		return len(moves)
	}

	for _, move := range moves {
		game.MakeMove(move)
		nodes += game.perft(depth - 1)
		game.UndoMove(move)
	}

	return nodes
}

func TestCreateMove(t *testing.T) {
	t.Run("Basic Pawn Push", func(t *testing.T) {
		start := getStartGame()
		move := start.CreateMoveStr("e2", "e4")

		if move.from.GetAlgebra() != "e2" {
			t.Errorf("expected origin square to be %v, got %v", "e2", move.from.GetAlgebra())
		}

		if move.to.GetAlgebra() != "e4" {
			t.Errorf("expected target square to be %v, got %v", "e4", move.to.GetAlgebra())
		}

		if move.piece != White|Pawn {
			t.Errorf("expected origin piece to be %v, got %v", Piece(White|Pawn).GetRune(), move.piece.GetRune())
		}

		if move.capture != 0 {
			t.Errorf("expected captured piece to be %v, got %v", 0, move.capture)
		}
	})

	t.Run("Rook x Rook", func(t *testing.T) {
		baseRooks := [8][8]Piece{
			{White | Rook},
			{},
			{},
			{},
			{},
			{},
			{},
			{Black | Rook},
		}
		board := Board{
			Board: &baseRooks,
		}

		move := board.CreateMoveStr("a1", "a8")

		if move.piece != White|Rook {
			t.Errorf("expected origin piece to be %v, got %v", Piece(White|Rook).GetRune(), move.piece.GetRune())
		}
		if move.capture != Black|Rook {
			t.Errorf("expected captured piece to be %v, got %v", Piece(Black|Rook).GetRune(), move.capture.GetRune())
		}
	})
}

func testForBothColors(t *testing.T, test func(color Piece)) {
	t.Run("White", func(t *testing.T) {
		test(White)
	})
	t.Run("Black", func(t *testing.T) {
		test(Black)
	})
}
