package ast

import (
	"fmt"
	"go/ast"
)

func extractFuncLastResultIdent(funcDecl *ast.FuncDecl) (*ast.Ident, bool) {
	expr, ok := extractFuncLastResultExpr(funcDecl)
	if !ok {
		return nil, false
	}
	ident, ok := expr.(*ast.Ident)
	return ident, ok
}

func extractFuncLastResultExpr(funcDecl *ast.FuncDecl) (ast.Expr, bool) {
	if funcDecl == nil {
		panic(fmt.Sprintf("funcDecl is nil"))
	}
	if funcDecl.Type == nil {
		panic(fmt.Sprintf("funcDecl.Type is nil: %v", funcDecl.Name))
	}
	if funcDecl.Type.Results == nil {
		return nil, false
	}
	results := funcDecl.Type.Results.List
	if len(results) == 0 {
		return nil, false
	}
	return results[len(results)-1].Type, true
}
