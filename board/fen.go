package board

import (
	"errors"
	"fmt"
	"strings"
)

func GeneratePieceString(arr [8][8]Piece) string {
	var sb strings.Builder
	for row := len(arr) - 1; row >= 0; row-- {
		sb.WriteString(generatePieceRow(arr[row]))
		if row > 0 {
			sb.WriteRune('/')
		}
	}

	return sb.String()
}

func generatePieceRow(row [8]Piece) string {
	var sb strings.Builder
	empty := 0
	for _, piece := range row {
		if piece != 0 {
			if empty > 0 {
				sb.WriteRune(rune('0' + empty))
			}
			empty = 0
			sb.WriteRune(piece.GetRune())
			continue
		}

		empty++
	}

	if empty != 0 {
		sb.WriteRune(rune('0' + empty))
	}

	return sb.String()
}

func GenerateBoard(str string) (*[8][8]Piece, error) {
	board := &[8][8]Piece{}
	ranks := strings.Split(str, "/")

	if len(ranks) != 8 {
		return board, fmt.Errorf("invalid number of rows (%d)", len(ranks))
	}

	for rank, row := range ranks {
		result, err := generateBoardRow(row)

		if err != nil {
			return board, err
		}

		board[7-rank] = *result
	}

	return board, nil
}

func generateBoardRow(str string) (*[8]Piece, error) {
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
