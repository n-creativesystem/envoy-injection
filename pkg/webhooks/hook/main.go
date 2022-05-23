package hook

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/n-creativesystem/kubernetes-extensions/pkg/cert"
	"github.com/n-creativesystem/kubernetes-extensions/pkg/config"
	"github.com/n-creativesystem/kubernetes-extensions/pkg/helper"
	"github.com/n-creativesystem/kubernetes-extensions/pkg/logger"
	webhookhttp "github.com/slok/kubewebhook/pkg/http"
	"github.com/slok/kubewebhook/pkg/webhook/mutating"
	"github.com/slok/kubewebhook/pkg/webhook/validating"
)

type Mutator interface {
	Mutating() (mutating.MutatorFunc, mutating.WebhookConfig)
}

type Validator interface {
	Validating() (validating.ValidatorFunc, validating.WebhookConfig)
}

type Webhook interface{}

func StartServer(ctx context.Context, hook Webhook) error {
	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	cert, err := cert.New(config.WebhookConfig.CertFile, config.WebhookConfig.KeyFile)
	if err != nil {
		return err
	}
	if err := cert.Watch(); err != nil {
		return err
	}

	listen, err := net.Listen("tcp", ":8443")
	if err != nil {
		return fmt.Errorf("Error to create listen: %s", err)
	}

	mux := http.NewServeMux()
	if v, ok := hook.(Mutator); ok {
		mutatorFunc, mutatorConfig := v.Mutating()
		mutator := mutating.MutatorFunc(mutatorFunc)
		webhook, err := mutating.NewWebhook(mutatorConfig, mutator, nil, nil, logger.GetLogger())
		if err != nil {
			return fmt.Errorf("Error to create mutate webhook: %s", err)
		}
		handler, err := webhookhttp.HandlerFor(webhook)
		if err != nil {
			return fmt.Errorf("Error to create mutate webhook handler: %s", err)
		}
		mux.Handle("/mutate", handler)
	}

	if v, ok := hook.(Validator); ok {
		validatorFunc, validatorConfig := v.Validating()
		validator := validating.ValidatorFunc(validatorFunc)
		webhook, err := validating.NewWebhook(validatorConfig, validator, nil, nil, logger.GetLogger())
		if err != nil {
			return fmt.Errorf("Error to create validate webhook: %s", err)
		}
		handler, err := webhookhttp.HandlerFor(webhook)
		if err != nil {
			return fmt.Errorf("Error to create validate webhook handler: %s", err)
		}
		mux.Handle("/validate", handler)
	}

	mux.Handle("/healthz", helper.HealthCheck)
	mux.Handle("/readyz", helper.HealthCheck)

	server := &http.Server{
		Handler: mux,
		TLSConfig: &tls.Config{
			GetCertificate: cert.GetCertificate,
		},
	}

	// graceful shutdown
	trap := make(chan os.Signal, 1)
	signal.Notify(trap, os.Interrupt, syscall.SIGTERM)
	defer func() {
		signal.Stop(trap)
		cancelFunc()
	}()
	go func() {
		timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancelFunc()
		select {
		case sig := <-trap:
			logger.Infof("signal shutdown signal: %s", sig)
			if err := server.Shutdown(timeoutCtx); err != nil {
				logger.Errorf("error server shutdown: %s", err)
			}
		case <-ctx.Done():
		}
	}()

	logger.Infof("start server on :8443")
	if err := server.ServeTLS(listen, config.WebhookConfig.CertFile, config.WebhookConfig.KeyFile); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("Error to start server: %s", err)
	}
	return nil
}
