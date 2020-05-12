package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/mpppk/gollup/testdata/atc001A/lib"
)

func solve(N int64, h []int64) int64 {
	m := lib.Int64Map(map[int64]int64{})
	m[0] = 0
	m[1] = h[0]
	for i := int64(1); i < N; i++ {
		m.ChMin(i+1, m[i]+h[i])
		m.ChMin(i+1, m[i-1]+h[i])
	}
	return m.MustGet(N - 1)
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
	h := make([]int64, N)
	for i := int64(0); i < N; i++ {
		scanner.Scan()
		h[i], _ = strconv.ParseInt(scanner.Text(), 10, 64)
	}
	fmt.Println(solve(N, h))
}
