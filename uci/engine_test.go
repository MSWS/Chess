package uci

import (
	"io"
	"strings"
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

	expectLine(t, in, "id name GoChess2")
	expectLine(t, in, "id author MS")
	expectLine(t, in, "uciok")
}

func expectLine(t *testing.T, in io.Reader, expected string) {
	line, err := readLine(in)
	if err != nil {
		t.Fatalf("Error reading line: %v", err)
	}

	if line != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, line)
	}
}

func readLine(in io.Reader) (string, error) {
	buf := make([]byte, 1024)
	n, err := in.Read(buf)
	if err != nil {
		return "", err
	}

	return strings.Trim(string(buf[:n]), "\n"), nil
}
