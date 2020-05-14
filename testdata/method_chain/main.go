package main

import (
	"fmt"

	"github.com/mpppk/gollup/testdata/method_chain/lib"
)

func main() {
	fmt.Println(lib.F().Get())
}
