package main

import (
	"fmt"

	"github.com/mpppk/gollup/testdata/const/lib"
)

const ANSWER = 42

func main() {
	fmt.Println(lib.F1() + ANSWER)
}
