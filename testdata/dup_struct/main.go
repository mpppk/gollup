package main

import (
	"fmt"

	"github.com/mpppk/gollup/testdata/struct/lib"
)

type S struct{}

func (S *S) F() int {
	return 1
}

func main() {
	s1 := S{}
	s2 := lib.S{}
	fmt.Println(s1.F(), s2.F())
}
