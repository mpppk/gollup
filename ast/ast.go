package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"

	"github.com/go-toolsmith/astcopy"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/loader"
)

func NewProgram(fileName string) (*loader.Program, error) {
	lo := &loader.Config{
		Fset:       token.NewFileSet(),
		ParserMode: parser.ParseComments}
	dirPath := filepath.Dir(fileName)
	packages, err := parser.ParseDir(lo.Fset, dirPath, nil, parser.Mode(0))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse dir: "+dirPath)
	}

	var files []*ast.File
	for _, pkg := range packages {
		for _, file := range pkg.Files {
			files = append(files, file)
		}
	}

	lo.CreateFromFiles("main", files...)
	return lo.Load()
}

func ConvertErrorFuncToMustFunc(prog *loader.Program, currentPkg *loader.PackageInfo, orgFuncDecl *ast.FuncDecl) (*ast.FuncDecl, bool) {
	funcDecl := astcopy.FuncDecl(orgFuncDecl)
	if !IsErrorFunc(funcDecl) {
		return nil, false
	}
	results := funcDecl.Type.Results.List
	funcDecl.Type.Results.List = results[:len(results)-1]
	replaceReturnStmtByPanicIfErrorExist(prog, currentPkg, funcDecl)
	addPrefixToFunc(funcDecl, "Must")
	return funcDecl, true
}

func GenerateErrorFuncWrapper(currentPkg *loader.PackageInfo, orgFuncDecl *ast.FuncDecl) (*ast.FuncDecl, bool) {
	funcDecl := astcopy.FuncDecl(orgFuncDecl)
	if !IsErrorFunc(funcDecl) {
		return nil, false
	}

	results := getFuncDeclResults(funcDecl)
	funcDecl.Type.Results.List = funcDecl.Type.Results.List[:len(funcDecl.Type.Results.List)-1]

	wrappedCallExpr := generateCallExpr(extractRecvName(funcDecl), funcDecl.Name.Name, funcDecl.Type.Params.List)
	var lhs []string
	for _, result := range results {
		for _, name := range result.Names {
			lhs = append(lhs, name.Name)
		}
	}

	if len(lhs) == 0 {
		for range results {
			tempValueName := getAvailableValueName(currentPkg.Pkg, "v", funcDecl.Body.Pos())
			lhs = append(lhs, tempValueName)
		}

		tempErrValueName := getAvailableValueName(currentPkg.Pkg, "err", funcDecl.Body.Pos())
		lhs[len(lhs)-1] = tempErrValueName
	}

	funcDecl.Body = &ast.BlockStmt{
		List: []ast.Stmt{
			generateAssignStmt(lhs, wrappedCallExpr),
			generatePanicIfErrorExistStmtAst(lhs[len(lhs)-1]),
			&ast.ReturnStmt{Results: identsToExprs(newIdents(lhs[:len(lhs)-1]))},
		},
	}
	addPrefixToFunc(funcDecl, "Must")
	return funcDecl, true
}

func replaceReturnStmtByPanicIfErrorExist(orgProg *loader.Program, currentPkg *loader.PackageInfo, funcDecl *ast.FuncDecl) *ast.FuncDecl {
	newFuncDeclNode := astutil.Apply(funcDecl, func(cr *astutil.Cursor) bool {
		returnStmt, ok := cr.Node().(*ast.ReturnStmt)
		if !ok {
			return true
		}

		returnResults := returnStmt.Results
		if len(returnResults) == 0 {
			return true
		}

		lastReturnResult := returnResults[len(returnResults)-1]
		returnStmt.Results = returnResults[:len(returnResults)-1]
		if lastReturnResultIdent, ok := lastReturnResult.(*ast.Ident); ok {
			if lastReturnResultIdent.Name == "nil" {
				return true
			}

			panicIfErrExistIfStmt := generatePanicIfErrorExistStmtAst(lastReturnResultIdent.Name)
			cr.InsertBefore(panicIfErrExistIfStmt)
			return true
		}

		if lastReturnResultCallExpr, ok := lastReturnResult.(*ast.CallExpr); ok {
			typeNames, err := getCallExprReturnTypes(orgProg, currentPkg, lastReturnResultCallExpr)
			if err != nil {
				panic(err)
			}
			var lhs []string
			for range typeNames {
				tempValueName := getAvailableValueName(currentPkg.Pkg, "v", lastReturnResultCallExpr.Pos())
				lhs = append(lhs, tempValueName)
			}

			tempErrValueName := getAvailableValueName(currentPkg.Pkg, "err", lastReturnResultCallExpr.Pos())
			lhs[len(lhs)-1] = tempErrValueName

			for _, lh := range lhs[:len(lhs)-1] {
				returnStmt.Results = append(returnStmt.Results, ast.NewIdent(lh))
			}

			assignStmt := generateAssignStmt(lhs, lastReturnResultCallExpr)
			panicIfErrExistIfStmt := generatePanicIfErrorExistStmtAst(tempErrValueName)
			cr.InsertBefore(assignStmt)
			cr.InsertBefore(panicIfErrExistIfStmt)
		}

		return true
	}, nil)
	newFuncDecl := newFuncDeclNode.(*ast.FuncDecl)
	return newFuncDecl
}

func identsToExprs(idents []*ast.Ident) (exprs []ast.Expr) {
	for _, ident := range idents {
		exprs = append(exprs, ast.Expr(ident))
	}
	return
}

func newIdents(identNames []string) (idents []*ast.Ident) {
	for _, identName := range identNames {
		idents = append(idents, &ast.Ident{
			Name: identName,
		})
	}
	return
}

func getAvailableValueName(currentPkg *types.Package, valName string, pos token.Pos) string {
	innerMost := currentPkg.Scope().Innermost(pos)
	s, _ := innerMost.LookupParent(valName, pos)
	if s == nil {
		return valName
	}

	cnt := 0
	valNameWithNumber := fmt.Sprintf("%v%v", valName, cnt)
	for {
		s, _ := innerMost.LookupParent(valNameWithNumber, pos)
		if s == nil {
			return valNameWithNumber
		}
		cnt++
		valNameWithNumber = fmt.Sprintf("%v%v", valName, cnt)
	}
}

func extractRecvName(funcDecl *ast.FuncDecl) string {
	if funcDecl.Recv == nil || len(funcDecl.Recv.List) <= 0 {
		return ""
	}
	names := funcDecl.Recv.List[0].Names
	if len(names) <= 0 {
		panic(fmt.Sprintf("unexpected recv names: %v from %v", names, funcDecl.Name.Name))
	}
	return names[0].Name
}

func MergeImportDeclsFromPackageInfo(packageInfo *loader.PackageInfo) (importDecl *ast.GenDecl) {
	for _, file := range packageInfo.Files {
		imports := ExtractImportDeclsFromDecls(file.Decls)
		if importDecl == nil {
			importDecl = imports[0]
		}
		for _, genDecl := range imports {
			importDecl.Specs = append(importDecl.Specs, genDecl.Specs...)
		}
	}
	return
}

func ExtractImportDeclsFromDecls(decls []ast.Decl) (importDecls []*ast.GenDecl) {
	for _, decl := range decls {
		if importDecl, ok := declToImportDecl(decl); ok {
			importDecls = append(importDecls, importDecl)
		}
	}
	return
}

func declToImportDecl(decl ast.Decl) (*ast.GenDecl, bool) {
	if genDecl, ok := decl.(*ast.GenDecl); ok {
		if genDecl.Tok == token.IMPORT {
			return genDecl, true
		}
	}
	return nil, false
}

func ImportDeclsToDecls(importDecls []*ast.GenDecl) (decls []ast.Decl) {
	for _, decl := range importDecls {
		decls = append(decls, decl)
	}
	return
}
