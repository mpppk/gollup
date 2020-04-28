package main

import (
	"fmt"

	"github.com/mpppk/gollup/testdata/multi_pkg/lib"
)

const ANSWER = 42

func main() {
	fmt.Println(F1(), lib.F1())
}

func F1() int {
	return f()
}

func f() int {
	return ANSWER
}
