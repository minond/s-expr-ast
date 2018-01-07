package main

import (
	"fmt"

	"github.com/minond/gong/gong"
)

func main() {
	toks := gong.Scan(`

(define atom?
  (lambda (x)
    (and (not (null? x))
         (not (pair? x)))))

`)

	statements := gong.Parse(toks)

	for _, statement := range statements {
		fmt.Println(statement)
	}
}
