package ast

import (
	"go/ast"
	"go/token"
	"go/types"
	"sort"

	"golang.org/x/tools/go/packages"

	"github.com/go-toolsmith/astcopy"
)

type Decls struct {
	Objects       []types.Object
	Decls         []ast.Decl
	ImportObjects []types.Object
	Imports       []*ast.GenDecl
	ConstObjects  []types.Object
	Consts        []*ast.GenDecl
	TypeObjects   []types.Object
	Types         []*ast.GenDecl
	VarObjects    []types.Object
	Vars          []*ast.GenDecl
	FuncObjects   []types.Object
	Funcs         []*ast.FuncDecl
}

func NewDecls(pkgs *Packages, objects []types.Object) *Decls {
	var decls []ast.Decl
	for _, object := range objects {
		decl := pkgs.FindDeclByObject(object)
		decls = append(decls, decl)
	}
	sdecls := &Decls{Decls: decls, Objects: objects}
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
		sortSpecs(constDecl.Specs)
		sdecls.Consts = append(sdecls.Consts, constDecl)
	}
	return sdecls
}

func (d *Decls) findDeclFromObject(object types.Object) (ast.Decl, bool) {
	for i, o := range d.Objects {
		if o.Id() == object.Id() {
			return d.Decls[i], true
		}
	}
	return nil, false
}

func getFuncFromIdent(pkg *packages.Package, ident *ast.Ident) (*types.Func, bool) {
	obj := pkg.TypesInfo.ObjectOf(ident)
	if obj == nil {
		return nil, false
	}
	f, ok := obj.(*types.Func)
	return f, ok
}

func mergeImportDecls(files []*ast.File) (importDecl *ast.GenDecl) {
	for _, file := range files {
		imports := selectGenDeclsFromDecls(file.Decls, token.IMPORT)
		if importDecl == nil {
			importDecl = imports[0]
		}
		for _, genDecl := range imports {
			importDecl.Specs = append(importDecl.Specs, genDecl.Specs...)
		}
	}
	return
}

func selectGenDeclsFromDecls(decls []ast.Decl, tkn token.Token) (importDecls []*ast.GenDecl) {
	for _, decl := range decls {
		if importDecl, ok := declToGenDecl(decl, tkn); ok {
			importDecls = append(importDecls, importDecl)
		}
	}
	return
}

func declToGenDecl(decl ast.Decl, tkn token.Token) (*ast.GenDecl, bool) {
	if genDecl, ok := decl.(*ast.GenDecl); ok {
		if genDecl.Tok == tkn {
			return genDecl, true
		}
	}
	return nil, false
}

func GenDeclToDecl(genDecls []*ast.GenDecl) (decls []ast.Decl) {
	for _, decl := range genDecls {
		decls = append(decls, decl)
	}
	return
}

func declToFuncDecl(decls []ast.Decl) (funcDecls []*ast.FuncDecl) {
	for _, decl := range decls {
		if fd, ok := decl.(*ast.FuncDecl); ok {
			funcDecls = append(funcDecls, fd)
		}
	}
	return
}

func funcDeclToDecl(funcDecls []*ast.FuncDecl) (decls []ast.Decl) {
	for _, fd := range funcDecls {
		decls = append(decls, fd)
	}
	return
}

func CopyFuncDeclsAsDecl(funcDecls []*ast.FuncDecl) (newFuncDecls []ast.Decl) {
	for _, decl := range funcDecls {
		newFuncDecls = append(newFuncDecls, astcopy.FuncDecl(decl))
	}
	return
}

func renameFunc(pkg *types.Package, funcName string) string {
	if pkg == nil || pkg.Name() == "main" {
		return funcName
	}
	return pkg.Name() + "_" + funcName
}

func SortGenDecls(genDecls []*ast.GenDecl) {
	sort.Slice(genDecls, func(i, j int) bool {
		spec1 := genDecls[i].Specs[0]
		spec2 := genDecls[j].Specs[0]
		return specToString(spec1) < specToString(spec2)
	})
}

func sortSpecs(specs []ast.Spec) {
	// FIXME: sort ValueSpecs names
	sort.Slice(specs, func(i, j int) bool {
		return specToString(specs[i]) < specToString(specs[j])
	})
}

func specToString(spec ast.Spec) string {
	switch s := spec.(type) {
	case *ast.ImportSpec:
		return s.Name.Name
	case *ast.ValueSpec:
		return s.Names[0].Name
	case *ast.TypeSpec:
		return s.Name.Name
	}
	return ""
}

func SortFuncDeclsFromDecls(decls []ast.Decl) []ast.Decl {
	funcDecls := declToFuncDecl(decls)
	getRecvName := func(funcDecl *ast.FuncDecl) (recv string) {
		if funcDecl.Recv == nil {
			return
		}
		for _, field := range funcDecl.Recv.List {
			for _, name := range field.Names {
				recv += name.Name
			}
		}
		return
	}
	sort.Slice(funcDecls, func(i, j int) bool {
		return getRecvName(funcDecls[i])+funcDecls[i].Name.Name < getRecvName(funcDecls[j])+funcDecls[j].Name.Name
	})
	return funcDeclToDecl(funcDecls)
}

// findFuncDeclByFuncType は指定された名前の関数をfilesから検索して返す。なければnil
func findFuncDeclByFuncType(files []*ast.File, f *types.Func) (funcDecl *ast.FuncDecl) {
	sig := f.Type().(*types.Signature)
	for _, file := range files {
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				if funcDecl.Name.Name != f.Name() {
					continue
				}

				if funcDecl.Recv == nil {
					if sig.Recv() != nil {
						continue
					} else {
						return funcDecl
					}
				} else {
					t := funcDecl.Recv.List[0].Type
					if sig.Recv() == nil {
						continue
					}
					recv := sig.Recv()
					switch expr := t.(type) {
					case *ast.StarExpr:
						if ident, ok := expr.X.(*ast.Ident); ok {
							if ptr, ok := recv.Type().(*types.Pointer); ok {
								named := ptr.Elem().(*types.Named)
								name := named.Obj().Name()
								if ident.Name == name {
									return funcDecl
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func RemoveCommentsFromFuncDecls(funcDecls []*ast.FuncDecl) {
	for _, decl := range funcDecls {
		decl.Doc = nil
	}
}
