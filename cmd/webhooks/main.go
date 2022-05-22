package webhooks

import (
	"flag"

	"github.com/n-creativesystem/kubernetes-extensions/cmd/helper"
	"github.com/n-creativesystem/kubernetes-extensions/pkg/config"
	"github.com/spf13/cobra"
)

func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use: "webhooks",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()
			helper.ViperForFlags(flags)
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmd.AddCommand(
		certificateRunCommand(),
		// envoyRunCommand(),
	)
	flags := cmd.PersistentFlags()
	flags.StringVarP(&config.WebhookConfig.CertFile, "cert-file", "c", "", "Certificate file name of TLS")
	flags.StringVarP(&config.WebhookConfig.KeyFile, "key-file", "k", "", "Key file name of TLS")
	flags.AddGoFlagSet(flag.CommandLine)
	return cmd
}
