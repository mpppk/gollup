// Package util provides some utilities
package util

import (
	"fmt"
	"github.com/comail/colog"
	"strconv"
)

// ConvertStringSliceToIntSlice converts string slices to int slices.
func ConvertStringSliceToIntSlice(stringSlice []string) (intSlice []int, err error) {
	for _, s := range stringSlice {
		num, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("failed to convert string slice to int slice: %w", err)
		}
		intSlice = append(intSlice, num)
	}
	return
}

// InitializeLog initialize log settings
func InitializeLog(verbose bool) {
	colog.Register()
	colog.SetDefaultLevel(colog.LDebug)
	colog.SetMinLevel(colog.LInfo)

	if verbose {
		colog.SetMinLevel(colog.LDebug)
	}
}
