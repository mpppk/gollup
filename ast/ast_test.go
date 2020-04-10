package ast

import (
	"go/ast"
	"strconv"
	"testing"
)

func TestIsErrorFunc(t *testing.T) {
	type args struct {
		funcDecl *ast.FuncDecl
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "FooBar",
			args: args{
				funcDecl: &ast.FuncDecl{
					Type: &ast.FuncType{
						Results: &ast.FieldList{
							Opening: 0,
							List: []*ast.Field{
								{
									Type: &ast.Ident{
										Name:    "error",
									},
								},
							},
							Closing: 0,
						},
					},
				},
			},
			want: true,
		},
		{
			name: "FooBar",
			args: args{
				funcDecl: &ast.FuncDecl{
					Type: &ast.FuncType{
						Results: &ast.FieldList{
							Opening: 0,
							List: []*ast.Field{
								{
									Type: &ast.Ident{
										Name:    strconv.Quote("int"),
									},
								},
								{
									Type: &ast.Ident{
										Name:    "error",
									},
								},
							},
							Closing: 0,
						},
					},
				},
			},
			want: true,
		},
		{
			name: "FooBar",
			args: args{
				funcDecl: &ast.FuncDecl{
					Type: &ast.FuncType{
						Results: &ast.FieldList{
							Opening: 0,
							List: []*ast.Field{
								{
									Type: &ast.Ident{
										Name:    strconv.Quote("int"),
									},
								},
							},
							Closing: 0,
						},
					},
				},
			},
			want: false,
		},
		{
			name: "FooBar",
			args: args{
				funcDecl: &ast.FuncDecl{
					Type: &ast.FuncType{
						Results: &ast.FieldList{
							Opening: 0,
							List: []*ast.Field{
								{
									Type: &ast.Ident{
										Name:    "error",
									},
								},
								{
									Type: &ast.Ident{
										Name:    "int",
									},
								},
							},
							Closing: 0,
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsErrorFunc(tt.args.funcDecl); got != tt.want {
				t.Errorf("IsErrorFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
