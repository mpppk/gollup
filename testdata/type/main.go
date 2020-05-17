package main

import (
	"fmt"

	"github.com/mpppk/gollup/testdata/type/lib"
)

func main() {
	m := lib.S(map[int64]lib.Int{})
	m2 := lib.M{}
	fmt.Println(m[0].Get(), m2.Get(), m2[0].Get())
	fmt.Println(m2.Get())
}
