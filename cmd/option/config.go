// Package option provides utilities of option handling
package option

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// BaseFlag represents base command line flag
type BaseFlag struct {
	IsPersistent bool
	IsRequired   bool
	Shorthand    string
	Name         string
	Usage        string
	ViperName    string
}

// Flag represents flag which has base flag
type Flag interface {
	getBaseFlag() *BaseFlag
}

func (f *BaseFlag) getBaseFlag() *BaseFlag {
	return f
}

func (f *BaseFlag) getViperName() string {
	if f.ViperName == "" {
		return f.Name
	}
	return f.ViperName
}

func getFlagSet(cmd *cobra.Command, flagConfig *BaseFlag) (flagSet *pflag.FlagSet) {
	if flagConfig.IsPersistent {
		return cmd.PersistentFlags()
	} else {
		return cmd.Flags()
	}
}

func markAttributes(cmd *cobra.Command, flagConfig *StringFlag) error {
	if err := markAsFileName(cmd, flagConfig); err != nil {
		return err
	}
	if err := markAsDirName(cmd, flagConfig); err != nil {
		return err
	}
	if err := markAsRequired(cmd, flagConfig.BaseFlag); err != nil {
		return err
	}
	return nil
}

func markAsFileName(cmd *cobra.Command, stringFlag *StringFlag) error {
	if stringFlag.IsFileName {
		if stringFlag.IsPersistent {
			if err := cmd.MarkPersistentFlagFilename(stringFlag.Name); err != nil {
				return err
			}
		} else {
			if err := cmd.MarkFlagFilename(stringFlag.Name); err != nil {
				return err
			}
		}
	}
	return nil
}

func markAsDirName(cmd *cobra.Command, stringFlag *StringFlag) error {
	if stringFlag.IsDirName {
		if stringFlag.IsPersistent {
			if err := cmd.MarkPersistentFlagDirname(stringFlag.Name); err != nil {
				return err
			}
		} else {
			if err := cmd.MarkFlagDirname(stringFlag.Name); err != nil {
				return err
			}
		}
	}
	return nil
}

func markAsRequired(cmd *cobra.Command, flagConfig *BaseFlag) error {
	if flagConfig.IsRequired {
		if flagConfig.IsPersistent {
			if err := cmd.MarkPersistentFlagRequired(flagConfig.Name); err != nil {
				return err
			}
		} else {
			if err := cmd.MarkFlagRequired(flagConfig.Name); err != nil {
				return err
			}
		}
	}
	return nil
}
