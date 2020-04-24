package ast

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ast/astutil"

	"golang.org/x/tools/go/packages"

	"github.com/mpppk/gollup/util"
)

func NewProgramFromPackages(packageNames []string) (*Packages, error) {
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
	return NewPackages(pkgs), nil
	//m := map[string]*packages.Package{}
	//for _, pkg := range pkgs {
	//	m[pkg.PkgPath] = pkg
	//}
	//return m, nil
}

func NewMergedFileFromPackageInfo(files []*ast.File) *ast.File {
	importDecl := mergeImportDecls(files)
	constDecls := mergeConstDecls(files)

	var imports []*ast.ImportSpec
	for _, file := range files {
		imports = append(imports, file.Imports...)
	}
	return &ast.File{
		Name: &ast.Ident{
			Name: "main",
		},
		Decls:      append([]ast.Decl{importDecl}, GenDeclToDecl(constDecls)...),
		Scope:      nil,
		Imports:    imports,
		Unresolved: nil,
		Comments:   nil,
	}
}

func ExtractObjectsFromFuncDeclRecursive(pkgs map[string]*packages.Package, f *types.Func, objects []types.Object) ([]types.Object, error) {
	//func ExtractObjectsFromFuncDeclRecursive(pkgs map[string]*Packages.Package, packageName, funcName string, objects []types.Object) ([]types.Object, error) {
	pkg := pkgs[f.Pkg().Path()]
	funcDecl := findFuncDeclByName(pkg.Syntax, f.Name())
	if funcDecl == nil {
		return nil, errors.New("specified function is not found: " + f.Name())
	}

	calledFuncs := extractCalledFuncsFromFuncDecl(pkg.TypesInfo, funcDecl)
	objects = append(objects, f)
	//calledFuncs := extractCalledFuncsFromFunc(f)
	for _, f := range calledFuncs {
		if f.Pkg() == nil || util.IsStandardPackage(f.Pkg().Path()) {
			continue
		}

		// 既に発見済みの関数の場合はスキップ
		if _, ok := findObject(objects, f); ok {
			continue
		}

		objs, err := ExtractObjectsFromFuncDeclRecursive(pkgs, f, objects)
		if err != nil {
			return nil, err
		}
		objects = objs
	}
	return objects, nil
}

//func ExtractObjectsFromFuncDeclRecursive(pkgs map[string]*Packages.Package, packageName, funcName string, foundedFuncNames []string) (funcs map[string]map[string]*types.Func, funcDecls map[string][]*ast.FuncDecl, err error) {
//	pkg := pkgs[packageName]
//	funcDecl := findFuncDeclByName(pkg.Syntax, funcName)
//	if funcDecl == nil {
//		return nil, nil, errors.New("specified function is not found: " + funcName)
//	}
//
//	funcDecls = map[string][]*ast.FuncDecl{}
//	funcDecls[packageName] = append(funcDecls[packageName], funcDecl)
//	calledFuncs := extractCalledFuncsFromFuncDecl(pkg.TypesInfo, funcDecl)
//	newFoundedFuncs := make([]string, len(foundedFuncNames))
//	copy(newFoundedFuncs, foundedFuncNames)
//	newFoundedFuncs = append(newFoundedFuncs, funcDecl.Name.Name)
//
//	funcs = map[string]map[string]*types.Func{}
//	for targetPkgName, fs := range calledFuncs {
//		for _, f := range fs {
//			if f.Pkg() == nil || util.IsStandardPackage(f.Pkg().Path()) {
//				continue
//			}
//
//			// 既に発見済みの関数の場合はスキップ
//			if isFuncName(foundedFuncNames, f.Name()) {
//				continue
//			}
//
//			if funcs[targetPkgName] == nil {
//				funcs[targetPkgName] = map[string]*types.Func{}
//			}
//			funcs[targetPkgName][f.Name()] = f
//			newFuncs, newFuncDecls, err := ExtractObjectsFromFuncDeclRecursive(pkgs, targetPkgName, f.Name(), newFoundedFuncs)
//			if err != nil {
//				return nil, nil, err
//			}
//			funcs = mergeFuncMapMap(funcs, newFuncs)
//
//			for pkgName, decls := range newFuncDecls {
//				funcDecls[pkgName] = append(funcDecls[pkgName], decls...)
//			}
//		}
//	}
//	return
//}

//func extractCalledFuncsFromFuncDecl(info *types.Info, targetFuncDecl *ast.FuncDecl) (funcs map[string][]*types.Func) {
//	funcs = map[string][]*types.Func{}
//	ast.Inspect(targetFuncDecl, func(node ast.Node) bool {
//		if callExpr, ok := node.(*ast.CallExpr); ok {
//			if ident, ok := callExpr.Fun.(*ast.Ident); ok {
//				obj := info.ObjectOf(ident)
//				tFunc, ok := obj.(*types.Func)
//				if !ok {
//					return true
//				}
//				if tFunc.Pkg() != nil {
//					funcs[tFunc.Pkg().Name()] = append(funcs[tFunc.Pkg().Name()], tFunc)
//				}
//				return true
//			}
//			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
//				obj := info.ObjectOf(selectorExpr.Sel)
//				// 自分自身は無視
//				if obj.Name() == targetFuncDecl.Name.Name {
//					return true
//				}
//
//				tFunc, ok := obj.(*types.Func)
//				if !ok {
//					return true
//				}
//				if tFunc.Pkg() != nil {
//					funcs[tFunc.Pkg().Name()] = append(funcs[tFunc.Pkg().Name()], tFunc)
//					return true
//				}
//				xident, ok := selectorExpr.X.(*ast.Ident)
//				if !ok {
//					panic("not ident x")
//				}
//				xObj := info.ObjectOf(xident)
//				funcs[xObj.Pkg().Name()] = append(funcs[tFunc.Pkg().Name()], tFunc)
//				return true
//			}
//		}
//		return true
//	})
//	return funcs
//}

func extractCalledFuncsFromFunc(f *types.Func) (funcs []*types.Func) {
	fmt.Println(f.Scope().Names())
	for _, name := range f.Scope().Names() {
		obj := f.Scope().Lookup(name)
		if childFunc, ok := obj.(*types.Func); ok {
			funcs = append(funcs, childFunc)
		}
	}
	return funcs
}

// extractCalledFuncsFromFuncDecl は指定したパッケージの指定したfuncDecl内で呼び出されている関数を、その関数が属するパッケージ名をキーとしたmapとして返す。
func extractCalledFuncsFromFuncDecl(info *types.Info, targetFuncDecl *ast.FuncDecl) (funcs []*types.Func) {
	ast.Inspect(targetFuncDecl, func(node ast.Node) bool {
		if callExpr, ok := node.(*ast.CallExpr); ok {
			if ident, ok := callExpr.Fun.(*ast.Ident); ok {
				obj := info.ObjectOf(ident)
				tFunc, ok := obj.(*types.Func)
				if !ok {
					return true
				}
				if tFunc.Pkg() != nil {
					funcs = append(funcs, tFunc)
				}
				return true
			}
			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				obj := info.ObjectOf(selectorExpr.Sel)
				// 自分自身は無視
				if obj.Name() == targetFuncDecl.Name.Name {
					return true
				}

				tFunc, ok := obj.(*types.Func)
				if !ok {
					return true
				}
				funcs = append(funcs, tFunc)
				return true
			}
		}
		return true
	})
	return funcs
}

func RenameExternalPackageFunctions(pkgs *Packages, sdecls *Decls) {
	//var decls []ast.Decl
	//for _, object := range objects {
	//	decl := pkgs.FindDeclByObject(object)
	//	decls = append(decls, decl)
	//}
	//sdecls := NewDecls(decls)
	for i, funcDecl := range sdecls.Funcs {
		object := sdecls.FuncObjects[i]
		astutil.Apply(funcDecl, func(cursor *astutil.Cursor) bool {
			if callExpr, ok := cursor.Node().(*ast.CallExpr); ok {
				if newCallExpr := removePackageFromCallExpr(callExpr, pkgs.getPkg(object.Pkg().Path())); newCallExpr != nil {
					cursor.Replace(newCallExpr)
				}
			}
			//if compositeLit, ok := cursor.Node().(*ast.CompositeLit); ok {
			//	if newCompositeLit := removePackageFromCompositeLit(compositeLit, funcDeclPkgPath, funcMapMap); newCompositeLit != nil {
			//		cursor.Replace(newCompositeLit)
			//	}
			//}
			return true
		}, nil)

		funcDecl.Name = ast.NewIdent(renameFunc(object.Pkg().Name(), funcDecl.Name.Name))
	}

	//for funcDeclPkgPath, fDecls := range decls {
	//	for _, fDecl := range fDecls {
	//		astutil.Apply(fDecl, func(cursor *astutil.Cursor) bool {
	//			if callExpr, ok := cursor.Node().(*ast.CallExpr); ok {
	//				if newCallExpr := removePackageFromCallExpr(callExpr, pkgs.getPkg(funcDeclPkgPath)); newCallExpr != nil {
	//					cursor.Replace(newCallExpr)
	//				}
	//			}
	//			//if compositeLit, ok := cursor.Node().(*ast.CompositeLit); ok {
	//			//	if newCompositeLit := removePackageFromCompositeLit(compositeLit, funcDeclPkgPath, funcMapMap); newCompositeLit != nil {
	//			//		cursor.Replace(newCompositeLit)
	//			//	}
	//			//}
	//			return true
	//		}, nil)
	//		fDecl.Name.Name = renameFunc(funcDeclPkgPath, fDecl.Name.Name)
	//	}
	//}
}

//func RenameExternalPackageFunctions(decls map[string][]*ast.FuncDecl, funcMapMap map[string]map[string]*types.Func) {
//	for funcDeclPkgName, fDecls := range decls {
//		for _, fDecl := range fDecls {
//			astutil.Apply(fDecl, func(cursor *astutil.Cursor) bool {
//				if callExpr, ok := cursor.Node().(*ast.CallExpr); ok {
//					if newCallExpr := removePackageFromCallExpr(callExpr, funcDeclPkgName, funcMapMap); newCallExpr != nil {
//						cursor.Replace(newCallExpr)
//					}
//				}
//				//if compositeLit, ok := cursor.Node().(*ast.CompositeLit); ok {
//				//	if newCompositeLit := removePackageFromCompositeLit(compositeLit, funcDeclPkgName, funcMapMap); newCompositeLit != nil {
//				//		cursor.Replace(newCompositeLit)
//				//	}
//				//}
//				return true
//			}, nil)
//			fDecl.Name.Name = renameFunc(funcDeclPkgName, fDecl.Name.Name)
//		}
//	}
//}

//func removePackageFromCallExpr(callExpr *ast.CallExpr, currentPkgName string, funcMapMap map[string]map[string]*types.Func) *ast.CallExpr {
//	if ident, ok := callExpr.Fun.(*ast.Ident); ok {
//		f, ok := funcMapMap[currentPkgName][ident.Name]
//		if !ok {
//			return nil
//		}
//		// 置き換え
//		return &ast.CallExpr{
//			Fun: &ast.BasicLit{
//				Kind:  token.STRING,
//				Value: renameFunc(f.Pkg().Name(), ident.Name),
//			},
//			Args: callExpr.Args,
//		}
//	}
//
//	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
//	if !ok {
//		return nil
//	}
//	x, ok := selExpr.X.(*ast.Ident)
//	if !ok {
//		return nil
//	}
//
//	f, ok := funcMapMap[x.Name][selExpr.Sel.Name]
//	if !ok {
//		return nil
//	}
//
//	// 置き換え
//	return &ast.CallExpr{
//		Fun: &ast.BasicLit{
//			Kind:  token.STRING,
//			Value: renameFunc(f.Pkg().Name(), selExpr.Sel.Name),
//		},
//		Args: callExpr.Args,
//	}
//}

// package名の部分を削除したCallExprを返します(非破壊). 存在しない名前の関数である場合や想定しない構造の場合はnilを返します.
func removePackageFromCallExpr(callExpr *ast.CallExpr, pkg *packages.Package) *ast.CallExpr {
	//pkg := pkgs.getPkg(pkgPath)
	if ident, ok := callExpr.Fun.(*ast.Ident); ok {
		obj := pkg.TypesInfo.ObjectOf(ident)
		//f, ok := getFuncFromIdent(pkg, ident)
		if !ok {
			return nil
		}
		// 置き換え
		return &ast.CallExpr{
			Fun: &ast.BasicLit{
				Kind:  token.STRING,
				Value: renameFunc(obj.Pkg().Name(), ident.Name),
			},
			Args: callExpr.Args,
		}
	}

	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	obj := pkg.TypesInfo.ObjectOf(selExpr.Sel)
	//f, ok := pkgs.getFuncFromSelectorExpr(pkgPath, selExpr)
	//if !ok {
	//	return nil
	//}
	//x, ok := selExpr.X.(*ast.Ident)
	//if !ok {
	//	return nil
	//}
	//
	//xObj := pkg.TypesInfo.ObjectOf(x)
	//if xObj == nil {
	//	return nil
	//}
	//pkg2 := pkgs.getPkg(xObj.Pkg().Path())
	//
	//f, ok := getFuncFromIdent(pkg2, selExpr.Sel)
	//if !ok {
	//	return nil
	//}

	// 置き換え
	return &ast.CallExpr{
		Fun: &ast.BasicLit{
			Kind:  token.STRING,
			Value: renameFunc(obj.Pkg().Name(), selExpr.Sel.Name),
		},
		Args: callExpr.Args,
	}
}

// package名の部分を削除したCompositeLitを返します(非破壊). 存在しない名前の関数である場合や想定しない構造の場合はnilを返します.
//func removePackageFromCompositeLit(compositeLit *ast.CompositeLit, currentPkgName string, files []*ast.File) *ast.CompositeLit {
//	if ident, ok := compositeLit.Type.(*ast.Ident); ok {
//		return compositeLit
//		return &ast.CompositeLit{
//			Type: ast.NewIdent(ident.Name),
//		}
//	}
//
//	selExpr, ok := compositeLit.Type.(*ast.SelectorExpr)
//	if !ok {
//		return nil
//	}
//	x, ok := selExpr.X.(*ast.Ident)
//	if !ok {
//		return nil
//	}
//
//	f, ok := funcMapMap[x.Name][selExpr.Sel.Name]
//	if !ok {
//		return nil
//	}
//
//	// 置き換え
//	return &ast.CompositeLit{
//		Type: ast.NewIdent(renameFunc(f.Pkg().Name(), selExpr.Sel.Name)),
//	}
//}

//func findPackageOfGenDecl(declName string, files []*ast.File) string {
//
//}
