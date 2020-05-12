package lib

import "fmt"

// Int64Map は、map[int64][int64]に便利メソッドを追加します.
type Int64Map map[int64]int64

// MustGetは、指定したkeyの値を返します. 指定したkeyの値が存在しない場合panicします.
func (m Int64Map) MustGet(key int64) int64 {
	v, ok := m[key]
	if !ok {
		panic(fmt.Sprintf("ivnalid key is specfied in Int64Map: %v", key))
	}
	return v
}

// ChMin は、与えられた値が既に存在する値よりも小さければ代入します.
// 指定したkeyの値が存在しない場合も代入します. この場合、2つめの戻り値はfalseになります.
func (m Int64Map) ChMin(key, value int64) (replaced bool, valueAlreadyExist bool) {
	if v, ok := m[key]; ok {
		if v > value {
			m[key] = value
			return true, true
		} else {
			return false, true
		}
	}
	m[key] = value
	return true, false
}
