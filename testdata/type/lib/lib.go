package lib

type Int int64
type S map[int64]Int
type M map[int64]S

func (i Int) Get() int {
	return 1
}

func (S S) Get() int {
	return 1
}

func (m M) Get() int {
	return 1
}
