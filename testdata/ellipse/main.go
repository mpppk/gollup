package main

import (
	"fmt"

	"github.com/mpppk/gollup/testdata/ellipse/lib"
)

func main() {
	fmt.Println(lib.MustMinInt64(1, 2), lib.MustMinInt64([]int64{1, 2}...))
}
