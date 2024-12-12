package board

import (
	"errors"
	"fmt"
	"strings"
)

func GenerateBoard(str string) (*[8][8]Piece, error) {
	board := &[8][8]Piece{}
	ranks := strings.Split(str, "/")

	if len(ranks) != 8 {
		return board, fmt.Errorf("invalid number of rows (%d)", len(ranks))
	}

	for rank, row := range ranks {
		result, err := GenerateRow(row)

		if err != nil {
			return board, err
		}

		board[7-rank] = *result
	}

	return board, nil
}

func GenerateRow(str string) (*[8]Piece, error) {
	row := &[8]Piece{}
	if len(str) == 0 {
		return row, nil
	}

	column := 0

	for _, c := range str {
		if column >= len(row) {
			return row, errors.New("out of bounds column")
		}

		if c >= '1' && c <= '8' {
			column += int(c - '0')
			continue
		}

		piece, err := GetPiece(c)

		if err != nil {
			return row, err
		}

		row[column] = piece
		column++
	}

	return row, nil
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
