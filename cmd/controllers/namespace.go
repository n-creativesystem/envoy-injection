package controllers

import (
	"github.com/n-creativesystem/kubernetes-extensions/pkg/controllers/interfaces"
	"github.com/n-creativesystem/kubernetes-extensions/pkg/controllers/namespaces"
	"github.com/n-creativesystem/kubernetes-extensions/pkg/logger"
	"github.com/spf13/cobra"
)

func generateCertsControllerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "namespace-cert-generator",
		Short: "namespace create is generate certificate",
		Run:   generateCertsControllerRun,
	}
	return cmd
}

func generateCertsControllerRun(cmd *cobra.Command, args []string) {
	reconciler := []interfaces.ReconcilerConstructor{
		namespaces.NewNamespaceReconciler,
	}
	if err := controllers("ecaf1259.nsc.namespace-cert-generator", reconciler...); err != nil {
		logger.Error(err.Error())
	}
}
