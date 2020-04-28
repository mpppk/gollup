package cmd

import (
	"bytes"
	"fmt"
	"go/format"
	"go/token"
	"go/types"
	"io"
	"os"

	"github.com/pkg/errors"

	"github.com/mpppk/gollup/util"
	"golang.org/x/tools/imports"

	"github.com/mpppk/gollup/cmd/option"

	"github.com/spf13/afero"

	"github.com/mitchellh/go-homedir"
	ast2 "github.com/mpppk/gollup/ast"
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
		Args:              cobra.ArbitraryArgs,
		PersistentPreRunE: pPreRunE,
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewRootCmdConfigFromViper()
			if err != nil {
				return err
			}

			pkgDirs := args
			if len(args) == 0 {
				pkgDirs = []string{"."}
			}

			pkgs, err := ast2.NewProgramFromPackages(pkgDirs)
			if err != nil {
				return err
			}

			pkg, ok := pkgs.FindPkgByName(conf.TargetPackage)
			if !ok {
				panic("specified packages does not found: " + conf.TargetPackage)
			}

			targetPkg, ok := pkg.Types.Scope().Lookup(conf.TargetMethod).(*types.Func)
			if !ok {
				panic("target is not func: " + conf.TargetPackage + "." + conf.TargetMethod)
			}
			objects, err := ast2.ExtractObjectsFromFuncDeclRecursive(pkgs.Packages, targetPkg, []types.Object{})
			if err != nil {
				return err
			}

			sdecls := ast2.NewDecls(pkgs, objects)
			ast2.RenameExternalPackageFunctions(pkgs, sdecls)
			renamedFuncDecls := ast2.CopyFuncDeclsAsDecl(sdecls.Funcs)
			renamedFuncDecls = ast2.SortFuncDeclsFromDecls(renamedFuncDecls)

			file := ast2.NewMergedFileFromPackageInfo(pkg.Syntax)
			file.Decls = append(file.Decls, ast2.GenDeclToDecl(sdecls.Consts)...)
			file.Decls = append(file.Decls, ast2.GenDeclToDecl(sdecls.Vars)...)
			file.Decls = append(file.Decls, ast2.GenDeclToDecl(sdecls.Types)...)
			file.Decls = append(file.Decls, renamedFuncDecls...)

			buf := new(bytes.Buffer)
			if err := format.Node(buf, token.NewFileSet(), file); err != nil {
				return errors.Wrap(err, "failed to output")
			}
			newSrc, err := formatSrc(buf.Bytes())
			if err != nil {
				return err
			}

			if _, err := io.WriteString(cmd.OutOrStdout(), string(newSrc)); err != nil {
				return err
			}
			return nil
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

func formatSrc(bytes []byte) ([]byte, error) {
	options := &imports.Options{
		TabWidth:  8,
		TabIndent: true,
		Comments:  true,
		Fragment:  true,
	}
	return imports.Process("<standard input>", bytes, options)
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
		&option.StringFlag{
			BaseFlag: &option.BaseFlag{
				Name:  "entrypoint",
				Usage: "Entrypoint",
			},
			Value: "main.main",
		},
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
