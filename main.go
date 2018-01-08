package main

import (
	"fmt"

	"github.com/minond/gong/gong"
)

func main() {
	statements := gong.Parse(`

(define atom?
  (lambda (x)
    (and (not (null? x))
         (not (pair? x)))))

`)

	for _, statement := range statements {
		fmt.Println(statement)
	}
}
