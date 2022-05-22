package config

import "errors"

type CommandConfig struct{}

type WebhookOptions struct {
	CertFile string
	KeyFile  string
}

var WebhookConfig = WebhookOptions{}

func (opt WebhookOptions) IsValid() error {
	if opt.CertFile == "" {
		return errors.New("cert file is required parameter")
	}
	if opt.KeyFile == "" {
		return errors.New("key file is required parameter")
	}
	return nil
}
