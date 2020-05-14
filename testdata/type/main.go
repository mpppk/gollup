package main

import (
	"fmt"

	"github.com/mpppk/gollup/testdata/type/lib"
)

func main() {
	m := lib.S(map[int64]lib.Int{})
	fmt.Println(m.F(), m[0].Get())
}
