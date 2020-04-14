package option

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// RootCmdConfig is config for root command
type RootCmdConfig struct {
	RootRawCmdConfig
	TargetPackage string
	TargetMethod  string
}

// RootCmdConfig is config for root command
type RootRawCmdConfig struct {
	Verbose    bool
	EntryPoint string
}

// NewRootCmdConfigFromViper generate config for sum command from viper
func NewRootCmdConfigFromViper() (*RootCmdConfig, error) {
	var rawConf RootRawCmdConfig
	if err := viper.Unmarshal(&rawConf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from viper: %w", err)
	}
	pkg, method, err := parseEntryPoint(rawConf.EntryPoint)
	if err != nil {
		return nil, err
	}
	return &RootCmdConfig{
		RootRawCmdConfig: rawConf,
		TargetPackage:    pkg,
		TargetMethod:     method,
	}, nil
}

func parseEntryPoint(entryPoint string) (string, string, error) {
	pkgAndMethod := strings.Split(entryPoint, ".")
	l := len(pkgAndMethod)
	switch l {
	case 1:
		return "main", entryPoint, nil
	case 2:
		return pkgAndMethod[0], pkgAndMethod[1], nil
	default:
		return "", "", errors.New("invalid entryPoint: " + entryPoint)
	}
}
