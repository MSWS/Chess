package board

import (
	"fmt"
	"regexp"
	"strings"
)

type Move struct {
	From, To    Coordinate
	Piece       Piece
	Capture     Piece
	promotionTo Piece
	isEnPassant bool
}

func (move Move) IsCastle() bool {
	if move.Piece.GetType() != King {
		return false
	}

	_, fromCol := move.From.GetCoords()
	_, toCol := move.To.GetCoords()

	if fromCol > toCol {
		toCol, fromCol = fromCol, toCol
	}

	return toCol-fromCol > 1
}

func (move Move) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("mv{%v%v", move.From, move.To))

	if move.Capture != 0 {
		sb.WriteString(fmt.Sprintf(",X%v", move.Capture))
	}
	if move.promotionTo != 0 {
		sb.WriteString(fmt.Sprintf(",P%v", move.promotionTo))
	}

	sb.WriteString("}")

	return sb.String()
}

func (move Move) GetAlgebra(game *Game) string {
	var result strings.Builder
	_, fromCol := move.From.GetCoords()

	switch move.Piece.GetType() {
	case Pawn:
		result.WriteRune(rune(int(fromCol) + 'a'))
	case King:
		if move.IsCastle() {
			if fromCol == 4 {
				result.WriteString("O-O")
			} else {
				result.WriteString("O-O-O")
			}
		} else {
			result.WriteRune('K')
		}
	case Queen:
		result.WriteRune('Q')
	case Rook:
		result.WriteRune('R')
	case Bishop:
		result.WriteRune('B')
	case Knight:
		result.WriteRune('N')
	}

	if move.Capture != 0 {
		result.WriteString("x")
	}

	result.WriteString(move.To.GetAlgebra())

	// TODO: Checks, Mates
	return result.String()
}

func (board Game) CreateMove(from Coordinate, to Coordinate) Move {
	result := Move{
		From:    from,
		To:      to,
		Piece:   board.Get(from),
		Capture: board.Get(to),
	}

	return result
}

func (board Game) CreateMoveStr(from string, to string) Move {
	return board.CreateMove(
		CreateCoordAlgebra(from),
		CreateCoordAlgebra(to),
	)
}

func (board Game) CreateMoveAlgebra(algebra string) Move {
	algebra = regexp.MustCompile("[x+#]").ReplaceAllString(algebra, "")
	if len(algebra) == 2 {
		// Pawn move
		pawnDir := -1
		if board.Active == Black {
			pawnDir = 1
		}
		ourPawn := board.Active | Pawn

		target := CreateCoordAlgebra(algebra)
		if board.Get(target.Add(pawnDir, 0)) == ourPawn {
			return board.CreateMove(target.Add(pawnDir, 0), target)
		} else if board.Get(target.Add(pawnDir*2, 0)) == ourPawn {
			return board.CreateMove(target.Add(pawnDir*2, 0), target)
		} else {
			panic(fmt.Errorf("invalid move: %v", algebra))
		}
	}

	if algebra[0] == 'O' || algebra[0] == '0' {
		castleRow := 0
		if board.Active == Black {
			castleRow = 7
		}

		fromCol := 4
		toCol := 0

		if len(algebra) == 3 {
			toCol = 7
		}

		return board.CreateMove(
			CreateCoordInt(castleRow, fromCol),
			CreateCoordInt(castleRow, toCol),
		)
	}

	if len(algebra) == 5 {
		// Double disambiguation, we know exactly from -> to
		return board.CreateMoveStr(algebra[1:3], algebra[3:])
	}

	piece := King
	switch strings.ToLower(algebra)[0] {
	case 'n':
		piece = Knight
	case 'b':
		piece = Bishop
	case 'r':
		piece = Rook
	case 'q':
		piece = Queen
	case 'k':
		piece = King
	default:
		piece = Pawn
	}

	return board.createMoveFromTarget(piece, algebra[1:])
}

func (board Game) CreateMoveUCI(uci string) Move {
	from := CreateCoordAlgebra(uci[:2])
	to := CreateCoordAlgebra(uci[2:])

	move := board.CreateMove(from, to)

	if len(uci) == 5 {
		piece, err := GetPiece(rune(uci[4]))
		if err != nil {
			panic(err)
		}

		move.promotionTo = piece
	}

	return move
}

func (board Game) createMoveFromTarget(piece Piece, algebra string) Move {
	if len(algebra) == 3 {
		target := CreateCoordAlgebra(algebra[1:])
		if algebra[0] >= 'a' && algebra[0] <= 'h' {
			return board.createMoveFromDisambiguatedCol(piece, algebra[0]-'a', target)
		}
		if algebra[0] >= '1' && algebra[0] <= '8' {
			return board.createMoveFromDisambiguatedRow(piece, algebra[0]-'1', target)
		}
		panic(fmt.Errorf("invalid move: %v", algebra))
	}

	target := CreateCoordAlgebra(algebra)
	source := board.getSourceCoord(piece, func(coord Coordinate) bool {
		immediateMoves := board.getMovesFor(coord)
		for _, move := range immediateMoves {
			if move.To == target {
				return true
			}
		}

		return false
	})

	return board.CreateMove(source, CreateCoordAlgebra(algebra))
}

func (board Game) createMoveFromDisambiguatedCol(piece Piece, row byte, target Coordinate) Move {
	source := board.getSourceCoord(piece, func(coord Coordinate) bool {
		_, col := coord.GetCoords()
		return col == row
	})

	return board.CreateMove(source, target)
}

func (board Game) createMoveFromDisambiguatedRow(piece Piece, row byte, target Coordinate) Move {
	source := board.getSourceCoord(piece, func(coord Coordinate) bool {
		cRow, _ := coord.GetCoords()
		return row == cRow
	})

	return board.CreateMove(source, target)
}

func (board Game) getSourceCoord(piece Piece, filter func(Coordinate) bool) Coordinate {
	for row := 0; row < len(board.Board); row++ {
		for col := 0; col < len(board.Board[row]); col++ {
			p := board.Board[row][col]
			if p == 0 || p.GetType() != piece || p.GetColor() != board.Active {
				continue
			}
			if !filter(CreateCoordInt(row, col)) {
				continue
			}
			return CreateCoordInt(row, col)
		}
	}

	panic(fmt.Errorf("no matching piece found"))
}

func (game Game) GetImmediateMoves() []Move {
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

func (game Game) GetMoves() []Move {
	result := []Move{}

	board := game.Board
	pieces := 0
	startRow := 0
	rowDir := 1
	if game.Active == Black {
		startRow = 7
		rowDir = -1
	}

	for row := startRow; row < len(board) && row >= 0 && pieces < 16; row += rowDir {
		for col := 0; col < len(board[row]) && pieces < 16; col++ {
			piece := board[row][col]

			if piece == 0 || piece.GetColor() != game.Active {
				continue
			}

			pieces++

			psuedo := game.getMovesFor(CreateCoordInt(row, col))
			legalMoves := []Move{}

			for _, psuedoMove := range psuedo {
				if psuedoMove.Capture.GetType() == King {
					continue
				}

				game.MakeMove(psuedoMove)
				enemyMoves := game.GetImmediateMoves()

				legal := true
				for _, enemyMove := range enemyMoves {
					if enemyMove.Capture.GetType() == King {
						legal = false
						break
					}
					if psuedoMove.IsCastle() {
						enemyRow, enemyCol := enemyMove.To.GetCoords()
						targetRow, targetCol := psuedoMove.To.GetCoords()

						if enemyRow == targetRow && enemyCol == 4 {
							// Cannot castle out of check
							legal = false
							break
						}

						if enemyMove.Capture.GetType() == Rook {
							if targetRow != enemyRow {
								continue
							}
							if targetCol == 0 && enemyCol == 3 {
								legal = false
								break
							}
							if targetCol == 7 && enemyCol == 5 {
								legal = false
								break
							}
						}
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

func (game Game) getMovesFor(coord Coordinate) []Move {
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

func (game Game) getPawnMoves(coord Coordinate) []Move {
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
		return m.Capture == 0
	})

	// Capturing
	for _, dx := range []int{-1, 1} {
		if col+byte(dx) > 7 {
			continue
		}
		capture := game.CreateMove(coord, CreateCoordInt(int(row)+direction, int(col)+dx))
		if capture.Capture != 0 && capture.Piece.GetColor() != capture.Capture.GetColor() {
			moves = append(moves, capture)
		}
	}

	// Promoting
	for _, move := range moves {
		row, _ := move.To.GetCoords()

		if row == 0 || row == 7 {
			// A 0 promotionTo defaults to Queen for simplicity
			for _, piece := range []Piece{Knight, Bishop, Rook} {
				move.promotionTo = piece | move.Piece.GetColor()
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
				move := game.CreateMove(coord, *game.EnPassant)
				move.isEnPassant = true
				moves = append(moves, move)
			}
		}
	}

	return moves
}

func (game Game) getKnightMoves(coord Coordinate) []Move {
	moves := []Move{}

	offsets := [][]int{{-2, 1}, {-1, 2}, {1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}, {-2, -1}}
	sx, sy := coord.GetCoords()

	for _, offset := range offsets {
		tx := sx + byte(offset[0])
		ty := sy + byte(offset[1])

		if tx > 7 || ty > 7 {
			continue
		}

		moves = append(moves, game.CreateMove(coord, coord.Add(offset[0], offset[1])))
	}

	moves = game.filterAllies(moves)
	return moves
}

func (game Game) getBishopMoves(coord Coordinate) []Move {
	return game.getSlidingMovesOf(coord, [][]int{{-1, 1}, {1, 1}, {1, -1}, {-1, -1}})
}

func (game Game) getRookMoves(coord Coordinate) []Move {
	return game.getSlidingMovesOf(coord, [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}})
}

func (game Game) getQueenMoves(coord Coordinate) []Move {
	return append(game.getBishopMoves(coord), game.getRookMoves(coord)...)
}

func (game Game) getKingMoves(coord Coordinate) []Move {
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

func (game Game) getCastleMoves(coord Coordinate) []Move {
	moves := []Move{}
	castling := game.WhiteCastling
	castleRow := 0

	pawnCheckRow := 1

	if game.Active == Black {
		castling = game.BlackCastling
		castleRow = 7
		pawnCheckRow = 6
	}

	// Edge case where once the king has moved, the pawn would
	// no longer be able to capture.
	// But the pawn being there in the first place already marks
	// the king in check, regardless of the pawn's ability to capture.
	enemyPawn := (^game.Active).GetColor() | Pawn
	for _, col := range []int{3, 5} {
		piece := game.Board[pawnCheckRow][col]
		if piece != enemyPawn {
			continue
		}
		return moves
	}

	if castling.KingSide {
		if game.Board[castleRow][5] == 0 && game.Board[castleRow][6] == 0 {
			moves = append(moves, game.CreateMove(coord, CreateCoordInt(castleRow, 7)))
		}
	}

	if castling.QueenSide {
		for col := 1; col <= 3; col++ {
			if game.Board[castleRow][col] != 0 {
				return moves
			}
		}

		moves = append(moves, game.CreateMove(coord, CreateCoordInt(castleRow, 0)))
	}
	return moves
}

func (game Game) getSlidingMovesOf(coord Coordinate, offsets [][]int) []Move {
	moves := []Move{}

	for _, offset := range offsets {
		moves = append(moves, game.getSlidingMoves(coord, offset[0], offset[1])...)
	}

	return moves
}

func (game Game) getSlidingMoves(coord Coordinate, offsetX int, offsetY int) []Move {
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

		if move.Capture == 0 {
			moves = append(moves, move)
			continue
		}

		if move.Piece.GetColor() != move.Capture.GetColor() {
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

func (board Game) filterEnemies(moves []Move) []Move {
	return filter(moves, func(m Move) bool {
		return m.Capture == 0 || m.Piece.GetColor() == m.Capture.GetColor()
	})
}

func (board Game) filterAllies(moves []Move) []Move {
	return filter(moves, func(m Move) bool {
		return m.Capture == 0 || m.Piece.GetColor() != m.Capture.GetColor()
	})
}
