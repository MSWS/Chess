package main

import (
	"fmt"

	"github.com/msws/chess/board"
)

func main() {
	board, err := board.FromFEN("8/8/4k3/8/8/8/4p3/R3K2R w KQ - 0 1")
	if err != nil {
		panic(err)
	}

	baseMoves := board.GetMoves()

	for _, move := range baseMoves {
		board.MakeMove(move)

		newMoves := board.GetMoves()
		fmt.Printf("%v: %d moves\n", move, len(newMoves))
		fmt.Println(newMoves)

		board.UndoMove()
	}

	fmt.Printf("(%d) %v", len(baseMoves), baseMoves)
}
