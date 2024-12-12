package board

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
