package cmd

import (
	"flag"

	"github.com/n-creativesystem/kubernetes-extensions/cmd/controllers"
	"github.com/n-creativesystem/kubernetes-extensions/cmd/webhooks"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "kubernetes-injector",
	Short:         "kubernetes-injector is webhook server to inject envoy sidecar",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	cobra.OnInitialize()
	rootCmd.AddCommand(
		controllers.Commands(),
		webhooks.Commands(),
	)
	flags := rootCmd.PersistentFlags()
	flags.AddGoFlagSet(flag.CommandLine)
}

func Execute() error {
	return rootCmd.Execute()
}
