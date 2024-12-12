package board

import "fmt"

type Piece byte

const (
	White Piece = iota
	Black Piece = 1 << (iota - 1)

	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

func (piece Piece) GetLegalMoves(board Board) []Move {
	return []Move{}
}

func GetPiece(c rune) (Piece, error) {
	white := true
	if c > 'Z' {
		white = false
		c = 'A' + (c - 'a')
	}
	var result Piece
	switch c {
	case 'P':
		result = Pawn
	case 'N':
		result = Knight
	case 'B':
		result = Bishop
	case 'R':
		result = Rook
	case 'Q':
		result = Queen
	case 'K':
		result = King
	}

	if result == 0 {
		return 0, fmt.Errorf("unknown piece type: %c (%d)", c, c)
	}

	if !white {
		result |= Black
	}

	return result, nil
}

func (piece Piece) GetRune() rune {
	color := piece & (White | Black)
	p := piece & ^(White | Black)

	var result rune
	switch p {
	case King:
		result = 'K'
	case Queen:
		result = 'Q'
	case Rook:
		result = 'R'
	case Bishop:
		result = 'B'
	case Knight:
		result = 'N'
	case Pawn:
		result = 'P'
	default:
		panic(fmt.Errorf("failed to get rune for piece: %x", piece))
	}

	if color == Black {
		// Upper -> lower
		result += 'a' - 'A'
	}

	return result
}
