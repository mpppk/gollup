package lib

import "math"

func F1() float64 {
	return f2()
}

func f2() float64 {
	return math.Sqrt(42)
}
