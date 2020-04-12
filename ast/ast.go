package ast

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/mpppk/gollup/util"

	"github.com/go-toolsmith/astcopy"
)

type Packages struct {
	M map[string]*packages.Package
}

func NewProgramFromPackages(packageNames []string) ([]*packages.Package, error) {
	config := &packages.Config{
		Mode: packages.NeedCompiledGoFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.LoadAllSyntax,
	}
	pkgs, err := packages.Load(config, packageNames...)
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, errors.New("error occurred in NewProgramFromPackages")
	}
	return pkgs, nil
}

func listGoFilesFromDirs(dirs []string) (filePaths []string, err error) {
	for _, dir := range dirs {
		fps, err := listGoFiles(dir)
		if err != nil {
			return nil, err
		}
		filePaths = append(filePaths, fps...)
	}
	return
}
func listGoFiles(dir string) (filePaths []string, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if !strings.HasSuffix(path, ".go") ||
			strings.Contains(path, "_test") {
			return nil
		}

		filePaths = append(filePaths, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return
}

func MergeImportDeclsFromPackageInfo(files []*ast.File) (importDecl *ast.GenDecl) {
	for _, file := range files {
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

func NewMergedFileFromPackageInfo(files []*ast.File) *ast.File {
	importDecl := MergeImportDeclsFromPackageInfo(files)

	var imports []*ast.ImportSpec
	for _, file := range files {
		imports = append(imports, file.Imports...)
	}
	return &ast.File{
		Name: &ast.Ident{
			Name: "main",
		},
		Decls:      []ast.Decl{importDecl},
		Scope:      nil,
		Imports:    imports,
		Unresolved: nil,
		Comments:   nil,
	}

}

func CopyFuncDeclsAsDecl(funcDecls []*ast.FuncDecl) (newFuncDecls []ast.Decl) {
	for _, decl := range funcDecls {
		newFuncDecls = append(newFuncDecls, astcopy.FuncDecl(decl))
	}
	return
}

func ExtractCalledFuncsFromFuncDeclRecursive(files []*ast.File, info *types.Info, funcName string, foundedFuncNames []string) (funcs []*types.Func, funcDecls []*ast.FuncDecl, err error) {
	funcDecl := findFuncDeclByName(files, funcName)
	if funcDecl == nil {
		return nil, nil, errors.New("specified function is not found: " + funcName)
	}

	funcDecls = append(funcDecls, funcDecl)
	calledFuncs := extractCalledFuncsFromFuncDecl(info, funcDecl)
	newFoundedFuncs := make([]string, len(foundedFuncNames))
	copy(newFoundedFuncs, foundedFuncNames)
	newFoundedFuncs = append(newFoundedFuncs, funcDecl.Name.Name)

	for _, f := range calledFuncs {

		// TODO: 他のpackageでも組み込みでなければ探索する
		// 他のライブラリの場合どうなるか要確認(多分動くと思っている)
		// 最終的には別リポジトリのコードでもバンドルできるようにしたい
		if util.IsStandardPackage(f.Pkg().Name()) {
			continue
		}

		// 既に発見済みの関数の場合はスキップ
		if isFuncName(foundedFuncNames, f.Name()) {
			continue
		}

		funcs = append(funcs, f)

		newFuncs, newFuncDecls, err := ExtractCalledFuncsFromFuncDeclRecursive(files, info, f.Name(), newFoundedFuncs)
		if err != nil {
			return nil, nil, err
		}
		funcs = append(funcs, newFuncs...)
		funcDecls = append(funcDecls, newFuncDecls...)
	}
	return
}

func isFuncName(funcNames []string, funcName string) bool {
	for _, n := range funcNames {
		if n == funcName {
			return true
		}
	}
	return false
}

// findFuncDeclByName は指定された名前の関数をfilesから検索して返す。なければnil
func findFuncDeclByName(files []*ast.File, name string) *ast.FuncDecl {
	for _, file := range files {
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				if funcDecl.Name.Name == name {
					return funcDecl
				}
			}
		}
	}
	return nil
}

// extractCalledFuncsFromFuncDecl は指定したパッケージの指定したfuncDecl内で呼び出されている関数一覧を返す。
func extractCalledFuncsFromFuncDecl(info *types.Info, targetFuncDecl *ast.FuncDecl) (funcs []*types.Func) {
	ast.Inspect(targetFuncDecl, func(node ast.Node) bool {
		if t, _ := node.(*ast.Ident); t != nil {

			obj := info.ObjectOf(t)
			if tFunc, _ := obj.(*types.Func); tFunc != nil {
				// 自分自身は無視
				if obj.Name() != targetFuncDecl.Name.Name {
					funcs = append(funcs, tFunc)
				}
			}
			return false
		}
		return true
	})
	return funcs
}
