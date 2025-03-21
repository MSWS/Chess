package board

import (
	"fmt"
	"strings"
)

type Move struct {
	from, to    Coordinate
	piece       Piece
	capture     Piece
	promotionTo Piece
}

func (move Move) String() string {
	fx, fy := move.from.GetCoords()
	tx, ty := move.to.GetCoords()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("mv{%v,%v->%v,%v", fx, fy, tx, ty))

	if move.capture != 0 {
		sb.WriteString(fmt.Sprintf(",X%v", move.capture))
	}
	if move.promotionTo != 0 {
		sb.WriteString(fmt.Sprintf(",P%v", move.promotionTo))
	}

	sb.WriteString("}")

	return sb.String()
}

func (board Board) CreateMove(from Coordinate, to Coordinate) Move {
	result := Move{
		from:    from,
		to:      to,
		piece:   board.Get(from),
		capture: board.Get(to),
	}

	return result
}

func (board Board) CreateMoveStr(from string, to string) Move {
	return board.CreateMove(
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

			if piece == 0 || piece.GetColor() != game.Active {
				continue
			}

			result = append(result, game.getMovesFor(CreateCoordInt(row, col))...)
		}
	}

	return result
}

func (game Board) getMovesFor(coord Coordinate) []Move {
	piece := game.Get(coord)

	switch piece.GetType() {
	case Pawn:
		return game.getPawnMoves(coord)
	case Knight:
		return game.getKnightMoves(coord)
	case Bishop:
		return game.getBishopMoves(coord)
	case Rook:
		return game.getRookMoves(coord)
	case Queen:
		return game.getQueenMoves(coord)
	case King:
		return game.getKingMoves(coord)
	default:
		panic(fmt.Errorf("unknown piece type: %v", piece))
	}
}

func (game Board) getPawnMoves(coord Coordinate) []Move {
	piece := game.Get(coord)
	moves := []Move{}

	if piece.GetColor() == White {
		row, col := coord.GetCoords()
		if row == 1 {
			moves = append(moves, game.CreateMove(coord, CreateCoordByte(row+2, col)))
		}

		moves = append(moves, game.CreateMove(coord, CreateCoordByte(row+1, col)))
	} else {
		row, col := coord.GetCoords()
		if row == 6 {
			moves = append(moves, game.CreateMove(coord, CreateCoordByte(row-2, col)))
		}

		moves = append(moves, game.CreateMove(coord, CreateCoordByte(row-1, col)))
	}

	moves = filter(moves, func(m Move) bool {
		target := game.Get(m.to)
		return target == 0
	})

	return moves
}

func (game Board) getKnightMoves(coord Coordinate) []Move {
	moves := []Move{}

	offsets := [][]int{{-2, 1}, {-1, 2}, {1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}, {-2, -1}}
	sx, sy := coord.GetCoords()

	for _, offset := range offsets {
		tx := sx + byte(offset[0])
		ty := sy + byte(offset[1])

		if tx > 7 || ty > 7 {
			continue
		}

		toCoord := CreateCoordByte(tx, ty)
		moves = append(moves, game.CreateMove(coord, toCoord))
	}

	moves = game.filterAllies(moves)
	return moves
}

func (game Board) getBishopMoves(coord Coordinate) []Move {
	return game.getSlidingMovesOf(coord, [][]int{{-1, 1}, {1, 1}, {1, -1}, {-1, -1}})
}

func (game Board) getRookMoves(coord Coordinate) []Move {
	return game.getSlidingMovesOf(coord, [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}})
}

func (game Board) getQueenMoves(coord Coordinate) []Move {
	return append(game.getBishopMoves(coord), game.getRookMoves(coord)...)
}

func (game Board) getKingMoves(coord Coordinate) []Move {
	moves := []Move{}

	offsets := [][]int{{-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1}, {-1, 0}}
	sx, sy := coord.GetCoords()

	for _, offset := range offsets {
		tx := sx + byte(offset[0])
		ty := sy + byte(offset[1])

		if tx > 7 || ty > 7 {
			continue
		}

		toCoord := CreateCoordByte(tx, ty)
		moves = append(moves, game.CreateMove(coord, toCoord))
	}

	moves = game.filterAllies(moves)
	return moves
}

func (game Board) getSlidingMovesOf(coord Coordinate, offsets [][]int) []Move {
	moves := []Move{}

	for _, offset := range offsets {
		moves = append(moves, game.getSlidingMoves(coord, offset[0], offset[1])...)
	}

	return moves
}

func (game Board) getSlidingMoves(coord Coordinate, offsetX int, offsetY int) []Move {
	moves := []Move{}

	current := coord

	for {
		cRow, cCol := current.GetCoords()

		cRow += byte(offsetY)
		cCol += byte(offsetX)

		if cRow > 7 || cCol > 7 {
			break
		}

		current = CreateCoordByte(cRow, cCol)

		move := game.CreateMove(coord, current)

		if move.capture == 0 {
			moves = append(moves, move)
			continue
		}

		if move.piece.GetColor() != move.capture.GetColor() {
			moves = append(moves, move)
		}
		break
	}

	return moves
}

func filter[T any](arr []T, predicate func(T) bool) []T {
	ret := []T{}
	for _, t := range arr {
		if !predicate(t) {
			continue
		}

		ret = append(ret, t)
	}

	return ret
}

func (board Board) filterEnemies(moves []Move) []Move {
	return filter(moves, func(m Move) bool {
		from := board.Get(m.from)
		target := board.Get(m.to)
		return target == 0 || target.GetColor() == from.GetColor()
	})
}

func (board Board) filterAllies(moves []Move) []Move {
	return filter(moves, func(m Move) bool {
		from := board.Get(m.from)
		target := board.Get(m.to)
		return target == 0 || target.GetColor() != from.GetColor()
	})
}
