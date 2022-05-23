package commands

import (
	"github.com/n-creativesystem/kubernetes-extensions/cmd/helper"
	"github.com/spf13/cobra"
)

func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commands",
		Short: "CommandLine",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()
			helper.ViperForFlags(flags)
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	return cmd
}
