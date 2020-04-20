package ast

import (
	"errors"
	"go/ast"
	"go/token"
	"go/types"
	"reflect"
	"sort"

	"golang.org/x/tools/go/ast/astutil"

	"golang.org/x/tools/go/packages"

	"github.com/mpppk/gollup/util"
)

type Packages struct {
	M map[string]*packages.Package
}

func NewProgramFromPackages(packageNames []string) (map[string]*packages.Package, error) {
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
	m := map[string]*packages.Package{}
	for _, pkg := range pkgs {
		m[pkg.Name] = pkg
	}
	return m, nil
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

func ExtractCalledFuncsFromFuncDeclRecursive(pkgs map[string]*packages.Package, packageName, funcName string, foundedFuncNames []string) (funcs map[string]map[string]*types.Func, funcDecls map[string][]*ast.FuncDecl, err error) {
	pkg := pkgs[packageName]
	funcDecl := findFuncDeclByName(pkg.Syntax, funcName)
	if funcDecl == nil {
		return nil, nil, errors.New("specified function is not found: " + funcName)
	}

	funcDecls = map[string][]*ast.FuncDecl{}
	funcDecls[packageName] = append(funcDecls[packageName], funcDecl)
	calledFuncs := extractCalledFuncsFromFuncDecl(pkg.TypesInfo, funcDecl)
	newFoundedFuncs := make([]string, len(foundedFuncNames))
	copy(newFoundedFuncs, foundedFuncNames)
	newFoundedFuncs = append(newFoundedFuncs, funcDecl.Name.Name)

	funcs = map[string]map[string]*types.Func{}
	for targetPkgName, fs := range calledFuncs {
		for _, f := range fs {
			if f.Pkg() == nil || util.IsStandardPackage(f.Pkg().Path()) {
				continue
			}

			// 既に発見済みの関数の場合はスキップ
			if isFuncName(foundedFuncNames, f.Name()) {
				continue
			}

			if funcs[targetPkgName] == nil {
				funcs[targetPkgName] = map[string]*types.Func{}
			}
			funcs[targetPkgName][f.Name()] = f
			newFuncs, newFuncDecls, err := ExtractCalledFuncsFromFuncDeclRecursive(pkgs, targetPkgName, f.Name(), newFoundedFuncs)
			if err != nil {
				return nil, nil, err
			}
			funcs = mergeFuncMapMap(funcs, newFuncs)

			for pkgName, decls := range newFuncDecls {
				funcDecls[pkgName] = append(funcDecls[pkgName], decls...)
			}
		}
	}
	return
}

func ListUsedStructNames(funcDecls map[string][]*ast.FuncDecl) map[string][]string {
	m := map[string][]string{}
	for pkg, decls := range funcDecls {
		for _, decl := range decls {
			if reflect.ValueOf(decl.Recv).IsNil() {
				continue
			}
			recv := decl.Recv.List[0]
			if starExpr, ok := recv.Type.(*ast.StarExpr); ok {
				if xident, ok := starExpr.X.(*ast.Ident); ok {
					m[pkg] = append(m[pkg], xident.Name)
				}
			}
		}
	}
	return m
}

func mergeFuncMapMap(m1, m2 map[string]map[string]*types.Func) map[string]map[string]*types.Func {
	nm := copyFuncMapMap(m1)
	for pkgName, cm := range m2 {
		nm[pkgName] = mergeFuncMap(nm[pkgName], cm)
	}
	return nm
}

func mergeFuncMap(m1, m2 map[string]*types.Func) map[string]*types.Func {
	cpm1 := copyFuncMap(m1)
	for funcName, f := range m2 {
		cpm1[funcName] = f
	}
	return cpm1
}

func copyFuncMapMap(m map[string]map[string]*types.Func) map[string]map[string]*types.Func {
	nm := map[string]map[string]*types.Func{}
	for pkgName, m2 := range m {
		nm[pkgName] = copyFuncMap(m2)
	}
	return nm
}

func copyFuncMap(m map[string]*types.Func) map[string]*types.Func {
	nm := map[string]*types.Func{}
	for funcName, f := range m {
		nm[funcName] = f
	}
	return nm
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

// FindTypeGenDeclByName は指定された名前の関数をfilesから検索して返す。なければnil
func FindTypeGenDeclByName(files []*ast.File, name string) *ast.GenDecl {
	for _, file := range files {
		for _, decl := range file.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok {
				if genDecl.Tok != token.TYPE {
					continue
				}
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if typeSpec.Name.Name == name {
							return genDecl
						}
					}
				}
			}
		}
	}
	return nil
}

// extractCalledFuncsFromFuncDecl は指定したパッケージの指定したfuncDecl内で呼び出されている関数を、その関数が属するパッケージ名をキーとしたmapとして返す。
func extractCalledFuncsFromFuncDecl(info *types.Info, targetFuncDecl *ast.FuncDecl) (funcs map[string][]*types.Func) {
	funcs = map[string][]*types.Func{}
	ast.Inspect(targetFuncDecl, func(node ast.Node) bool {
		if callExpr, ok := node.(*ast.CallExpr); ok {
			if ident, ok := callExpr.Fun.(*ast.Ident); ok {
				obj := info.ObjectOf(ident)
				tFunc, ok := obj.(*types.Func)
				if !ok {
					return true
				}
				if tFunc.Pkg() != nil {
					funcs[tFunc.Pkg().Name()] = append(funcs[tFunc.Pkg().Name()], tFunc)
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

				if tFunc.Pkg() != nil {
					funcs[tFunc.Pkg().Name()] = append(funcs[tFunc.Pkg().Name()], tFunc)
					return true
				}
				xident, ok := selectorExpr.X.(*ast.Ident)
				if !ok {
					panic("not ident x")
				}
				xObj := info.ObjectOf(xident)
				funcs[xObj.Pkg().Name()] = append(funcs[tFunc.Pkg().Name()], tFunc)
				return true
			}
		}
		return true
	})
	return funcs
}

func RenameExternalPackageFunctions(decls map[string][]*ast.FuncDecl, funcMapMap map[string]map[string]*types.Func) {
	for funcDeclPkgName, fDecls := range decls {
		for _, fDecl := range fDecls {
			astutil.Apply(fDecl, func(cursor *astutil.Cursor) bool {
				if callExpr, ok := cursor.Node().(*ast.CallExpr); ok {
					if newCallExpr := removePackageFromCallExpr(callExpr, funcDeclPkgName, funcMapMap); newCallExpr != nil {
						cursor.Replace(newCallExpr)
					}
				}
				return true
			}, nil)
			fDecl.Name.Name = renameFunc(funcDeclPkgName, fDecl.Name.Name)
		}
	}
}

// package名の部分を削除したCallExprを返します(非破壊). 存在しない名前の関数である場合や想定しない構造の場合はnilを返します.
func removePackageFromCallExpr(callExpr *ast.CallExpr, currentPkgName string, funcMapMap map[string]map[string]*types.Func) *ast.CallExpr {
	if ident, ok := callExpr.Fun.(*ast.Ident); ok {
		f, ok := funcMapMap[currentPkgName][ident.Name]
		if !ok {
			return nil
		}
		// 置き換え
		return &ast.CallExpr{
			Fun: &ast.BasicLit{
				Kind:  token.STRING,
				Value: renameFunc(f.Pkg().Name(), ident.Name),
			},
			Args: callExpr.Args,
		}
	}

	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	x, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return nil
	}

	f, ok := funcMapMap[x.Name][selExpr.Sel.Name]
	if !ok {
		return nil
	}

	// 置き換え
	return &ast.CallExpr{
		Fun: &ast.BasicLit{
			Kind:  token.STRING,
			Value: renameFunc(f.Pkg().Name(), selExpr.Sel.Name),
		},
		Args: callExpr.Args,
	}
}

func renameFunc(pkgName, funcName string) string {
	if pkgName == "main" {
		return funcName
	}
	return pkgName + "_" + funcName
}

func SortFuncDeclsFromDecls(decls []ast.Decl) []ast.Decl {
	funcDecls := declToFuncDecl(decls)
	sort.Slice(funcDecls, func(i, j int) bool {
		return funcDecls[i].Name.Name < funcDecls[j].Name.Name
	})
	return funcDeclToDecl(funcDecls)
}
