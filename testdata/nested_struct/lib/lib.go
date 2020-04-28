package lib

type S struct {
	S2 *S2
}

type S2 struct {
	Num int
}

func NewS() *S {
	return &S{
		S2: &S2{1},
	}
}

func (S *S) F() int {
	return S.S2.Num
}
