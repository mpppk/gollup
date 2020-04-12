package option

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// RootCmdConfig is config for root command
type RootCmdConfig struct {
	RootRawCmdConfig
	Dirs []string
}

// RootCmdConfig is config for root command
type RootRawCmdConfig struct {
	Verbose bool
	Dirs    string
}

// NewRootCmdConfigFromViper generate config for sum command from viper
func NewRootCmdConfigFromViper() (*RootCmdConfig, error) {
	var rawConf RootRawCmdConfig
	if err := viper.Unmarshal(&rawConf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from viper: %w", err)
	}
	return &RootCmdConfig{
		RootRawCmdConfig: rawConf,
		Dirs:             strings.Split(rawConf.Dirs, ","),
	}, nil
}
