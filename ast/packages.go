package ast

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Packages struct {
	Packages map[string]*packages.Package
}

func NewPackages(pkgs []*packages.Package) *Packages {
	m := map[string]*packages.Package{}
	for _, pkg := range pkgs {
		m[pkg.PkgPath] = pkg
	}
	return &Packages{
		Packages: m,
	}
}

func (p *Packages) FindPkgByName(name string) (*packages.Package, bool) {
	for _, pkg := range p.Packages {
		if pkg.Name == name {
			return pkg, true
		}
	}
	return nil, false
}

func (p *Packages) ObjectsToDecls(objects []types.Object) []ast.Decl {
	var decls []ast.Decl
	for _, object := range objects {
		decl := p.FindDeclByObject(object)
		decls = append(decls, decl)
	}
	return decls
}

func (p *Packages) getFuncFromSelectorExpr(currentPkgPath string, selExpr *ast.SelectorExpr) (*types.Func, bool) {
	pkg := p.getPkg(currentPkgPath)

	x, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return nil, false
	}

	xObj := pkg.TypesInfo.ObjectOf(x)
	if xObj == nil {
		return nil, false
	}
	pkg2 := p.getPkg(xObj.Pkg().Path())

	f, ok := getFuncFromIdent(pkg2, selExpr.Sel)
	if !ok {
		return nil, false
	}
	return f, ok
}

func (p *Packages) getPkg(path string) *packages.Package {
	return p.Packages[path]
}

func (p *Packages) FindDeclByObject(object types.Object) ast.Decl {
	pkg := p.getPkg(object.Pkg().Path())
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				if d.Tok == token.IMPORT {
					continue
				}
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.ImportSpec:
						continue
					case *ast.ValueSpec:
						for _, name := range s.Names {
							if name.Name == object.Name() {
								return decl
							}
						}
					case *ast.TypeSpec:
						if s.Name.Name == object.Name() {
							return decl
						}
					}
				}
			case *ast.FuncDecl:
				if d.Name.Name == object.Name() {
					return decl
				}
			}
		}
	}
	return nil
}

func (p *Packages) findObject(pkgPath string, ident *ast.Ident) types.Object {
	pkg := p.Packages[pkgPath]
	return pkg.TypesInfo.ObjectOf(ident)
}

// 与えられたidentが参照しているDeclと所属しているパッケージ名を返します。
func (p *Packages) findDeclOfIdent(pkgPath string, ident *ast.Ident) (string, ast.Decl) {
	obj := p.findObject(pkgPath, ident)
	targetPkgName := obj.Pkg().Name()
	pkg := p.Packages[targetPkgName]
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				if d.Tok == token.IMPORT {
					continue
				}
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.ImportSpec:
						continue
					case *ast.ValueSpec:
						for _, name := range s.Names {
							if name.Name == ident.Name {
								return obj.Pkg().Name(), decl
							}
						}
					case *ast.TypeSpec:
						if s.Name.Name == ident.Name {
							return targetPkgName, decl
						}
					}
				}
			case *ast.FuncDecl:
				if d.Name.Name == ident.Name {
					return targetPkgName, decl
				}
			}
		}
	}
	return "", nil
}
