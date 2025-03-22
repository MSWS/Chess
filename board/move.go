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

func (move Move) IsCastle() bool {
	if move.piece.GetType() != King {
		return false
	}

	_, fromCol := move.from.GetCoords()
	_, toCol := move.to.GetCoords()

	if fromCol > toCol {
		toCol, fromCol = fromCol, toCol
	}

	return toCol-fromCol > 1
}

func (move Move) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("mv{%v%v", move.from, move.to))

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

func (game Board) GetImmediateMoves() []Move {
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

func (game Board) GetMoves() []Move {
	result := []Move{}

	board := game.Board
	for row := 0; row < len(board); row++ {
		for col := 0; col < len(board[row]); col++ {
			piece := board[row][col]

			if piece == 0 || piece.GetColor() != game.Active {
				continue
			}

			psuedo := game.getMovesFor(CreateCoordInt(row, col))
			legalMoves := []Move{}

			for _, psuedoMove := range psuedo {
				if psuedoMove.capture.GetType() == King {
					continue
				}
				game.MakeMove(psuedoMove)

				enemyMoves := game.GetImmediateMoves()

				legal := true
				for _, enemyMove := range enemyMoves {
					if enemyMove.IsCastle() {
						continue
					}
					if psuedoMove.IsCastle() && enemyMove.capture.GetType() == Rook {
						targetRow, targetCol := psuedoMove.to.GetCoords()
						enemyRow, enemyCol := enemyMove.to.GetCoords()
						if targetRow != enemyRow {
							continue
						}
						if targetCol == 0 && enemyCol == 2 {
							legal = false
							break
						}
						if targetCol == 7 && enemyCol == 5 {
							legal = false
							break
						}
					}

					if enemyMove.capture.GetType() == King {
						legal = false
						break
					}
				}

				if legal {
					legalMoves = append(legalMoves, psuedoMove)
				}

				game.UndoMove()
			}

			result = append(result, legalMoves...)
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
	direction := 1
	if piece.GetColor() == Black {
		direction = -1
	}

	// Basic pushing
	row, col := coord.GetCoords()
	if piece.GetColor() == White {
		if row == 1 && game.Board[row+1][col] == 0 {
			moves = append(moves, game.CreateMove(coord, CreateCoordByte(row+2, col)))
		}
	} else {
		if row == 6 && game.Board[row-1][col] == 0 {
			moves = append(moves, game.CreateMove(coord, CreateCoordByte(row-2, col)))
		}
	}
	moves = append(moves, game.CreateMove(coord, CreateCoordInt(int(row)+direction, int(col))))
	moves = filter(moves, func(m Move) bool {
		return m.capture == 0
	})

	// Capturing
	for _, dx := range []int{-1, 1} {
		if col+byte(dx) > 7 {
			continue
		}
		capture := game.CreateMove(coord, CreateCoordInt(int(row)+direction, int(col)+dx))
		if capture.capture != 0 && capture.piece.GetColor() != capture.capture.GetColor() {
			moves = append(moves, capture)
		}
	}

	// Promoting
	for _, move := range moves {
		row, _ := move.to.GetCoords()

		if row == 0 || row == 7 {
			// A 0 promotionTo defaults to Queen for simplicity
			for _, piece := range []Piece{Knight, Bishop, Rook} {
				move.promotionTo = piece | move.piece.GetColor()
				moves = append(moves, move)
			}
		}
	}

	// En Passant

	if game.EnPassant != nil {
		enRow, enCol := (*game.EnPassant).GetCoords()

		if int(enRow) == int(row)+direction {
			diff := int(enCol) - int(col)
			if diff == -1 || diff == 1 {
				moves = append(moves, game.CreateMove(coord, *game.EnPassant))
			}
		}
	}

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
	return append(moves, game.getCastleMoves(coord)...)
}

func (game Board) getCastleMoves(coord Coordinate) []Move {
	moves := []Move{}
	castling := game.WhiteCastling
	castleRow := 0

	if game.Active == Black {
		castling = game.BlackCastling
		castleRow = 7
	}

	if castling.CanKingSide {
		if game.Board[castleRow][5] == 0 && game.Board[castleRow][6] == 0 {
			moves = append(moves, game.CreateMove(coord, CreateCoordInt(castleRow, 7)))
		}
	}

	if castling.CanQueenSide {
		for col := 1; col <= 3; col++ {
			if game.Board[castleRow][col] != 0 {
				return moves
			}
		}

		moves = append(moves, game.CreateMove(coord, CreateCoordInt(castleRow, 0)))
	}
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
		return m.capture == 0 || m.piece.GetColor() == m.capture.GetColor()
	})
}

func (board Board) filterAllies(moves []Move) []Move {
	return filter(moves, func(m Move) bool {
		return m.capture == 0 || m.piece.GetColor() != m.capture.GetColor()
	})
}
