package webhook

import (
	"context"

	"github.com/n-creativesystem/kubernetes-extensions/pkg/helper"
	"github.com/n-creativesystem/kubernetes-extensions/pkg/webhooks/hook"
	"github.com/slok/kubewebhook/pkg/webhook/mutating"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const prefix = "ns-injector.ncs.dev"

type Webhook struct {
}

var (
	_ hook.Mutator = (*Webhook)(nil)
)

func (h *Webhook) Mutating() (mutating.MutatorFunc, mutating.WebhookConfig) {
	config := mutating.WebhookConfig{
		Name: "namespaceMutating",
		Obj:  &corev1.Namespace{},
	}

	return h.mutatingFunc, config
}

func (h *Webhook) mutatingFunc(ctx context.Context, obj metav1.Object) (bool, error) {
	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		return false, nil
	}

	annotation := helper.NewAnnotationHelper(ns.Annotations)
	if !annotation.IsEnabled() {
		return false, nil
	}
	return false, nil
}
