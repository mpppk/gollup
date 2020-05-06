package lib

import (
	"fmt"
)

type UnionFindRune struct {
	nodes map[rune]rune
}

func (u *UnionFindRune) GetRoot(value rune) (rune, int) {
	v := value
	newV := u.nodes[v]
	cnt := 0
	for newV != v {
		cnt++
		oldV := v
		v = newV
		newV = u.nodes[newV]
		u.nodes[oldV] = newV
	}
	return newV, cnt
}

func (u *UnionFindRune) IsSameGroup(v1, v2 rune) bool {
	v1Root, _ := u.GetRoot(v1)
	v2Root, _ := u.GetRoot(v2)
	return v1Root == v2Root
}

func (u *UnionFindRune) Unite(v1, v2 rune) (rune, bool) {
	v1Root, v1HopNum := u.GetRoot(v1)
	v2Root, v2HopNum := u.GetRoot(v2)
	if v1Root == v2Root {
		return v1Root, false
	}
	if v1HopNum >= v2HopNum {
		u.nodes[v2Root] = v1Root
		return v1Root, true
	}
	u.nodes[v1Root] = v2Root
	return v2Root, true
}

type UnionFindInt struct {
	nodes map[int]int
}

func NewUnionFindInt(values []int) *UnionFindInt {
	m := map[int]int{}
	for _, v := range values {
		m[v] = v
	}
	return &UnionFindInt{nodes: m}
}

// GetRoot は、与えられた値の根の値を返します.
func (u *UnionFindInt) GetRoot(value int) (int, int) {
	v := value
	newV := u.nodes[v]
	cnt := 0
	for newV != v {
		cnt++
		oldV := v
		v = newV
		newV = u.nodes[newV]
		u.nodes[oldV] = newV
	}
	return newV, cnt
}

// Unite は、v1とv2のグループをマージします.
func (u *UnionFindInt) Unite(v1, v2 int) (int, bool) {
	v1Root, v1HopNum := u.GetRoot(v1)
	v2Root, v2HopNum := u.GetRoot(v2)
	if v1Root == v2Root {
		return v1Root, false
	}
	if v1HopNum >= v2HopNum {
		u.nodes[v2Root] = v1Root
		return v1Root, true
	}
	u.nodes[v1Root] = v2Root
	return v2Root, true
}

// IsSameGroup は、v1とv2が同じグループに所属しているかを返します.
func (u *UnionFindInt) IsSameGroup(v1, v2 int) bool {
	v1Root, _ := u.GetRoot(v1)
	v2Root, _ := u.GetRoot(v2)
	return v1Root == v2Root
}

func IntRange(start, end, step int) ([]int, error) {
	if end < start {
		return nil, fmt.Errorf("end(%v) is bigger than start(%v)", end, start)
	}
	s := make([]int, 0, int(1+(end-start)/step))
	for start < end {
		s = append(s, start)
		start += step
	}
	return s, nil
}

func MustIntRange(start, end, step int) []int {
	_v0, _err := IntRange(start, end, step)
	if _err != nil {
		panic(_err)
	}
	return _v0
}

func TernaryOPString(ok bool, v1, v2 string) string {
	if ok {
		return v1
	}
	return v2
}
