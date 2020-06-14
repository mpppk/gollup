package ast

import (
	"go/ast"
	"go/token"
	"go/types"
	"sort"

	"golang.org/x/tools/go/packages"

	"github.com/go-toolsmith/astcopy"
)

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
			expr := unwrapStarExpr(field.Type)
			if ident, ok := expr.(*ast.Ident); ok {
				recv += ident.Name
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
						if expr.X.(*ast.Ident).Name == getRecvTypeName(recv) {
							return funcDecl
						}
					case *ast.Ident:
						if expr.Name == getRecvTypeName(recv) {
							return funcDecl
						}
					}
				}
			}
		}
	}
	return nil
}

func getRecvTypeName(recv *types.Var) string {
	switch recvType := recv.Type().(type) {
	case *types.Pointer:
		recvNewType := unwrapPointer(recvType)
		switch t := recvNewType.(type) {
		case *types.Named:
			return t.Obj().Name()
		default:
			panic("unknown type in getRecvTypeName: " + t.String())
		}
	case *types.Named:
		return recvType.Obj().Name()
	default:
		panic("unknown recv type in getRecvTypeName: " + recvType.String())
	}
}

func unwrapStarExpr(expr ast.Expr) ast.Expr {
	cnt := 0
	for {
		switch newExpr := expr.(type) {
		case *ast.StarExpr:
			expr = newExpr.X
			cnt++
			if cnt > 100 {
				panic(expr)
			}
		default:
			return newExpr
		}
	}
}

func unwrapPointer(ptr *types.Pointer) (retType types.Type) {
	for {
		if p, ok := ptr.Elem().(*types.Pointer); ok {
			ptr = p
		} else {
			retType = ptr.Elem()
			break
		}
	}
	return
}

func removeCommentsFromFuncDecls(funcDecls []*ast.FuncDecl) {
	for _, decl := range funcDecls {
		decl.Doc = nil
	}
}
