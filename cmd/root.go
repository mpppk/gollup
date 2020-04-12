package cmd

import (
	"fmt"
	"go/format"
	"go/token"
	"os"

	"golang.org/x/tools/go/packages"

	"github.com/pkg/errors"

	"github.com/mpppk/gollup/util"

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
		PersistentPreRunE: pPreRunE,
		Args:              cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewRootCmdConfigFromViper()
			if err != nil {
				return err
			}

			pkgs, err := ast2.NewProgramFromPackages(conf.Dirs)
			if err != nil {
				return err
			}
			var main *packages.Package
			for _, p := range pkgs {
				if p.Name == "main" {
					main = p
					break
				}
			}

			if main == nil {
				panic("main is nil")
			}

			_, funcDecls, err := ast2.ExtractCalledFuncsFromFuncDeclRecursive(main.Syntax, main.TypesInfo, "main", []string{})
			if err != nil {
				return err
			}

			newFuncDecls := ast2.CopyFuncDeclsAsDecl(funcDecls)
			file := ast2.NewMergedFileFromPackageInfo(main.Syntax)
			file.Decls = append(file.Decls, newFuncDecls...)
			if err := format.Node(os.Stdout, token.NewFileSet(), file); err != nil {
				return errors.Wrap(err, "failed to output ast to stdout")
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
				Name:  "dirs",
				Usage: "Packages dirs(comma separated)",
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
