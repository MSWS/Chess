package uci

import (
	"io"

	"github.com/msws/chess/board"
)

type Engine interface {
	Debug(bool bool)
	Init()
	CopyProtection()
	SetOption(name string, id string, value string)
	Print(message string)
	Position(fen string, moves []string)
}

type MyEngine struct {
	in    io.Reader
	out   io.Writer
	game  *board.Game
	debug bool
}

// CopyProtection implements Engine.
func (e *MyEngine) CopyProtection() {
	e.Println("copyprotection ok")
}

func (e *MyEngine) Init() {
	e.Println("id name GoChess2")
	e.Println("id author MS")
	e.Println("uciok")
}

func NewEngine(in io.Reader, out io.Writer) Engine {
	return &MyEngine{
		in:  in,
		out: out,
	}
}

func (e *MyEngine) Print(message string) {
	e.out.Write([]byte(message))
}

func (e *MyEngine) Println(message string) {
	e.Print(message + "\n")
}

func (e *MyEngine) Debug(bool bool) {
	e.debug = bool
}

func (e *MyEngine) SetOption(name string, id string, value string) {
}

func (e *MyEngine) Position(fen string, moves []string) {
	game, err := board.FromFEN(fen)
	if err != nil {
		if e.debug {
			e.Print(err.Error())
		}
		return
	}
	e.game = game
}
