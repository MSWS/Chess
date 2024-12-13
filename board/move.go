package board

type Move struct {
	from, to    Coordinate
	capture     Piece
	promotionTo Piece
}

func CreateMove(from Coordinate, to Coordinate) Move {
	result := Move{
		from: from,
		to:   to,
	}

	return result
}

func CreateMoveStr(from string, to string) Move {
	return CreateMove(
		CreateCoordAlgebra(from),
		CreateCoordAlgebra(to),
	)
}

func (game Board) GetMoves() []Move {
	result := []Move{}

	board := game.Board
	for row := 0; row < len(board); row++ {
		for col := 0; col < len(board[row]); col++ {
			piece := board[row][col]

			if piece&game.Active != game.Active {
				continue
			}

			// result = append(result, piece.GetLegalMoves(game)...)
		}
	}

	return result
}

func (game Board) getMovesFor(coord Coordinate) []Move {
	piece := game.Get(coord)

	switch piece.GetType() {
	case Pawn:
		return game.getPawnMoves(coord)
	}

	return []Move{}
}

func (game Board) getPawnMoves(coord Coordinate) []Move {

}
