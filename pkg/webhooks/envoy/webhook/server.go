package webhook

import (
	"context"

	"github.com/n-creativesystem/kubernetes-extensions/pkg/helper"
	"github.com/n-creativesystem/kubernetes-extensions/pkg/webhooks/hook"
	"github.com/slok/kubewebhook/pkg/log"
	"github.com/slok/kubewebhook/pkg/webhook/mutating"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const prefix = "envoy-injector.ncs.dev"

type envoySidecar struct {
	dockerImage string
	configMaps  []string
}

type webHook struct {
	logger *log.Std
}

var (
	_ hook.Mutator = (*webHook)(nil)
)

func (w *webHook) Mutating() (mutating.MutatorFunc, mutating.WebhookConfig) {
	return w.sideCarInjectMutation, mutating.WebhookConfig{
		Name: "envoySideCarInject",
		Obj:  &corev1.Pod{},
	}
}

func (w *webHook) sideCarInjectMutation(ctx context.Context, obj metav1.Object) (bool, error) {
	w.logger.Debugf("mutation request")
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		// pod以外は何もしない
		return false, nil
	}
	w.logger.Debugf("pod: %#v", pod)
	helper := helper.NewAnnotationHelper(pod.Annotations)
	if !helper.IsEnabled() {
		return false, nil
	}
	dockerImage := helper.GetValueOrDefault("docker-image", viper.GetString("docker-image"))
	sideCar := corev1.Container{
		Name:  "envoy-sidecar",
		Image: dockerImage,
	}
	resources := corev1.ResourceRequirements{
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse("100m"),
			corev1.ResourceMemory: resource.MustParse("200Mi"),
		},
		Limits: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceMemory: resource.MustParse("1Gi"),
		},
	}
	sideCar.Resources = resources

	pod.Spec.Containers = append(pod.Spec.Containers, sideCar)
	return false, nil
}
