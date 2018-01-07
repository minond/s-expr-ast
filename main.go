package main

import (
	"fmt"

	"github.com/minond/gong/gong"
)

func main() {
	toks := gong.Scan(`(f->xyz 1 2 3 456 0.5 0b001 0x32 "My Name Is Marcos M\"inond")`)

	for _, tok := range toks {
		fmt.Println(tok)
	}
}
