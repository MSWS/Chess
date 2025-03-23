package main

import (
	"fmt"

	"github.com/msws/chess/board"
)

func main() {
	board, err := board.FromFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/P7/1PP1NnPP/RNBQK2R b KQ - 0 8")
	if err != nil {
		panic(err)
	}

	baseMoves := board.GetMoves()

	depth := 5
	total := 0

	for _, move := range baseMoves {
		board.MakeMove(move)

		perft := board.Perft(depth - 1)
		total += perft
		fmt.Printf("%v%v: %d\n", move.From.GetAlgebra(), move.To.GetAlgebra(), perft)

		board.UndoMove()
	}

	fmt.Printf("(%d)", total)
}
