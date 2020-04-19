package main

import (
	"fmt"

	"github.com/mpppk/gollup/testdata/test3/lib"
)

const ANSWER = 42

func main() {
	s := lib.S{}
	fmt.Println(F1(), s.F())
}

func F1() int {
	return ANSWER
}
