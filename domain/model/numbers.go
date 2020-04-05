package model

import (
	"github.com/mpppk/cli-template/util"
	"math"
)

// Numbers represents numbers
type Numbers []int

// NewNumbers is constructor for Numbers
func NewNumbers(nums []int) Numbers {
	return nums
}

// CalcSum calc sum of numbers
func (n Numbers) CalcSum() int {
	return Sum(n)
}

// CalcL1Norm calc L1 norm of numbers
func (n Numbers) CalcL1Norm() int {
	return L1Norm(n)
}

// Sum returns sum of numbers
func Sum(numbers []int) (sum int) {
	for _, number := range numbers {
		sum += number
	}
	return
}

// L1Norm returns L1 norm of numbers
func L1Norm(numbers []int) (l1norm int) {
	var absNumbers []int
	for _, number := range numbers {
		absNumbers = append(absNumbers, int(math.Abs(float64(number))))
	}
	return Sum(absNumbers)
}

// SumFromString returns sum numbers which be converted from strings
func SumFromString(stringNumbers []string) (sum int, err error) {
	numbers, err := util.ConvertStringSliceToIntSlice(stringNumbers)
	if err != nil {
		return 0, err
	}
	return Sum(numbers), nil
}
