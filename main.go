package main

import (
	"bufio"
	"os"

	"github.com/msws/chess/uci"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	engine := uci.NewEngine(reader, os.Stdout)

	for {
		cmd, err := reader.ReadString('\n')

		if err != nil {
			panic(err)
		}

		go engine.Print(cmd)
	}
}
