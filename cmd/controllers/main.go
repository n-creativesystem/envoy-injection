package controllers

import (
	"github.com/n-creativesystem/kubernetes-extensions/cmd/helper"
	"github.com/spf13/cobra"
)

type Config struct {
	MetricsBindAddress     string
	HealthProbeBindAddress string
	LeaderElect            bool
}

var config Config

func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use: "controllers",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()
			helper.ViperForFlags(flags)
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmd.AddCommand(
		generateCertsControllerCommand(),
	)
	flag := cmd.PersistentFlags()
	flag.StringVar(&config.MetricsBindAddress, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&config.HealthProbeBindAddress, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&config.LeaderElect, "leader-elect", false, "Enable leader election for controller manager. \nEnabling this will ensure there is only one active controller manager.")
	return cmd
}
