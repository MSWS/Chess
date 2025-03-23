package testdata

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/msws/chess/board"
)

type TestData struct {
	Description string      `json:"description"`
	Tests       []TestCases `json:"testCases"`
}

type TestCases struct {
	Start    TestStart         `json:"start"`
	Expected []TestExpectation `json:"expected"`
}

type TestStart struct {
	Fen         string `json:"fen"`
	Description string `json:"description"`
}

type TestExpectation struct {
	Move string `json:"move"`
	Fen  string `json:"fen"`
}

func TestJSON(t *testing.T) {
	data, err := os.ReadDir(".")

	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range data {
		t.Run(entry.Name(), func(t *testing.T) {
			jsonData, err := os.Open(entry.Name())
			if err != nil {
				t.Error(err)
				return
			}

			defer jsonData.Close()

			byteData, err := io.ReadAll(jsonData)
			if err != nil {
				t.Error(err)
				return
			}
			testJson(t, byteData)
		})
	}
}

func testJson(t *testing.T, bytes []byte) {
	var data TestData
	json.Unmarshal(bytes, &data)
	t.Log(data.Description)
	testData(t, data)
}

func testData(t *testing.T, data TestData) {
	for _, test := range data.Tests {
		t.Log(test.Start.Description)
		t.Run(strings.ReplaceAll(test.Start.Fen, "/", "."), func(t *testing.T) {
			testCase(t, test.Start, test.Expected)
		})
	}
}

func testCase(t *testing.T, start TestStart, expected []TestExpectation) {
	t.Run("MoveCount", func(t *testing.T) {
		testMoveCount(t, start, expected)
	})
}

func testMoveCount(t *testing.T, start TestStart, expected []TestExpectation) {
	board, err := board.FromFEN(start.Fen)

	if err != nil {
		t.Error(err)
		return
	}

	moves := board.GetMoves()

	if len(moves) != len(expected) {
		t.Errorf("expected %d moves, got %d", len(expected), len(moves))
	}
}
