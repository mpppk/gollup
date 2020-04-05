package cmd

import (
	"github.com/mpppk/cli-template/util"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newVersionCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		//Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(util.Version)
		},
	}
	return cmd, nil
}

func init() {
	cmdGenerators = append(cmdGenerators, newVersionCmd)
}
