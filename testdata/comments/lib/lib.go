package lib

import "math"

// F1 is function
func F1() float64 {
	return f2()
}

// F2 is function
func f2() float64 {
	return math.Sqrt(42)
}
