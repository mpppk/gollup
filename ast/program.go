package ast

import (
	"go/ast"
	"go/token"
	"go/types"
)

type Program struct {
	Packages      *Packages
	Objects       []types.Object
	Decls         []ast.Decl
	ImportObjects []types.Object
	Imports       []*ast.GenDecl
	ConstObjects  []types.Object
	Const         *ast.GenDecl
	TypeObjects   []types.Object
	Types         []*ast.GenDecl
	VarObjects    []types.Object
	Vars          []*ast.GenDecl
	FuncObjects   []types.Object
	Funcs         []*ast.FuncDecl
}

func NewProgram(pkgs *Packages, objects []types.Object) *Program {
	var decls []ast.Decl
	for _, object := range objects {
		decl := pkgs.FindDeclByObject(object)
		decls = append(decls, decl)
	}
	sdecls := &Program{Packages: pkgs, Decls: decls, Objects: objects}
	constDecl := &ast.GenDecl{Tok: token.CONST}
	for i, decl := range decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
				sdecls.Imports = append(sdecls.Imports, d)
				sdecls.ImportObjects = append(sdecls.ImportObjects, objects[i])
			case token.CONST:
				o := objects[i]
				for _, spec := range d.Specs {
					vspec := spec.(*ast.ValueSpec)
					for i, name := range vspec.Names {
						if o.Name() == name.Name {
							constDecl.Specs = append(constDecl.Specs, &ast.ValueSpec{
								Names:  []*ast.Ident{ast.NewIdent(name.Name)},
								Values: []ast.Expr{vspec.Values[i]},
							})
						}
					}
				}
				sdecls.ConstObjects = append(sdecls.ConstObjects, objects[i])
			case token.TYPE:
				sdecls.Types = append(sdecls.Types, d)
				sdecls.TypeObjects = append(sdecls.TypeObjects, objects[i])
			case token.VAR:
				sdecls.Vars = append(sdecls.Vars, d)
				sdecls.VarObjects = append(sdecls.VarObjects, objects[i])
			}
		case *ast.FuncDecl:
			sdecls.Funcs = append(sdecls.Funcs, d)
			sdecls.FuncObjects = append(sdecls.FuncObjects, objects[i])
		}
	}
	if len(constDecl.Specs) > 0 {
		sdecls.Const = constDecl
	}
	return sdecls
}

func (p *Program) Bundle(files []*ast.File) *ast.File {
	// rename functions
	p.renameExternalPackageFunctions()
	removeCommentsFromFuncDecls(p.Funcs)
	renamedFuncDecls := CopyFuncDeclsAsDecl(p.Funcs)
	renamedFuncDecls = SortFuncDeclsFromDecls(renamedFuncDecls)

	p.addPackagePrefixToConst()

	// rename consts
	if p.Const != nil {
		sortSpecs(p.Const.Specs)
	}
	SortGenDecls(p.Vars)
	SortGenDecls(p.Types)

	file := newMergedFileFromPackageInfo(files)
	if p.Const != nil {
		file.Decls = append(file.Decls, p.Const)
	}
	file.Decls = append(file.Decls, GenDeclToDecl(p.Vars)...)
	file.Decls = append(file.Decls, GenDeclToDecl(p.Types)...)
	file.Decls = append(file.Decls, renamedFuncDecls...)
	return file
}

func (p *Program) findDeclFromObject(object types.Object) (ast.Decl, bool) {
	for i, o := range p.Objects {
		if o.Id() == object.Id() {
			return p.Decls[i], true
		}
	}
	return nil, false
}

func (p *Program) renameExternalPackageFunctions() {
	for i, funcDecl := range p.Funcs {
		object := p.FuncObjects[i]
		pkg := p.Packages.getPkg(object.Pkg().Path())
		renameExternalPackageFunction(funcDecl, object, pkg)
		renameExternalPackageConst(funcDecl, pkg)
	}
}

func (p *Program) addPackagePrefixToConst() {
	if p.Const == nil {
		return
	}
	for i, spec := range p.Const.Specs {
		pkgName := p.ConstObjects[i].Pkg().Name()
		if pkgName != "main" { // FIXME
			addPrefixToSpec(spec, pkgName+"_")
		}
	}
}
