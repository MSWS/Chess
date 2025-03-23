package board

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Single byte to store 2D coordinates,
// as the max value we need to store in a given
// dimension is 8.
// We store the X coordinate in the upper 4 bits,
// and the Y coordinate in the lower 4 bits
type Coordinate byte

func (coord Coordinate) String() string {
	return coord.GetAlgebra()
}

func (coord Coordinate) GetCoords() (byte, byte) {
	return byte(coord) >> 4, byte(coord) & 0b1111
}

func (coord Coordinate) Add(row int, col int) Coordinate {
	thisRow, thisCol := coord.GetCoords()
	return CreateCoordInt(int(thisRow)+row, int(thisCol)+col)
}

func CreateCoordInt(row int, col int) Coordinate {
	return CreateCoordByte(byte(row), byte(col))
}

func CreateCoordByte(row byte, col byte) Coordinate {
	if row > 7 || col > 7 {
		panic(fmt.Errorf("invalid row, col (%v, %v)", row, col))
	}
	return Coordinate(row<<4 + col)
}

func CreateCoordAlgebra(alg string) Coordinate {
	col := alg[0] - 'a'
	row := alg[1] - '1'

	return CreateCoordByte(row, col)
}

func (coord Coordinate) GetAlgebra() string {
	row, col := coord.GetCoords()
	return fmt.Sprintf("%c%c", col+'a', row+'1')
}

type Castling struct {
	QueenSide bool
	KingSide  bool
}

func (board Game) Get(coord Coordinate) Piece {
	row, col := coord.GetCoords()
	return board.Board[row][col]
}

func (board Game) GetStr(coord string) Piece {
	return board.Get(CreateCoordAlgebra(coord))
}

func (board Game) Set(coord Coordinate, piece Piece) {
	row, col := coord.GetCoords()
	board.Board[row][col] = piece
}

func (board Game) Move(from Coordinate, to Coordinate) {
	fromRow, fromCol := from.GetCoords()
	toRow, toCol := to.GetCoords()
	board.Board[toRow][toCol] = board.Board[fromRow][fromCol]
	board.Board[fromRow][fromCol] = 0
}

func (board *Game) MakeMove(move Move) {
	board.WhiteCastleHistory = append(board.WhiteCastleHistory, board.WhiteCastling)
	board.BlackCastleHistory = append(board.BlackCastleHistory, board.BlackCastling)
	board.PreviousEnpassant = board.EnPassant
	captured := board.Get(move.To)
	move.Capture = captured

	board.Move(move.From, move.To)
	_, fromCol := move.From.GetCoords()

	if move.Piece.GetType() == Pawn {
		toRow, _ := move.To.GetCoords()
		if toRow == 0 || toRow == 7 { // Promotions
			if move.promotionTo == 0 {
				move.promotionTo = Queen | move.Piece.GetColor()
			}
			board.Set(move.To, move.promotionTo|move.Piece.GetColor())
		}

	}
	board.applyEnPassant(&move)

	castlability := &board.WhiteCastling
	if move.Piece.GetColor() == Black {
		castlability = &board.BlackCastling
	}

	if move.Piece.GetType() == Rook {
		if fromCol == 0 {
			castlability.QueenSide = false
		} else if fromCol == 7 {
			castlability.KingSide = false
		}
	}

	if move.Piece.GetType() == King {
		castlability.KingSide = false
		castlability.QueenSide = false

		if move.IsCastle() {
			board.applyCastle(move)
		}
	}

	if move.Capture.GetType() == Rook {
		// If capturing a rook that would've otherwise allowed castling
		// make sure we update the enemy's castlablity

		enemyCastling := &board.BlackCastling
		if move.Capture.GetColor() == White {
			enemyCastling = &board.WhiteCastling
		}

		toRow, toCol := move.To.GetCoords()

		enemyHomeRow := 7
		if move.Capture.GetColor() == White {
			enemyHomeRow = 0
		}

		if toRow == byte(enemyHomeRow) {
			if toCol == 0 {
				enemyCastling.QueenSide = false
			} else if toCol == 7 {
				enemyCastling.KingSide = false
			}
		}
	}

	board.Active = (^board.Active).GetColor()
	board.Moves = append(board.Moves, move)
	// return move
}

func (board *Game) UndoMove() {
	move := board.Moves[len(board.Moves)-1]
	board.Set(move.To, move.Capture)
	board.Set(move.From, move.Piece)

	castling := &board.WhiteCastling
	if move.Piece.GetColor() == Black {
		castling = &board.BlackCastling
	}

	toRow, toCol := move.To.GetCoords()
	fromRow, _ := move.From.GetCoords()

	if move.IsCastle() {
		if toCol == 0 {
			castling.QueenSide = true
			board.Set(CreateCoordInt(int(toRow), 1), 0)
			board.Set(CreateCoordInt(int(toRow), 2), 0)
			board.Set(CreateCoordInt(int(toRow), 3), 0)
		} else {
			castling.KingSide = true
			board.Set(CreateCoordInt(int(toRow), 5), 0)
			board.Set(CreateCoordInt(int(toRow), 6), 0)
		}
	}

	if move.isEnPassant {
		board.Set(move.To, 0)
		captured := CreateCoordByte(fromRow, toCol)
		board.Set(captured, move.Capture)
		board.EnPassant = &move.To
	} else {
		board.EnPassant = board.PreviousEnpassant
	}

	board.Active = (^board.Active).GetColor()
	board.Moves = board.Moves[0 : len(board.Moves)-1]
	board.WhiteCastling = board.WhiteCastleHistory[len(board.WhiteCastleHistory)-1]
	board.BlackCastling = board.BlackCastleHistory[len(board.BlackCastleHistory)-1]

	board.WhiteCastleHistory = board.WhiteCastleHistory[0 : len(board.WhiteCastleHistory)-1]
	board.BlackCastleHistory = board.BlackCastleHistory[0 : len(board.BlackCastleHistory)-1]
}

func (board *Game) applyEnPassant(move *Move) {
	toRow, toCol := move.To.GetCoords()
	fromRow, fromCol := move.From.GetCoords()
	if move.Piece.GetType() != Pawn {
		board.EnPassant = nil
		return
	}

	enemyPiece := Pawn | ((^move.Piece).GetColor())
	if board.EnPassant != nil && move.To == *board.EnPassant {
		// En passant!
		captured := CreateCoordByte(fromRow, toCol)
		if board.Get(captured) != enemyPiece {
			panic(fmt.Sprintf("En Passanted non-enemy piece on %v, got %v, expected %v", captured, board.Get(captured), enemyPiece))
		}
		board.Set(captured, 0)
		move.Capture = enemyPiece
		move.isEnPassant = true
	}
	board.EnPassant = nil

	if (toRow == 3 || toRow == 4) && (fromRow == 1 || fromRow == 6) {
		rowDir := (int(toRow) - int(fromRow)) / 2
		if fromCol > 0 && board.Get(move.To.Add(0, -1)) == enemyPiece {
			enpassant := move.From.Add(rowDir, 0)
			board.EnPassant = &enpassant
		} else if fromCol < 7 && board.Get(move.To.Add(0, 1)) == enemyPiece {
			enpassant := move.From.Add(rowDir, 0)
			board.EnPassant = &enpassant
		}
	}
}

func (board Game) applyCastle(move Move) {
	kingSquare := 2
	rookSquare := 3
	castleRow, castleCol := move.To.GetCoords()

	if castleCol == 7 {
		kingSquare = 6
		rookSquare = 5
	}

	board.Set(move.To, 0)
	board.Set(CreateCoordByte(castleRow, byte(kingSquare)), move.Piece)
	board.Set(CreateCoordByte(castleRow, byte(rookSquare)), move.Capture)
}

func (board *Game) MakeMoveStr(str string) {
	switch str {
	case "e4":
		board.MakeMove(board.CreateMoveStr("e2", "e4"))
	case "e5":
		board.MakeMove(board.CreateMoveStr("e7", "e5"))
	default:
		panic(fmt.Errorf("unsupported make move call with %v", str))
	}
}

// A game that represents a given game.
// Represents all data that is stored in a FEN string.
// i.e. given a Game, you can find the FEN, and vice-versa
type Game struct {
	// The game's pieces, in [row][col]
	// with the first row, first column being the bottom
	// left of the board from white's perspective (i.e. a1)
	Board *[8][8]Piece

	// Active player's turn
	Active Piece

	WhiteCastling Castling
	BlackCastling Castling

	// The square that can be captured en passant
	EnPassant *Coordinate

	// Moves both players have made since last capture or pawn move
	HalfMoves int

	// Number of full moves (incremented after black moves)
	FullMoves int

	Moves              []Move
	WhiteCastleHistory []Castling
	BlackCastleHistory []Castling
	PreviousEnpassant  *Coordinate
}

func (board Game) Equal(other Game) bool {
	return board.ToFEN() == other.ToFEN()
}

func (board Game) String() string {
	return fmt.Sprintf("Board{%s}", board.ToFEN())
}

func (board Game) PrettyPrint() string {
	var builder strings.Builder

	for row := 0; row < len(board.Board); row++ {
		for col := 0; col < len(board.Board[row]); col++ {
			piece := board.Board[row][col]
			if piece == 0 {
				builder.WriteRune('0')
			} else {
				builder.WriteRune(rune(board.Board[row][col].GetRune()))
			}
		}
		builder.WriteRune('\n')
	}

	return builder.String()
}

func FromFEN(str string) (*Game, error) {
	result := Game{}

	records := strings.Split(strings.TrimSpace(str), " ")

	if len(records) != 6 {
		return nil, errors.New("malformed FEN string, expected 6 records")
	}

	board, err := GenerateBoard(records[0])

	if err != nil {
		return &result, err
	}

	result.Board = board

	switch records[1][0] {
	case 'w':
		result.Active = White
	case 'b':
		result.Active = Black
	default:
		return &result, fmt.Errorf("invalid active color %c", records[1][0])
	}

	result.WhiteCastling = formatCastling(records[2], White)
	result.BlackCastling = formatCastling(records[2], Black)

	if records[3][0] != '-' {
		coords := CreateCoordAlgebra(records[3])
		result.EnPassant = &coords
	}

	halfMoves, err := strconv.Atoi(records[4])

	if err != nil {
		return &result, err
	}

	result.HalfMoves = halfMoves

	fullMoves, err := strconv.Atoi(records[5])

	if err != nil {
		return &result, err
	}

	result.FullMoves = fullMoves

	return &result, nil
}

func (board Game) ToFEN() string {
	var result strings.Builder
	result.WriteString(GeneratePieceString(*board.Board))
	result.WriteRune(' ')
	if board.Active == White {
		result.WriteRune('w')
	} else {
		result.WriteRune('b')
	}

	result.WriteRune(' ')
	oldLen := result.Len()

	if board.WhiteCastling.KingSide {
		result.WriteRune('K')
	}

	if board.WhiteCastling.QueenSide {
		result.WriteRune('Q')
	}

	if board.BlackCastling.KingSide {
		result.WriteRune('k')
	}

	if board.BlackCastling.QueenSide {
		result.WriteRune('q')
	}

	if oldLen == result.Len() {
		result.WriteRune('-')
	}

	result.WriteRune(' ')

	if board.EnPassant == nil {
		result.WriteRune('-')
	} else {
		result.WriteString(board.EnPassant.GetAlgebra())
	}

	result.WriteRune(' ')
	result.WriteString(fmt.Sprint(board.HalfMoves))

	result.WriteRune(' ')
	result.WriteString(fmt.Sprint(board.FullMoves))

	return result.String()
}

func (game Game) Perft(depth int) int {
	if depth == 0 {
		return 1
	}

	var nodes int

	moves := game.GetMoves()

	if depth == 1 {
		return len(moves)
	}

	for _, move := range moves {
		game.MakeMove(move)
		nodes += game.Perft(depth - 1)
		game.UndoMove()
	}

	return nodes
}

func formatCastling(str string, color Piece) Castling {
	if color == White {
		return Castling{
			QueenSide: strings.Contains(str, "Q"),
			KingSide:  strings.Contains(str, "K"),
		}
	} else if color == Black {
		return Castling{
			QueenSide: strings.Contains(str, "q"),
			KingSide:  strings.Contains(str, "k"),
		}
	}

	return Castling{}
}
