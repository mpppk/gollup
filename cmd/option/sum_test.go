package option_test

import (
	"testing"

	"github.com/mpppk/cli-template/cmd/option"
)

func TestSumCmdConfig_HasOut(t *testing.T) {
	type fields struct {
		Norm bool
		Out  string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "should return true if Out is not default value",
			fields: fields{
				Norm: false,
				Out:  "test.txt",
			},
			want: true,
		},
		{
			name: "should return false if out is default value",
			fields: fields{
				Norm: false,
				Out:  "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &option.SumCmdConfig{
				Norm: tt.fields.Norm,
				Out:  tt.fields.Out,
			}
			if got := c.HasOut(); got != tt.want {
				t.Errorf("SumCmdConfig.HasOut() = %v, want %v", got, tt.want)
			}
		})
	}
}
