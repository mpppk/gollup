package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/mpppk/gollup/testdata/compro/lib"
)

func solve(N int64) string {
	fmt.Println(lib.MaxInt(1, 2))
	return ""
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	const initialBufSize = 4096
	const maxBufSize = 1000000
	scanner.Buffer(make([]byte, initialBufSize), maxBufSize)
	scanner.Split(bufio.ScanWords)
	var N int64
	scanner.Scan()
	N, _ = strconv.ParseInt(scanner.Text(), 10, 64)
	fmt.Println(solve(N))
}
