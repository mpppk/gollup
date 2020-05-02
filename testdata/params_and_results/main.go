package main

import (
	"fmt"

	"github.com/mpppk/gollup/testdata/params_and_results/lib"
)

func main() {
	s1 := lib.NewS()
	fmt.Println(F1(s1))

	s2 := lib.S{}
	fmt.Println(F2(s2))
}

func F1(s *lib.S) int {
	return s.F()
}
func F2(s lib.S) int {
	return s.F()
}
