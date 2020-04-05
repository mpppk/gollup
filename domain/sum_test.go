package domain_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/mpppk/cli-template/domain"

	"github.com/mpppk/cli-template/domain/model"

	"github.com/mpppk/cli-template/util"
)

func TestNewNumbers(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want model.Numbers
	}{
		{
			name: "",
			args: args{
				nums: []int{},
			},
			want: []int{},
		},
		{
			name: "",
			args: args{
				nums: []int{1},
			},
			want: []int{1},
		},
		{
			name: "",
			args: args{
				nums: []int{1, 2},
			},
			want: []int{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := model.NewNumbers(tt.args.nums); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNumbersFromStringSlice(t *testing.T) {
	type args struct {
		strNumbers []string
	}
	tests := []struct {
		name    string
		args    args
		want    model.Numbers
		wantErr bool
	}{
		{
			args: args{
				strNumbers: []string{"1"},
			},
			want:    []int{1},
			wantErr: false,
		},
		{
			args: args{
				strNumbers: []string{"1", "2"},
			},
			want:    []int{1, 2},
			wantErr: false,
		},
		{
			args: args{
				strNumbers: []string{"-1", "2"},
			},
			want:    []int{-1, 2},
			wantErr: false,
		},
		{
			args: args{
				strNumbers: []string{"1", "a"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewNumbersFromStringSlice(tt.args.strNumbers)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNumbersFromStringSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNumbersFromStringSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumbers_CalcL1Norm(t *testing.T) {
	tests := []struct {
		name string
		n    model.Numbers
		want int
	}{
		{
			n:    []int{1, 2},
			want: 3,
		},
		{
			n:    []int{-1, 2},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.CalcL1Norm(); got != tt.want {
				t.Errorf("CalcL1Norm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumbers_CalcSum(t *testing.T) {
	tests := []struct {
		name string
		n    model.Numbers
		want int
	}{
		{
			n:    []int{1, 2},
			want: 3,
		},
		{
			n:    []int{-1, 2},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.CalcSum(); got != tt.want {
				t.Errorf("CalcSum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleSum() {
	numbers := []int{1, -2, 3}
	fmt.Println(model.Sum(numbers))
	// Output:
	// 2
}

func TestSum(t *testing.T) {
	type args struct {
		numbers []int
	}
	tests := []struct {
		name    string
		args    args
		wantSum int
	}{
		{
			name: "return sum of numbers",
			args: args{
				numbers: []int{1, 2, 3},
			},
			wantSum: 6,
		},
		{
			name: "return sum of numbers",
			args: args{
				numbers: []int{1, -2, 3},
			},
			wantSum: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSum := model.Sum(tt.args.numbers); gotSum != tt.wantSum {
				t.Errorf("Sum() = %v, want %v", gotSum, tt.wantSum)
			}
		})
	}
}

func TestSumFromFile(t *testing.T) {
	type args struct {
		numbers []int
	}

	type tCase struct {
		name    string
		args    args
		wantSum int
	}

	f, err := os.Open("../testdata/sum.txt")
	if err != nil {
		t.Fatal(err)
	}
	contentBytes, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	contents := string(contentBytes)
	lines := strings.Split(strings.Replace(contents, "\r\n", "\n", -1), "\n")
	var tests []tCase
	for i, line := range lines {
		strRow := strings.Split(line, " ")
		row, err := util.ConvertStringSliceToIntSlice(strRow)
		if err != nil {
			t.Fatal(err)
		}
		want := row[len(row)-1]
		nums := row[:len(row)-1]
		tc := tCase{
			name: fmt.Sprintf("case%d", i),
			args: args{
				numbers: nums,
			},
			wantSum: want,
		}
		tests = append(tests, tc)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSum := model.Sum(tt.args.numbers); gotSum != tt.wantSum {
				t.Errorf("Sum() = %v, want %v", gotSum, tt.wantSum)
			}
		})
	}
}

func TestSumFromString(t *testing.T) {
	type args struct {
		stringNumbers []string
	}
	tests := []struct {
		name    string
		args    args
		wantSum int
		wantErr bool
	}{
		{
			name: "return sum of numbers",
			args: args{
				stringNumbers: []string{"1", "2", "3"},
			},
			wantSum: 6,
			wantErr: false,
		},
		{
			name: "will be error if args includes not number string",
			args: args{
				stringNumbers: []string{"1", "2", "a"},
			},
			wantSum: 0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSum, err := model.SumFromString(tt.args.stringNumbers)
			if (err != nil) != tt.wantErr {
				t.Errorf("SumFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSum != tt.wantSum {
				t.Errorf("SumFromString() = %v, want %v", gotSum, tt.wantSum)
			}
		})
	}
}

func ExampleL1Norm() {
	numbers := []int{1, -2, 3}
	fmt.Println(model.L1Norm(numbers))
	// Output:
	// 6
}

func TestL1Norm(t *testing.T) {
	type args struct {
		numbers []int
	}
	tests := []struct {
		name    string
		args    args
		wantSum int
	}{
		{
			name: "return sum of numbers",
			args: args{
				numbers: []int{1, 2, 3},
			},
			wantSum: 6,
		},
		{
			name: "return sum of numbers",
			args: args{
				numbers: []int{1, -2, 3},
			},
			wantSum: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSum := model.L1Norm(tt.args.numbers); gotSum != tt.wantSum {
				t.Errorf("Sum() = %v, want %v", gotSum, tt.wantSum)
			}
		})
	}
}
