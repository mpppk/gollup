package option

import (
	"fmt"

	"github.com/spf13/viper"
)

// RootCmdConfig is config for root command
type RootCmdConfig struct {
	Verbose bool
}

// NewRootCmdConfigFromViper generate config for sum command from viper
func NewRootCmdConfigFromViper() (*RootCmdConfig, error) {
	var conf RootCmdConfig
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from viper: %w", err)
	}
	return &conf, nil
}
