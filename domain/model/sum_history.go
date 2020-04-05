package model

import "time"

// SumHistory represents history of sum calculation
type SumHistory struct {
	IsNorm  bool
	Date    time.Time
	Numbers Numbers
	Result  int
}
