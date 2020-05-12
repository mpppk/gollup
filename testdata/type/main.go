package main

import (
	"fmt"

	"github.com/mpppk/gollup/testdata/type/lib"
)

func main() {
	m := lib.S(map[int64]int64{})
	fmt.Println(m.F())
}
