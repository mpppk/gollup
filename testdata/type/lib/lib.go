package lib

type Int int64
type S map[int64]Int

func (i Int) Get() int {
	return 1
}

func (S *S) F() int {
	return 1
}
