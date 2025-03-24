package uci

import (
	"io"
	"testing"
)

func base() (Engine, io.Reader, io.Writer) {
	in, out := io.Pipe()
	e := NewEngine(in, out)

	return e, in, out
}

func TestInit(t *testing.T) {
	e, in, _ := base()

	go e.Init()

	buf := make([]byte, 1024)
	n, err := in.Read(buf)

	if err != nil {
		t.Error(err)
	}

	if n == 0 {
		t.Error("No data read")
	}

	if string(buf[:n]) != "id name GoChess2\n" {
		t.Error("Unexpected output, got:", string(buf[:n]))
	}
}
