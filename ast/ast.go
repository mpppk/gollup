package ast

import (
	"errors"
	"go/ast"
	"go/token"
	"go/types"
	"log"

	"golang.org/x/tools/go/ast/astutil"

	"golang.org/x/tools/go/packages"

	"github.com/mpppk/gollup/util"
)

func NewProgramFromPackages(packageNames []string) (*Packages, *token.FileSet, error) {
	fset := token.NewFileSet()
	config := &packages.Config{
		Mode: packages.NeedCompiledGoFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.LoadAllSyntax,
		Fset: fset,
	}
	pkgs, err := packages.Load(config, packageNames...)
	if err != nil {
		return nil, nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, nil, errors.New("error occurred in NewProgramFromPackages")
	}
	return NewPackages(pkgs), fset, nil
}

func NewMergedFileFromPackageInfo(files []*ast.File) *ast.File {
	importDecl := mergeImportDecls(files)

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

func ExtractObjectsFromFuncDeclRecursive(pkgs map[string]*packages.Package, f *types.Func, objects []types.Object) ([]types.Object, error) {
	log.Println("searching objects from func", f.Pkg().Name()+"."+f.Name())
	pkg := pkgs[f.Pkg().Path()]
	if pkg == nil {
		return nil, errors.New("specified function is not found in pkgs: " + f.Name())
	}
	funcDecl := findFuncDeclByFuncType(pkg.Syntax, f)
	if funcDecl == nil {
		return nil, errors.New("specified function is not found: " + f.Name())
	}

	calledFuncs := extractCalledFuncsFromFuncDecl(pkg.TypesInfo, funcDecl)
	newObjects := extractNonStandardObjectFromFuncDecl(pkg.TypesInfo, funcDecl)
	objects = append(objects, newObjects...)
	objects = append(objects, f)
	for _, f := range calledFuncs {
		if !util.HasPkg(f) || util.IsStandardPackage(f.Pkg().Path()) {
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
	objects = distinctObjects(objects)
	return objects, nil
}

// extractCalledFuncsFromFuncDecl は指定したパッケージの指定したfuncDecl内で呼び出されている関数を返す
func extractCalledFuncsFromFuncDecl(info *types.Info, targetFuncDecl *ast.FuncDecl) (funcs []*types.Func) {
	ast.Inspect(targetFuncDecl, func(node ast.Node) bool {
		if callExpr, ok := node.(*ast.CallExpr); ok {
			if f := callExprToFunc(info, callExpr); f != nil && f.Name() != targetFuncDecl.Name.Name {
				funcs = append(funcs, f)
			}
		}
		return true
	})
	return funcs
}

// extractStructFromFuncDecl は指定したパッケージの指定したfuncDecl内で呼び出されている関数から参照されている型名を返す
func extractNonStandardObjectFromFuncDecl(info *types.Info, targetFuncDecl *ast.FuncDecl) (objects []types.Object) {
	ast.Inspect(targetFuncDecl, func(node ast.Node) bool {
		ident, ok := node.(*ast.Ident)
		if !ok {
			return true
		}
		obj := info.ObjectOf(ident)
		if !util.HasPkg(obj) || util.IsStandardPackage(obj.Pkg().Path()) {
			return true
		}
		switch obj.(type) {
		case *types.Const, *types.Var, *types.TypeName:
			if obj.Pkg().Name() == "main" || obj.Exported() {
				objects = append(objects, obj)
			}
		}
		return true
	})
	return
}

func callExprToFunc(info *types.Info, callExpr *ast.CallExpr) *types.Func {
	switch fun := callExpr.Fun.(type) {
	case *ast.Ident:
		obj := info.ObjectOf(fun)
		tFunc, ok := obj.(*types.Func)
		if !ok {
			return nil
		}
		if tFunc.Pkg() != nil {
			return tFunc
		}
		return nil
	case *ast.SelectorExpr:
		obj := info.ObjectOf(fun.Sel)
		tFunc, ok := obj.(*types.Func)
		if !ok {
			return nil
		}
		return tFunc
	}
	return nil
}

func RenameExternalPackageFunctions(pkgs *Packages, sdecls *Decls) {
	for i, funcDecl := range sdecls.Funcs {
		object := sdecls.FuncObjects[i]
		pkg := pkgs.getPkg(object.Pkg().Path())
		renameExternalPackageFunction(funcDecl, object, pkg)
	}
}

func renameExternalPackageFunction(funcDecl *ast.FuncDecl, object types.Object, pkg *packages.Package) {
	astutil.Apply(funcDecl, func(cursor *astutil.Cursor) bool {
		if callExpr, ok := cursor.Node().(*ast.CallExpr); ok {
			if newCallExpr := removePackageFromCallExpr(callExpr, pkg); newCallExpr != nil {
				cursor.Replace(newCallExpr)
			}
		}
		if compositeLit, ok := cursor.Node().(*ast.CompositeLit); ok {
			if newCompositeLit := removePackageFromCompositeLit(compositeLit, pkg); newCompositeLit != nil {
				cursor.Replace(newCompositeLit)
			}
		}
		return true
	}, nil)

	// 構造体のメソッドはrenameしない
	if funcDecl.Recv == nil {
		funcDecl.Name = ast.NewIdent(renameFunc(object.Pkg(), funcDecl.Name.Name))
	}

	renameFuncDeclParams(funcDecl, pkg)
	renameFuncDeclResults(funcDecl, pkg)
}

// 他ライブラリの構造体などを引数に取っていればrename
func renameFuncDeclParams(funcDecl *ast.FuncDecl, pkg *packages.Package) {
	if funcDecl.Type.Params == nil {
		return
	}
	for i, result := range funcDecl.Type.Params.List {
		switch t := result.Type.(type) {
		case *ast.SelectorExpr:
			obj := pkg.TypesInfo.ObjectOf(t.Sel)
			if !util.HasPkg(obj) || util.IsStandardPackage(obj.Pkg().Path()) {
				continue
			}
			funcDecl.Type.Params.List[i].Type = ast.NewIdent(t.Sel.Name)
		case *ast.StarExpr:
			selExpr, ok := t.X.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			obj := pkg.TypesInfo.ObjectOf(selExpr.Sel)
			if !util.HasPkg(obj) || util.IsStandardPackage(obj.Pkg().Path()) {
				continue
			}
			funcDecl.Type.Params.List[i].Type = &ast.StarExpr{
				X: ast.NewIdent(selExpr.Sel.Name),
			}
		}
	}
}

// 他ライブラリの構造体などが戻り値であればrename
func renameFuncDeclResults(funcDecl *ast.FuncDecl, pkg *packages.Package) {
	if funcDecl.Type.Results == nil {
		return
	}
	for i, result := range funcDecl.Type.Results.List {
		switch t := result.Type.(type) {
		case *ast.SelectorExpr:
			ident, ok := result.Type.(*ast.Ident)
			if !ok {
				continue
			}
			obj := pkg.TypesInfo.ObjectOf(ident)
			if !util.HasPkg(obj) || util.IsStandardPackage(obj.Pkg().Path()) {
				continue
			}
			funcDecl.Type.Results.List[i].Type = ast.NewIdent(t.Sel.Name)
		case *ast.StarExpr:
			selExpr, ok := t.X.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			ident, ok := result.Type.(*ast.Ident)
			if !ok {
				continue
			}
			obj := pkg.TypesInfo.ObjectOf(ident)
			if !util.HasPkg(obj) || util.IsStandardPackage(obj.Pkg().Path()) {
				continue
			}
			funcDecl.Type.Results.List[i].Type = &ast.StarExpr{
				X: ast.NewIdent(selExpr.Sel.Name),
			}
		}
	}
}

// package名の部分を削除したCallExprを返します(非破壊). 存在しない名前の関数である場合や想定しない構造の場合はnilを返します.
// 標準パッケージの呼び出しである場合は書き換えを行いません。
func removePackageFromCallExpr(callExpr *ast.CallExpr, pkg *packages.Package) *ast.CallExpr {
	if ident, ok := callExpr.Fun.(*ast.Ident); ok {
		obj := pkg.TypesInfo.ObjectOf(ident)
		if !util.HasPkg(obj) || util.IsStandardPackage(obj.Pkg().Path()) {
			return callExpr
		}

		// 置き換え
		return &ast.CallExpr{
			Fun: &ast.BasicLit{
				Kind:  token.STRING,
				Value: renameFunc(obj.Pkg(), ident.Name),
			},
			Ellipsis: callExpr.Ellipsis,
			Args:     callExpr.Args,
		}
	}

	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	obj := pkg.TypesInfo.ObjectOf(selExpr.Sel)
	if !util.HasPkg(obj) || util.IsStandardPackage(obj.Pkg().Path()) {
		return callExpr
	}

	// 構造体のメソッドを呼び出している場合は書き換えない
	if xident, ok := selExpr.X.(*ast.Ident); ok {
		xobj := pkg.TypesInfo.ObjectOf(xident)
		if _, ok := xobj.(*types.Var); ok {
			return callExpr
		}
	}

	// 置き換え
	return &ast.CallExpr{
		Fun: &ast.BasicLit{
			Kind:  token.STRING,
			Value: renameFunc(obj.Pkg(), selExpr.Sel.Name),
		},
		Args: callExpr.Args,
	}
}

// package名の部分を削除したCompositeLitを返します(非破壊). 存在しない名前の関数である場合や想定しない構造の場合はnilを返します.
func removePackageFromCompositeLit(compositeLit *ast.CompositeLit, pkg *packages.Package) *ast.CompositeLit {
	if _, ok := compositeLit.Type.(*ast.Ident); ok {
		return compositeLit
		//return &ast.CompositeLit{
		//	Type: ast.NewIdent(ident.Name),
		//}
	}

	selExpr, ok := compositeLit.Type.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	x, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return nil
	}
	obj := pkg.TypesInfo.ObjectOf(x)

	// 置き換え
	return &ast.CompositeLit{
		Type: ast.NewIdent(renameFunc(obj.Pkg(), selExpr.Sel.Name)),
	}
}
