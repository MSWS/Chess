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
	e, out, _ := base()

	go e.Init()

	expectLine(t, out, "id name GoChess2")
	expectLine(t, out, "id author MS")
	expectLine(t, out, "uciok")
}

func TestPrint(t *testing.T) {
	e, out, _ := base()

	go e.Print("Hello, World!")

	expectLine(t, out, "Hello, World!")
}

func TestCopyProtection(t *testing.T) {
	e, out, _ := base()

	go e.CopyProtection()

	expectLine(t, out, "copyprotection ok")
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
