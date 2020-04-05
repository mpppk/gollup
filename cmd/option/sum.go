package option

import (
	"fmt"
	"github.com/mpppk/cli-template/util"

	"github.com/spf13/viper"
)

// SumCmdConfig is config for sum command
type SumCmdConfig struct {
	Norm bool
	Out  string
	Numbers []int
}

// NewSumCmdConfigFromViper generate config for sum command from viper
func NewSumCmdConfigFromViper(args []string) (*SumCmdConfig, error) {
	var conf SumCmdConfig
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from viper: %w", err)
	}

	if err := conf.validate(); err != nil {
		return nil, fmt.Errorf("failed to create sum cmd config: %w", err)
	}

	numbers, err := util.ConvertStringSliceToIntSlice(args)
	if err != nil {
		return nil, fmt.Errorf("failed to convert args to numbers. args=%q: %w", args, err)
	}

	conf.Numbers = numbers
	return &conf, nil
}

// HasOut returns whether or not config has Out property
func (c *SumCmdConfig) HasOut() bool {
	return c.Out != ""
}

func (c *SumCmdConfig) validate() error {
	return nil
}
