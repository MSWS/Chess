package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/msws/chess/board"
)

func main() {
	board, err := board.FromFEN("r3k2r/p1pPqpb1/1n3np1/1b2N3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq - 0 2")
	if err != nil {
		panic(err)
	}

	baseMoves := board.GetMoves()

	depth := 1
	total := 0

	println(strings.Join(os.Args, ","))

	if len(os.Args) == 2 {
		val, err := strconv.Atoi(os.Args[1])
		if err != nil {
			panic(err)
		}
		depth = val
	}

	// fmt.Printf("Depth: %d\n", depth)

	for _, move := range baseMoves {
		board.MakeMove(move)

		perft := board.Perft(depth - 1)
		total += perft
		fmt.Printf("%v%v: %d\n", move.From.GetAlgebra(), move.To.GetAlgebra(), perft)

		board.UndoMove()
	}

	fmt.Printf("Depth: %d, (%d)", depth, total)
}
