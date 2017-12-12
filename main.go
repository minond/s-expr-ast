package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/minond/gong/gong"
)

func show(prog string) {
	fmt.Println("\n  Tokens:\n")
	for _, tok := range gong.Lex(prog) {
		fmt.Printf("    - %s\n", tok)
	}

	fmt.Print("\n")
}

func repl() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')

		switch strings.TrimSpace(text) {
		case "quit":
			fallthrough
		case "exit":
			return

		default:
			show(strings.TrimSpace(text))
		}
	}
}

func main() {
	repl()
}
