package main

import (
	"fmt"

	"github.com/minond/gong/lexer"
)

func main() {
	fmt.Println(lexer.Lex(`"one two \"three\" four" 1 2 0b001 0xFFF 32 s true1 true fda    . `))
}
