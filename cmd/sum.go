package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/mpppk/cli-template/registry"

	"github.com/mpppk/cli-template/cmd/option"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

func newSumCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "sum",
		Short:   "Print sum of arguments",
		Long:    ``,
		Args:    cobra.MinimumNArgs(2),
		Example: "cli-template sum -- -1 2  ->  1",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if _, err := strconv.Atoi(arg); err != nil {
					return fmt.Errorf("failed to convert args to int from %q: %w", arg, err)
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewSumCmdConfigFromViper(args)
			if err != nil {
				return err
			}

			useCase := registry.InitializeSumUseCase(nil)

			var result int
			if conf.Norm {
				log.Println("start L1 Norm calculation")
				r := useCase.CalcL1Norm(conf.Numbers)
				log.Println("finish L1 Norm calculation")
				result = r
			} else {
				log.Println("start sum calculation")
				r := useCase.CalcSum(conf.Numbers)
				log.Println("finish sum calculation")
				result = r
			}

			if conf.HasOut() {
				s := strconv.Itoa(result)
				if err := afero.WriteFile(fs, conf.Out, []byte(s), 777); err != nil {
					return fmt.Errorf("failed to write file to %s: %w", conf.Out, err)
				}
				log.Println("result is written to " + conf.Out)
			} else {
				cmd.Println(result)
			}

			return nil
		},
	}

	if err := registerSumCommandFlags(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

func registerSumCommandFlags(cmd *cobra.Command) error {
	flags := []option.Flag{
		&option.BoolFlag{
			BaseFlag: &option.BaseFlag{
				Name:  "norm",
				Usage: "Calc L1 norm instead of sum",
			},
			Value: false,
		},
		&option.StringFlag{
			BaseFlag: &option.BaseFlag{
				Name:  "out",
				Usage: "Output file path",
			},
		},
	}
	return option.RegisterFlags(cmd, flags)
}

func init() {
	cmdGenerators = append(cmdGenerators, newSumCmd)
}
