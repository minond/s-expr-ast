package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
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

func stdio() {
	scanner := bufio.NewScanner(os.Stdin)
	content := ""

	for scanner.Scan() {
		content += "\n" + scanner.Text()
	}

	show(strings.TrimSpace(content))
}

func file(path string) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Printf("Error reading file: %v", err)
		os.Exit(2)
	} else {
		show(strings.TrimSpace(string(content)))
	}
}

func main() {
	if len(os.Args) > 1 {
		file(os.Args[1])
	} else if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
		stdio()
	} else {
		repl()
	}
}
