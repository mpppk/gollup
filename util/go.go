package util

import (
	"go/types"

	"golang.org/x/tools/go/packages"
)

var standardPackages = make(map[string]struct{})

func init() {
	pkgs, err := packages.Load(nil, "std")
	if err != nil {
		panic(err)
	}

	for _, p := range pkgs {
		standardPackages[p.PkgPath] = struct{}{}
	}
}

func HasPkg(obj types.Object) bool {
	return obj != nil && obj.Pkg() != nil
}

func IsStandardPackage(pkg string) bool {
	_, ok := standardPackages[pkg]
	return ok
}
