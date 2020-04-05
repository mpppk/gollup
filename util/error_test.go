package util_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/mpppk/cli-template/util"
)

func TestPrettyPrintError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				err: errors.New("sample error"),
			},
			want: fmt.Sprintln("Error: sample error"),
		},
		{
			name: "",
			args: args{
				err: fmt.Errorf("a: %w", errors.New("b")),
			},
			want: fmt.Sprintln("Error: a") + fmt.Sprintln("  b"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.PrettyPrintError(tt.args.err); got != tt.want {
				t.Errorf("PrettyPrintError() = %v, want %v", got, tt.want)
			}
		})
	}
}
