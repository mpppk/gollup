package lib

type S struct{}

func NewS() *S {
	return &S{}
}

func (S *S) F() int {
	return 1
}
