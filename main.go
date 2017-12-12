package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/minond/gong/gong"
)

func show(prog string) {
	fmt.Println(gong.Lex(prog))
	fmt.Print("\n\n")
}

func main() {
	// show(` false "one two \"three\" four" 1 2 0b001 0xFFF 32 s true1 true fda  +-*^/ [.] . ##@ `)
	// show(`  3`)
	// show(`0b0001 0b0101 &`)
	// show(`0b0001      0xF               ^`)
	// show(`(define add (lambda (x y) (+ x y)))`)
	// show(`  false "one"   `)
	// show(`  (a)  `)
	// show(` a->b = "c" `)
	// show(` 1/3 `)
	// show(` 1//3 ; 1 // 3;1    //     2`)
	// show(`3.1//0.4`)
	// show(`{1 2 3 4} [2 * print] map`)
	//
	// show(`1 2
	// 3
	// 4
	//
	//
	// 5`)
	// show(`one "tw\"o" 'thr\'ee' '"""""""' "''''''''''"`)
	// show(`iiii"i`)
	// show(`0b0001 1x2`)
	// show(`true || false ; true | false ; true && true ; true & true;`)
	// show(`true||false;true|false;true&&true;true&true;`)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		show(text)
	}
}
