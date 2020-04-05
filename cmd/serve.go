package cmd

import (
	"github.com/mpppk/cli-template/cmd/option"
	"github.com/mpppk/cli-template/registry"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

func newServeCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run server",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewServeCmdConfigFromViper()
			if err != nil {
				return err
			}
			e := registry.InitializeServer(nil)
			e.Logger.Fatal(e.Start(":" + conf.Port))
			return nil
		},
	}
	if err := registerServeCommandFlags(cmd); err != nil {
		return nil, err
	}
	return cmd, nil
}

func registerServeCommandFlags(cmd *cobra.Command) error {
	flags := []option.Flag{
		&option.Uint16Flag{
			BaseFlag: &option.BaseFlag{
				Name:  "port",
				Usage: "server port",
			},
			Value: 1323,
		},
	}

	if err := viper.BindEnv("port"); err != nil {
		return err
	}
	return option.RegisterFlags(cmd, flags)
}

func init() {
	cmdGenerators = append(cmdGenerators, newServeCmd)
}
