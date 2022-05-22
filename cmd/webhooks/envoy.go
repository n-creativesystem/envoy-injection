package webhooks

import (
	"github.com/n-creativesystem/kubernetes-extensions/pkg/logger"
	"github.com/spf13/cobra"
)

func envoyRunCommand() *cobra.Command {
	opt := &envoyWebhookOptions{}
	cmd := &cobra.Command{
		Use:   "envoy-injector",
		Short: "start envoy injector webhook server",
		Run:   opt.run,
	}
	flags := cmd.Flags()
	flags.String("docker-image", "", "sidecar container image name")
	return cmd
}

type envoyWebhookOptions struct {
	certFile string
	keyFile  string
}

func (opt *envoyWebhookOptions) run(cmd *cobra.Command, args []string) {
	if opt.certFile == "" {
		logger.Fatal("cert file is required parameter")
	}
	if opt.keyFile == "" {
		logger.Fatal("key file is required parameter")
	}

}
