package cmd

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"os"

	"github.com/go-toolsmith/astcopy"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/loader"

	"github.com/mpppk/gollup/util"

	"github.com/mpppk/gollup/cmd/option"

	"github.com/spf13/afero"

	"github.com/mitchellh/go-homedir"
	goofyast "github.com/mpppk/gollup/ast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// NewRootCmd generate root cmd
func NewRootCmd(fs afero.Fs) (*cobra.Command, error) {
	pPreRunE := func(cmd *cobra.Command, args []string) error {
		conf, err := option.NewRootCmdConfigFromViper()
		if err != nil {
			return err
		}
		util.InitializeLog(conf.Verbose)
		return nil
	}

	cmd := &cobra.Command{
		Use:               "gollup",
		Short:             "bundle golang sources into single file with tree-shaking",
		SilenceErrors:     true,
		SilenceUsage:      true,
		PersistentPreRunE: pPreRunE,
		Args:              cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args[0])
		},
	}

	if err := registerSubCommands(fs, cmd); err != nil {
		return nil, err
	}

	if err := registerFlags(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

func registerSubCommands(fs afero.Fs, cmd *cobra.Command) error {
	var subCmds []*cobra.Command
	for _, cmdGen := range cmdGenerators {
		subCmd, err := cmdGen(fs)
		if err != nil {
			return err
		}
		subCmds = append(subCmds, subCmd)
	}
	cmd.AddCommand(subCmds...)
	return nil
}

func registerFlags(cmd *cobra.Command) error {
	flags := []option.Flag{
		&option.BoolFlag{
			BaseFlag: &option.BaseFlag{
				Name:         "verbose",
				Shorthand:    "v",
				IsPersistent: true,
				Usage:        "Show more logs",
			}},
	}
	return option.RegisterFlags(cmd, flags)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd, err := NewRootCmd(afero.NewOsFs())
	if err != nil {
		panic(err)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Print(util.PrettyPrintError(err))
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gollup" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gollup")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func run(filePath string) error {
	prog, err := goofyast.NewProgram(filePath)
	if err != nil {
		return err
	}
	main := prog.Package("main")
	_, funcDecls := extractCalledFuncsFromFuncDeclRecursive(main, "main", "main", []string{})

	newFuncDecls := copyFuncDeclsAsDecl(funcDecls)
	file := newMergedFileFromPackageInfo(main)
	file.Decls = append(file.Decls, newFuncDecls...)
	if err := format.Node(os.Stdout, token.NewFileSet(), file); err != nil {
		return errors.Wrap(err, "failed to output ast to stdout")
	}
	return nil
}

func newMergedFileFromPackageInfo(packageInfo *loader.PackageInfo) *ast.File {
	importDecl := goofyast.MergeImportDeclsFromPackageInfo(packageInfo)

	var imports []*ast.ImportSpec
	for _, file := range packageInfo.Files {
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

func copyFuncDeclsAsDecl(funcDecls []*ast.FuncDecl) (newFuncDecls []ast.Decl) {
	for _, decl := range funcDecls {
		newFuncDecls = append(newFuncDecls, astcopy.FuncDecl(decl))
	}
	return
}

func extractCalledFuncsFromFuncDeclRecursive(packageInfo *loader.PackageInfo, packageName, funcName string, foundedFuncNames []string) (funcs []*types.Func, funcDecls []*ast.FuncDecl) {
	funcDecl := findFuncDeclByName(packageInfo.Files, funcName)
	funcDecls = append(funcDecls, funcDecl)
	calledFuncs := extractCalledFuncsFromFuncDecl(packageInfo, funcDecl)
	newFoundedFuncs := make([]string, len(foundedFuncNames))
	copy(newFoundedFuncs, foundedFuncNames)
	newFoundedFuncs = append(newFoundedFuncs, funcDecl.Name.Name)

	for _, f := range calledFuncs {
		if f.Pkg().Name() != packageName {
			continue
		}

		// 既に発見済みの関数の場合はスキップ
		if isFuncName(foundedFuncNames, f.Name()) {
			continue
		}

		funcs = append(funcs, f)

		newFuncs, newFuncDecls := extractCalledFuncsFromFuncDeclRecursive(packageInfo, packageName, f.Name(), newFoundedFuncs)
		funcs = append(funcs, newFuncs...)
		funcDecls = append(funcDecls, newFuncDecls...)
	}
	return
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

// extractCalledFuncsFromFuncDecl は指定したパッケージの指定したfuncDecl内で呼び出されている関数一覧を返す。
func extractCalledFuncsFromFuncDecl(packageInfo *loader.PackageInfo, targetFuncDecl *ast.FuncDecl) (funcs []*types.Func) {
	ast.Inspect(targetFuncDecl, func(node ast.Node) bool {
		if t, _ := node.(*ast.Ident); t != nil {
			obj := packageInfo.Info.ObjectOf(t)
			if tFunc, _ := obj.(*types.Func); tFunc != nil {
				// 自分自身は無視
				if obj.Name() != targetFuncDecl.Name.Name {
					funcs = append(funcs, tFunc)
				}
			}
			return false
		}
		return true
	})
	return funcs
}
