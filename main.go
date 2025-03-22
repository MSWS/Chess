package main

import (
	"fmt"

	"github.com/msws/chess/board"
)

func main() {
	board, err := board.FromFEN("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1")
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
