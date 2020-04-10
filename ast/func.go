package ast

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/loader"
)

func IsErrorFunc(funcDecl *ast.FuncDecl) bool {
	lastResultIdent, ok := extractFuncLastResultIdent(funcDecl)
	if !ok {
		return false
	}
	return lastResultIdent.Name == "error"
}

func getCallExprReturnTypes(prog *loader.Program, currentPkg *loader.PackageInfo, callExpr *ast.CallExpr) ([]string, error) {
	if currentPkg == nil {
		panic("current pkg is nil")
	}

	pkg := currentPkg
	var funcName string
	recvName := ""
	switch fun := callExpr.Fun.(type) {
	case *ast.SelectorExpr:
		xIdent, ok := fun.X.(*ast.Ident)
		funcName = fun.Sel.Name

		if !ok {
			return nil, errors.New("selectorExpr.X is not *ast.Ident")
		}

		foundPkg := false
		for cPkg := range prog.AllPackages {
			if cPkg.Name() == xIdent.Name {
				pkg = prog.AllPackages[cPkg]
				foundPkg = true
				break
			}
		}

		if foundPkg {
			break
		}

		innerMost := currentPkg.Pkg.Scope().Innermost(xIdent.Pos())
		_, obj := innerMost.LookupParent(xIdent.Name, xIdent.Pos())
		if ptrObj, ok := obj.Type().Underlying().(*types.Pointer); ok {
			if namedObj, ok := ptrObj.Elem().(*types.Named); ok {
				recvName = namedObj.Obj().Name()
				for i := 0; i < namedObj.NumMethods(); i++ {
					method := namedObj.Method(i)
					if method.Name() == fun.Sel.Name {
						pkgName := method.Pkg().Name()
						for cPkg := range prog.AllPackages {
							if cPkg.Name() == pkgName {
								pkg = prog.AllPackages[cPkg]
								foundPkg = true
								break
							}
						}
						if foundPkg {
							break
						}
					}
				}
			}
		}

	case *ast.Ident:
		funcName = fun.Name
	case nil:
		panic("callExpr is nil")
	}
	typeNames, ok := getFuncDeclResultTypes(pkg, recvName, funcName)
	if !ok {
		return nil, fmt.Errorf("func not found: name: %v, recv: %v", funcName, recvName)
	}

	return typeNames, nil
}

//func getStructMethodResultTypes(packageInfo *loader.PackageInfo, selector *ast.SelectorExpr) ([]string, bool) {
//	xIdent, ok := selector.X.(*ast.Ident)
//	if !ok {
//		return nil, false
//	}
//
//	if _, ok := getStruct(packageInfo, xIdent); !ok {
//		return nil, false
//	}
//
//	selName := selector.Sel.Name
//	return getFuncDeclResultTypes(packageInfo, selName)
//}

//func getStruct(packageInfo *loader.PackageInfo, ident *ast.Ident) (*types.Struct, bool) {
//	obj := packageInfo.ObjectOf(ident)
//	if obj == nil {
//		return nil, false
//	}
//
//	structObj, ok := obj.Type().Underlying().(*types.Struct)
//	if !ok {
//		return nil, false
//	}
//	return structObj, true
//}

func getFuncDeclResultTypes(packageInfo *loader.PackageInfo, recvName, funcName string) (typeNames []string, ok bool) {
	funcDecl, ok := getFuncDeclByRecvAndMethodName(packageInfo, recvName, funcName)
	if !ok {
		return nil, false
	}

	if funcDecl == nil {
		panic(fmt.Sprintf("funcDecl is nil"))
	}
	if funcDecl.Type == nil {
		panic(fmt.Sprintf("funcDecl.Type is nil: %v", funcDecl.Name))
	}
	if funcDecl.Type.Results == nil {
		panic(fmt.Sprintf("funcDecl.Type.Results is nil: %#v", funcDecl.Type))
	}
	if funcDecl.Type.Results.List == nil {
		panic("funcDecl.Type.Results.List is nil")
	}
	results := funcDecl.Type.Results.List
	for _, result := range results {
		if typeIdent, ok := result.Type.(*ast.Ident); ok {
			typeNames = append(typeNames, typeIdent.Name)
		}
	}
	return typeNames, true
}

func getFuncDeclByRecvAndMethodName(packageInfo *loader.PackageInfo, recvName, methodName string) (*ast.FuncDecl, bool) {
	return getFuncDecl(packageInfo, func(decl *ast.FuncDecl) bool {
		if decl.Name.Name != methodName {
			return false
		}

		if decl.Recv == nil {
			return recvName == ""
		}

		list := decl.Recv.List
		if len(list) != 1 {
			panic("len(decl.Recv.List) is not 1")
		}

		typeStarExpr, ok := list[0].Type.(*ast.StarExpr)
		if !ok {
			panic("typeIdent is not *ast.StarExpr")
		}

		ident, ok := typeStarExpr.X.(*ast.Ident)
		return ident.Name == recvName
	})
}

//func getFuncDeclByFuncName(packageInfo *loader.PackageInfo, funcName string) (*ast.FuncDecl, bool) {
//	return getFuncDecl(packageInfo, func(decl *ast.FuncDecl) bool {
//		return decl.Recv == nil && decl.Name.Name == funcName
//	})
//	//if packageInfo == nil {
//	//	panic("packageInfo is nil")
//	//}
//	//if packageInfo.Files == nil {
//	//	panic("packageInfo.Files is nil")
//	//}
//	//for _, file := range packageInfo.Files {
//	//	ast.FileExports(file)
//	//	for _, decl := range file.Decls {
//	//		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
//	//			//fmt.Printf("searching Func(%v) -> %v\n", funcName, funcDecl.Name.Name)
//	//			if funcDecl.Name.Name == funcName {
//	//				return funcDecl, true
//	//			}
//	//		}
//	//	}
//	//}
//	//return nil, false
//}

func getFuncDeclResults(funcDecl *ast.FuncDecl) (newResults []*ast.Field) {
	results := funcDecl.Type.Results.List
	for _, result := range results {
		newResults = append(newResults, result)
	}
	return
}

func getFuncDeclParamNames(funcDecl *ast.FuncDecl) (paramNames []string) {
	list := funcDecl.Type.Params.List
	for _, v := range list {
		for _, name := range v.Names {
			paramNames = append(paramNames, name.Name)
		}
	}
	return
}

func getFuncDecl(packageInfo *loader.PackageInfo, f func(decl *ast.FuncDecl) bool) (*ast.FuncDecl, bool) {
	if packageInfo == nil {
		panic("packageInfo is nil")
	}
	if packageInfo.Files == nil {
		panic("packageInfo.Files is nil")
	}
	for _, file := range packageInfo.Files {
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				if !ast.IsExported(funcDecl.Name.Name) {
					continue
				}

				if f(funcDecl) {
					return funcDecl, true
				}
			}
		}
	}
	return nil, false
}

func addPrefixToFunc(funcDecl *ast.FuncDecl, prefix string) {
	funcNameRunes := []rune(funcDecl.Name.Name)
	funcDecl.Name.Name = prefix + strings.ToUpper(string(funcNameRunes[0])) + string(funcNameRunes[1:])
}
