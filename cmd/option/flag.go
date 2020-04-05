package option

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// StringFlag represents flag which can be specified as string
type StringFlag struct {
	*BaseFlag
	Value      string
	IsDirName  bool
	IsFileName bool
}

// BoolFlag represents flag which can be specified as bool
type BoolFlag struct {
	*BaseFlag
	Value bool
}

// IntFlag represents flag which can be specified as int
type IntFlag struct {
	*BaseFlag
	Value int
}

// Int8Flag represents flag which can be specified as int8
type Int8Flag struct {
	*BaseFlag
	Value int8
}

// Int16Flag represents flag which can be specified as int16
type Int16Flag struct {
	*BaseFlag
	Value int16
}

// Int32Flag represents flag which can be specified as int32
type Int32Flag struct {
	*BaseFlag
	Value int32
}

// Int64Flag represents flag which can be specified as int64
type Int64Flag struct {
	*BaseFlag
	Value int64
}

// UintFlag represents flag which can be specified as uint
type UintFlag struct {
	*BaseFlag
	Value uint
}

// Uint8Flag represents flag which can be specified as uint8
type Uint8Flag struct {
	*BaseFlag
	Value uint8
}

// Uint16Flag represents flag which can be specified as uint16
type Uint16Flag struct {
	*BaseFlag
	Value uint16
}

// Uint32Flag represents flag which can be specified as uint32
type Uint32Flag struct {
	*BaseFlag
	Value uint32
}

// Uint64Flag represents flag which can be specified as uint64
type Uint64Flag struct {
	*BaseFlag
	Value uint64
}

// Float32Flag represents flag which can be specified as float32
type Float32Flag struct {
	*BaseFlag
	Value float32
}

// Float64Flag represents flag which can be specified as float64
type Float64Flag struct {
	*BaseFlag
	Value float64
}

// RegisterFlags register flags to provided cmd and viper
func RegisterFlags(cmd *cobra.Command, flags []Flag) error {
	for _, flag := range flags {
		if err := RegisterFlag(cmd, flag); err != nil {
			return err
		}
	}
	return nil
}

// RegisterFlag register flag to provided cmd and viper
func RegisterFlag(cmd *cobra.Command, flag Flag) error {
	baseFlag := flag.getBaseFlag()
	flagSet := getFlagSet(cmd, baseFlag)

	var rerr error
	switch f := flag.(type) {
	case *StringFlag:
		rerr = RegisterStringFlag(cmd, f)
	case *BoolFlag:
		rerr = RegisterBoolFlag(cmd, f)
	case *IntFlag:
		rerr = RegisterIntFlag(cmd, f)
	case *Int8Flag:
		rerr = RegisterInt8Flag(cmd, f)
	case *Int16Flag:
		rerr = RegisterInt16Flag(cmd, f)
	case *Int32Flag:
		rerr = RegisterInt32Flag(cmd, f)
	case *Int64Flag:
		rerr = RegisterInt64Flag(cmd, f)
	case *UintFlag:
		rerr = RegisterUintFlag(cmd, f)
	case *Uint8Flag:
		rerr = RegisterUint8Flag(cmd, f)
	case *Uint16Flag:
		rerr = RegisterUint16Flag(cmd, f)
	case *Uint32Flag:
		rerr = RegisterUint32Flag(cmd, f)
	case *Uint64Flag:
		rerr = RegisterUint64Flag(cmd, f)
	case *Float32Flag:
		rerr = RegisterFloat32Flag(cmd, f)
	case *Float64Flag:
		rerr = RegisterFloat64Flag(cmd, f)
	}

	if rerr != nil {
		return rerr
	}

	if err := markAsRequired(cmd, baseFlag); err != nil {
		return err
	}

	if err := viper.BindPFlag(baseFlag.getViperName(), flagSet.Lookup(baseFlag.Name)); err != nil {
		return err
	}
	return nil
}

// RegisterStringFlag register string flag to provided cmd and viper
func RegisterStringFlag(cmd *cobra.Command, flagConfig *StringFlag) error {
	flagSet := getFlagSet(cmd, flagConfig.BaseFlag)
	if flagConfig.Shorthand == "" {
		flagSet.String(flagConfig.Name, flagConfig.Value, flagConfig.Usage)
	} else {
		flagSet.StringP(flagConfig.Name, flagConfig.Shorthand, flagConfig.Value, flagConfig.Usage)
	}

	return markAttributes(cmd, flagConfig)
}

// RegisterBoolFlag register bool flag to provided cmd and viper
func RegisterBoolFlag(cmd *cobra.Command, flagConfig *BoolFlag) error {
	flagSet := getFlagSet(cmd, flagConfig.BaseFlag)
	if flagConfig.Shorthand == "" {
		flagSet.Bool(flagConfig.Name, flagConfig.Value, flagConfig.Usage)
	} else {
		flagSet.BoolP(flagConfig.Name, flagConfig.Shorthand, flagConfig.Value, flagConfig.Usage)
	}
	return nil
}

// RegisterIntFlag register int flag to provided cmd and viper
func RegisterIntFlag(cmd *cobra.Command, flagConfig *IntFlag) error {
	flagSet := getFlagSet(cmd, flagConfig.BaseFlag)
	if flagConfig.Shorthand == "" {
		flagSet.Int(flagConfig.Name, flagConfig.Value, flagConfig.Usage)
	} else {
		flagSet.IntP(flagConfig.Name, flagConfig.Shorthand, flagConfig.Value, flagConfig.Usage)
	}
	return nil
}

// RegisterInt8Flag register int8 flag to provided cmd and viper
func RegisterInt8Flag(cmd *cobra.Command, flagConfig *Int8Flag) error {
	flagSet := getFlagSet(cmd, flagConfig.BaseFlag)
	if flagConfig.Shorthand == "" {
		flagSet.Int8(flagConfig.Name, flagConfig.Value, flagConfig.Usage)
	} else {
		flagSet.Int8P(flagConfig.Name, flagConfig.Shorthand, flagConfig.Value, flagConfig.Usage)
	}
	return nil
}

// RegisterInt16Flag register int16 flag to provided cmd and viper
func RegisterInt16Flag(cmd *cobra.Command, flagConfig *Int16Flag) error {
	flagSet := getFlagSet(cmd, flagConfig.BaseFlag)
	if flagConfig.Shorthand == "" {
		flagSet.Int16(flagConfig.Name, flagConfig.Value, flagConfig.Usage)
	} else {
		flagSet.Int16P(flagConfig.Name, flagConfig.Shorthand, flagConfig.Value, flagConfig.Usage)
	}
	return nil
}

// RegisterInt32Flag register int32 flag to provided cmd and viper
func RegisterInt32Flag(cmd *cobra.Command, flagConfig *Int32Flag) error {
	flagSet := getFlagSet(cmd, flagConfig.BaseFlag)
	if flagConfig.Shorthand == "" {
		flagSet.Int32(flagConfig.Name, flagConfig.Value, flagConfig.Usage)
	} else {
		flagSet.Int32P(flagConfig.Name, flagConfig.Shorthand, flagConfig.Value, flagConfig.Usage)
	}
	return nil
}

// RegisterInt64Flag register int64 flag to provided cmd and viper
func RegisterInt64Flag(cmd *cobra.Command, flagConfig *Int64Flag) error {
	flagSet := getFlagSet(cmd, flagConfig.BaseFlag)
	if flagConfig.Shorthand == "" {
		flagSet.Int64(flagConfig.Name, flagConfig.Value, flagConfig.Usage)
	} else {
		flagSet.Int64P(flagConfig.Name, flagConfig.Shorthand, flagConfig.Value, flagConfig.Usage)
	}
	return nil
}

// RegisterUintFlag register int flag to provided cmd and viper
func RegisterUintFlag(cmd *cobra.Command, flagConfig *UintFlag) error {
	flagSet := getFlagSet(cmd, flagConfig.BaseFlag)
	if flagConfig.Shorthand == "" {
		flagSet.Uint(flagConfig.Name, flagConfig.Value, flagConfig.Usage)
	} else {
		flagSet.UintP(flagConfig.Name, flagConfig.Shorthand, flagConfig.Value, flagConfig.Usage)
	}
	return nil
}

// RegisterUint8Flag register int8 flag to provided cmd and viper
func RegisterUint8Flag(cmd *cobra.Command, flagConfig *Uint8Flag) error {
	flagSet := getFlagSet(cmd, flagConfig.BaseFlag)
	if flagConfig.Shorthand == "" {
		flagSet.Uint8(flagConfig.Name, flagConfig.Value, flagConfig.Usage)
	} else {
		flagSet.Uint8P(flagConfig.Name, flagConfig.Shorthand, flagConfig.Value, flagConfig.Usage)
	}
	return nil
}

// RegisterUint16Flag register int16 flag to provided cmd and viper
func RegisterUint16Flag(cmd *cobra.Command, flagConfig *Uint16Flag) error {
	flagSet := getFlagSet(cmd, flagConfig.BaseFlag)
	if flagConfig.Shorthand == "" {
		flagSet.Uint16(flagConfig.Name, flagConfig.Value, flagConfig.Usage)
	} else {
		flagSet.Uint16P(flagConfig.Name, flagConfig.Shorthand, flagConfig.Value, flagConfig.Usage)
	}
	return nil
}

// RegisterUint32Flag register int32 flag to provided cmd and viper
func RegisterUint32Flag(cmd *cobra.Command, flagConfig *Uint32Flag) error {
	flagSet := getFlagSet(cmd, flagConfig.BaseFlag)
	if flagConfig.Shorthand == "" {
		flagSet.Uint32(flagConfig.Name, flagConfig.Value, flagConfig.Usage)
	} else {
		flagSet.Uint32P(flagConfig.Name, flagConfig.Shorthand, flagConfig.Value, flagConfig.Usage)
	}
	return nil
}

// RegisterUint64Flag register int64 flag to provided cmd and viper
func RegisterUint64Flag(cmd *cobra.Command, flagConfig *Uint64Flag) error {
	flagSet := getFlagSet(cmd, flagConfig.BaseFlag)
	if flagConfig.Shorthand == "" {
		flagSet.Uint64(flagConfig.Name, flagConfig.Value, flagConfig.Usage)
	} else {
		flagSet.Uint64P(flagConfig.Name, flagConfig.Shorthand, flagConfig.Value, flagConfig.Usage)
	}
	return nil
}

// RegisterFloat32Flag register int32 flag to provided cmd and viper
func RegisterFloat32Flag(cmd *cobra.Command, flagConfig *Float32Flag) error {
	flagSet := getFlagSet(cmd, flagConfig.BaseFlag)
	if flagConfig.Shorthand == "" {
		flagSet.Float32(flagConfig.Name, flagConfig.Value, flagConfig.Usage)
	} else {
		flagSet.Float32P(flagConfig.Name, flagConfig.Shorthand, flagConfig.Value, flagConfig.Usage)
	}
	return nil
}

// RegisterFloat64Flag register int64 flag to provided cmd and viper
func RegisterFloat64Flag(cmd *cobra.Command, flagConfig *Float64Flag) error {
	flagSet := getFlagSet(cmd, flagConfig.BaseFlag)
	if flagConfig.Shorthand == "" {
		flagSet.Float64(flagConfig.Name, flagConfig.Value, flagConfig.Usage)
	} else {
		flagSet.Float64P(flagConfig.Name, flagConfig.Shorthand, flagConfig.Value, flagConfig.Usage)
	}
	return nil
}
