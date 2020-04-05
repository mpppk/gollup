package option

import (
	"fmt"

	"github.com/spf13/viper"
)

// ServeCmdConfig is config for serve command
type ServeCmdConfig struct {
	Port string
}

// NewServeCmdConfigFromViper generate config for serve command from viper
func NewServeCmdConfigFromViper() (*ServeCmdConfig, error) {
	var conf ServeCmdConfig
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from viper: %w", err)
	}
	return &conf, nil
}
