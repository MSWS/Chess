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
