package usecase_test

import (
	"testing"

	"github.com/mpppk/cli-template/registry"
)

func TestCalcSum(t *testing.T) {
	type args struct {
		strNumbers []int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "",
			args: args{
				strNumbers: []int{1, 2},
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				strNumbers: []int{-1, 2},
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sumUseCase := registry.InitializeSumUseCase(nil)
			got := sumUseCase.CalcSum(tt.args.strNumbers)
			if got != tt.want {
				t.Errorf("CalcSumFromStringSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalcL1Norm(t *testing.T) {
	type args struct {
		strNumbers []int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "",
			args: args{
				strNumbers: []int{1, 2},
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				strNumbers: []int{-1, 2},
			},
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sumUseCase := registry.InitializeSumUseCase(nil)
			got := sumUseCase.CalcL1Norm(tt.args.strNumbers)
			if got != tt.want {
				t.Errorf("CalcSumFromStringSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}
