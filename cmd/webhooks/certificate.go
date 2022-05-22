package webhooks

import (
	"github.com/n-creativesystem/kubernetes-extensions/pkg/config"
	"github.com/n-creativesystem/kubernetes-extensions/pkg/logger"
	"github.com/n-creativesystem/kubernetes-extensions/pkg/webhooks/certificate"
	"github.com/n-creativesystem/kubernetes-extensions/pkg/webhooks/hook"
	"github.com/spf13/cobra"
)

func certificateRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificate-injector",
		Short: "start tls certificate injector webhook server",
		Run:   certificateRun,
	}
	return cmd
}

func certificateRun(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	if err := config.WebhookConfig.IsValid(); err != nil {
		logger.Fatal(err.Error())
	}
	webhook := certificate.NewWebhook()
	if err := hook.StartServer(ctx, webhook); err != nil {
		logger.Fatal(err.Error())
	}
}
