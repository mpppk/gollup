package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/mpppk/gollup/testdata/atc001B/lib"
)

const YES = "Yes"
const NO = "No"

func solve(N int64, Q int64, P []int64, A []int64, B []int64) {
	unionFind := lib.NewUnionFindInt(lib.MustIntRange(0, 100001, 1))
	for i := 0; i < int(Q); i++ {
		p, a, b := P[i], int(A[i]), int(B[i])

		if p == 0 { // 連結クエリ
			unionFind.Unite(a, b)
		} else if p == 1 { // 判定クエリ
			fmt.Println(lib.TernaryOPString(unionFind.IsSameGroup(a, b), YES, NO))
			//lib.TernaryOPString(unionFind.IsSameGroup(a, b), "Yes", "No")
		}
	}
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
	var Q int64
	scanner.Scan()
	Q, _ = strconv.ParseInt(scanner.Text(), 10, 64)
	P := make([]int64, Q)
	A := make([]int64, Q)
	B := make([]int64, Q)
	for i := int64(0); i < Q; i++ {
		scanner.Scan()
		P[i], _ = strconv.ParseInt(scanner.Text(), 10, 64)
		scanner.Scan()
		A[i], _ = strconv.ParseInt(scanner.Text(), 10, 64)
		scanner.Scan()
		B[i], _ = strconv.ParseInt(scanner.Text(), 10, 64)
	}
	solve(N, Q, P, A, B)
}
