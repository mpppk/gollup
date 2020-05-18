package main

import (
	"github.com/mpppk/gollup/testdata/2dmap/lib"
)

func main() {
	m := lib.NewMap()
	m2 := lib.NewMap2()
	m.Get()
	m2.Get()
}
