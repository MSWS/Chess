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
		if start.Get(move.From).GetType() != piece.GetType() {
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

	t.Run("En Passant", func(t *testing.T) {
		board, err := FromFEN("rnbqkbnr/ppp1p1pp/8/3pPp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 3")
		if err != nil {
			t.Error(err)
		}

		moves := board.GetMoves()

		if len(moves) != 31 {
			t.Errorf("invalid number of legal moves, expected %d, got %d", 31, len(moves))
		}
	})
}

func TestPerfs(t *testing.T) {
	for name, test := range getPerfData() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			start, err := FromFEN(test.FEN)

			if err != nil {
				t.Fatal(err)
			}

			for ply := 1; ply <= len(test.knownPerfs); ply++ {
				success := t.Run(strconv.Itoa(ply), func(t *testing.T) {
					calculated := start.Perft(ply)
					if calculated != test.knownPerfs[ply-1] {
						t.Fatalf("expected %v nodes, got %v", test.knownPerfs[ply-1], calculated)
					}
				})

				if !success {
					break
				}
			}
		})
	}
}

func TestIsCastle(t *testing.T) {
	start := getStartGame()
	t.Run("King-Side", func(t *testing.T) {
		castle := start.CreateMoveStr("e1", "h1")

		if !castle.IsCastle() {
			t.Fail()
		}
	})
	t.Run("Queen-Side", func(t *testing.T) {
		castle := start.CreateMoveStr("e1", "a1")

		if !castle.IsCastle() {
			t.Fail()
		}
	})
}

func TestCreateMove(t *testing.T) {
	t.Run("Basic Pawn Push", func(t *testing.T) {
		start := getStartGame()
		move := start.CreateMoveStr("e2", "e4")

		if move.From.GetAlgebra() != "e2" {
			t.Errorf("expected origin square to be %v, got %v", "e2", move.From.GetAlgebra())
		}

		if move.To.GetAlgebra() != "e4" {
			t.Errorf("expected target square to be %v, got %v", "e4", move.To.GetAlgebra())
		}

		if move.Piece != White|Pawn {
			t.Errorf("expected origin piece to be %v, got %v", Piece(White|Pawn).GetRune(), move.Piece.GetRune())
		}

		if move.Capture != 0 {
			t.Errorf("expected captured piece to be %v, got %v", 0, move.Capture)
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
		board := Game{
			Board: &baseRooks,
		}

		move := board.CreateMoveStr("a1", "a8")

		if move.Piece != White|Rook {
			t.Errorf("expected origin piece to be %v, got %v", Piece(White|Rook).GetRune(), move.Piece.GetRune())
		}
		if move.Capture != Black|Rook {
			t.Errorf("expected captured piece to be %v, got %v", Piece(Black|Rook).GetRune(), move.Capture.GetRune())
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

func getPerfData() map[string]struct {
	FEN        string
	knownPerfs []int
} {
	return map[string]struct {
		FEN        string
		knownPerfs []int
	}{
		"Starting Position": {
			knownPerfs: []int{20, 400, 8902, 197281 /*4865609 /* 119060324, 3195901860*/},
			FEN:        START_POSITION,
		},
		"Nf3 g5": {
			knownPerfs: []int{22, 657, 15616, 478736},
			FEN:        "rnbqkbnr/pppp1ppp/8/4p3/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 0 2",
		},
		"a4 a5": {
			FEN:        "rnbqkbnr/1ppppppp/8/p7/P7/8/1PPPPPPP/RNBQKBNR w KQkq - 0 2",
			knownPerfs: []int{20, 401, 9062, 204508},
		},
		"Kiwipete": {
			knownPerfs: []int{48, 2039, 97862 /* 4085603 /*, 193690690 */},
			FEN:        "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0",
		},
		"Kiwipete - Passant": {
			knownPerfs: []int{44, 2149 /*90978*/ /*4387586*/},
			FEN:        "r3k2r/p1ppqpb1/bn2pnp1/3PN3/Pp2P3/2N2Q1p/1PPBBPPP/R3K2R b KQkq a3 0 1",
		},
		"Kiwipete - Checked": {
			knownPerfs: []int{6, 280, 12919 /*605604*/},
			FEN:        "r3k2r/p1pPqpb1/1n3np1/1b2N3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq - 0 2",
		},
		"Kiwipete - Minimized": {
			knownPerfs: []int{3, 44, 531, 7973 /*115009*/},
			FEN:        "4k2r/3P4/8/4N3/8/8/8/4K3 b k - 0 1",
		},
		"3": {
			knownPerfs: []int{14, 191, 2812, 43238},
			FEN:        "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
		},
		"3 g3+": {
			knownPerfs: []int{4, 54, 1014, 14747},
			FEN:        "8/2p5/3p4/KP5r/1R3p1k/6P1/4P3/8 b - - 0 1",
		},
		"4": {
			knownPerfs: []int{6, 264, 9467},
			FEN:        "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
		},
		"5": {
			knownPerfs: []int{44, 1486, 62379 /* 2103487 */},
			FEN:        "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
		},
		"6": {
			knownPerfs: []int{46, 2079, 89890 /*3894594*/},
			FEN:        "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
		},
		"Gaviota": {
			knownPerfs: []int{14, 97, 1585, 7630, 133028},
			FEN:        "1N6/6k1/8/8/7B/8/8/4K3 w - - 19 103",
		},
	}
}
