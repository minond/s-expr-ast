package main

import (
	"fmt"
)

func main() {
	statements := Parse(`

(define atom?
  (lambda (x)
    (and (not (null? x))
         (not (pair? x)))))

`)

	for _, statement := range statements {
		fmt.Println(statement)
	}
}
