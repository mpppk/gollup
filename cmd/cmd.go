package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type cmdGenerator func(fs afero.Fs) (*cobra.Command, error)

var cmdGenerators []cmdGenerator
