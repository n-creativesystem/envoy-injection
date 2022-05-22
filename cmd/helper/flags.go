package helper

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func ViperForFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(f *pflag.Flag) {
		_ = viper.BindPFlag(strings.ReplaceAll(strings.ToUpper(f.Name), "-", "_"), f)
	})
}
