package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"

	"github.com/pkg/errors"
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
