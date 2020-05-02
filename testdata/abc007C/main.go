package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/mpppk/gollup/testdata/abc007C/lib"
)

func solve(input *lib.Input) int {
	startIndices := input.MustGetIntLine(1)
	startRowIndex, startColIndex := startIndices[0]-1, startIndices[1]-1
	endIndices := input.MustGetIntLine(2)
	endRowIndex, endColIndex := endIndices[0]-1, endIndices[1]-1
	m := input.MustReadAsStringGridFrom(3)
	minStepMap := lib.NewIntGridMap(len(m), len(m[0]), math.MaxInt32)
	step := 0

	l := list.New()
	l.PushBack([2]int{startRowIndex, startColIndex})
	minStepMap[startRowIndex][startColIndex] = 0

	for l.Len() > 0 {
		pos, ok := l.Remove(l.Front()).([2]int)
		if !ok {
			panic("invalid pos")
		}

		step = minStepMap[pos[0]][pos[1]]

		if pos[0] == endRowIndex && pos[1] == endColIndex {
			return step
		}

		s := m[pos[0]][pos[1]]
		if s == "#" {
			continue
		}

		nextX, nextY := pos[0]-1, pos[1]
		if pos[0] > 0 && step+1 < minStepMap[nextX][nextY] {
			minStepMap[nextX][nextY] = step + 1
			l.PushBack([2]int{nextX, nextY})
		}

		nextX, nextY = pos[0]+1, pos[1]
		if len(m) > pos[0]+1 && step+1 < minStepMap[nextX][nextY] {
			minStepMap[nextX][nextY] = step + 1
			l.PushBack([2]int{nextX, nextY})
		}

		nextX, nextY = pos[0], pos[1]-1
		if pos[1] > 0 && step+1 < minStepMap[nextX][nextY] {
			minStepMap[nextX][nextY] = step + 1
			l.PushBack([2]int{nextX, nextY})
		}

		nextX, nextY = pos[0], pos[1]+1
		if len(m[0]) > pos[1]+1 && step+1 < minStepMap[nextX][nextY] {
			minStepMap[nextX][nextY] = step + 1
			l.PushBack([2]int{nextX, nextY})
		}
	}
	panic("no route")
}

func main() {
	input := lib.MustNewInputFromReader(bufio.NewReader(io.Reader(os.Stdin)))
	fmt.Println(solve(input))
}
