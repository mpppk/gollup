package lib

import "errors"

func MinInt64(values ...int64) (min int64, err error) {
	if len(values) == 0 {
		return 0, errors.New("empty slice is given")
	}
	min = values[0]
	for _, value := range values {
		if min > value {
			min = value
		}
	}
	return
}

func MustMinInt64(values ...int64) (min int64) {
	min, err := MinInt64(values...)
	if err != nil {
		panic(err)
	}
	return min
}
