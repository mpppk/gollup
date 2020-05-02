package ast

import (
	"errors"
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
	pkg := pkgs[f.Pkg().Path()]
	if pkg == nil {
		return nil, errors.New("specified function is not found in pkgs: " + f.Name())
	}
	funcDecl := findFuncDeclByName(pkg.Syntax, f.Name())
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
			objects = append(objects, obj)
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
		astutil.Apply(funcDecl, func(cursor *astutil.Cursor) bool {
			pkg := pkgs.getPkg(object.Pkg().Path())
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

		// 他ライブラリの構造体などを引数に取っていればrename
		if funcDecl.Type.Params != nil {
			for i, param := range funcDecl.Type.Params.List {
				if name, isPtr, ok := getFieldName(param); ok {
					if isPtr {
						funcDecl.Type.Params.List[i].Type = &ast.StarExpr{
							X: ast.NewIdent(name),
						}
					} else {
						funcDecl.Type.Params.List[i].Type = ast.NewIdent(name)
					}
				}
			}
		}

		// 他ライブラリの構造体などが戻り値であればrename
		if funcDecl.Type.Results != nil {
			for i, result := range funcDecl.Type.Results.List {
				if name, isPtr, ok := getFieldName(result); ok {
					if isPtr {
						funcDecl.Type.Results.List[i].Type = &ast.StarExpr{
							X: ast.NewIdent(name),
						}
					} else {
						funcDecl.Type.Results.List[i].Type = ast.NewIdent(name)
					}
				}
			}
		}
	}
}

func getFieldName(field *ast.Field) (string, bool, bool) {
	switch t := field.Type.(type) {
	case *ast.SelectorExpr:
		return t.Sel.Name, false, true
	case *ast.StarExpr:
		if selExpr, ok := t.X.(*ast.SelectorExpr); ok {
			return selExpr.Sel.Name, true, true
		}
	}
	return "", false, false
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
			Args: callExpr.Args,
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
