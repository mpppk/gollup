package main

import (
	"fmt"

	"github.com/mpppk/gollup/testdata/comments/lib"
)

const ANSWER = 42

func main() {
	fmt.Println(F1(), lib.F1())
}

// F1 is function
func F1() int {
	return f()
}

// f is function
func f() int {
	return ANSWER
}
