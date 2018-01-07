package main

import (
	"fmt"

	"github.com/minond/gong/gong"
)

func main() {
	toks := gong.Scan(`

(define atom
  (lambda (x)
    (and (not (null x))
         (not (pair x)))))

`)

	stms := gong.Parse(toks)

	// for _, tok := range toks {
	// 	fmt.Println(tok)
	// }

	fmt.Println(stms)
}
