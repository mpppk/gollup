package lib

type Map map[int]string
type Map2 map[int]string

func (m Map) Get() string {
	return "a"
}

func (m Map2) Get() string {
	return "a"
}

func NewMap() Map {
	return Map{}
}

func NewMap2() Map2 {
	return Map2{}
}
