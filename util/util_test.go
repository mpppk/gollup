package util

import (
	"reflect"
	"testing"
)

func TestConvertStringSliceToIntSlice(t *testing.T) {
	type args struct {
		stringSlice []string
	}
	tests := []struct {
		name         string
		args         args
		wantIntSlice []int
		wantErr      bool
	}{
		{
			name: "can convert string slice to int slice",
			args: args{
				stringSlice: []string{"1", "2", "3"},
			},
			wantIntSlice: []int{1, 2, 3},
			wantErr:      false,
		},
		{
			name: "will be error if string can not be convert to number",
			args: args{
				stringSlice: []string{"1", "2", "a"},
			},
			wantIntSlice: nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIntSlice, err := ConvertStringSliceToIntSlice(tt.args.stringSlice)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertStringSliceToIntSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIntSlice, tt.wantIntSlice) {
				t.Errorf("ConvertStringSliceToIntSlice() = %v, want %v", gotIntSlice, tt.wantIntSlice)
			}
		})
	}
}
