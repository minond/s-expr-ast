package main

import (
	"fmt"

	"github.com/minond/gong/lexer"
)

func main() {
	fmt.Println(lexer.Lex(` false "one two \"three\" four" 1 2 0b001 0xFFF 32 s true1 true fda  +-*^/ [.] . ##@ `))
	fmt.Println(lexer.Lex(`  3`))
	fmt.Println("\n\n")
	fmt.Println(lexer.Lex(`0b0001 0b0101 &`))
	fmt.Println("\n\n")
	fmt.Println(lexer.Lex(`0b0001      0xF               ^`))
	fmt.Println("\n\n")
	fmt.Println(lexer.Lex(`(define add (lambda (x y) (+ x y)))`))
	fmt.Println("\n\n")
	fmt.Println(lexer.Lex(`  false "one"   `))
	fmt.Println("\n\n")
	fmt.Println(lexer.Lex(`  (a)  `))
	fmt.Println("\n\n")
	fmt.Println(lexer.Lex(` a->b = "c" `))
	fmt.Println("\n\n")
	fmt.Println(lexer.Lex(` 1/3 `))
	fmt.Println("\n\n")
	fmt.Println(lexer.Lex(` 1//3 ; 1 // 3;1    //     2`))
	fmt.Println("\n\n")
	fmt.Println(lexer.Lex(`3.1//0.4`))
}
