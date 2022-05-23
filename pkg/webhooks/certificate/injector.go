package certificate

import (
	"context"
	"fmt"

	"github.com/n-creativesystem/kubernetes-extensions/pkg/helper"
	"github.com/n-creativesystem/kubernetes-extensions/pkg/logger"
	"github.com/n-creativesystem/kubernetes-extensions/pkg/webhooks/hook"
	"github.com/slok/kubewebhook/pkg/webhook/mutating"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type WebHook struct{}

var (
	_ hook.Mutator = (*WebHook)(nil)
)

func NewWebhook() *WebHook {
	return &WebHook{}
}

func (w *WebHook) Mutating() (mutating.MutatorFunc, mutating.WebhookConfig) {
	return w.mutation, mutating.WebhookConfig{
		Name: "certificateInject",
		Obj:  &appsv1.Deployment{},
	}
}

func (w *WebHook) mutation(ctx context.Context, obj metav1.Object) (bool, error) {
	const volumeName = "certificate-injection"
	logger.Debugf("mutation request")
	dep, ok := obj.(*appsv1.Deployment)
	if !ok {
		return false, nil
	}
	logger.Debugf(fmt.Sprintf("dep: %#v", dep))
	helper := helper.NewAnnotationHelper(dep.Annotations)
	if !helper.IsEnabled() {
		return false, nil
	}
	typ := helper.GetValue("type")
	// タイプがserverもしくはclientが設定されていない場合は何もしない
	if typ != "server" && typ != "client" {
		return false, nil
	}
	// 読み取り専用
	var readOnlyMode int32 = 0644
	items := []corev1.KeyToPath{
		{
			Key:  "tls.crt",
			Path: "server.crt",
			Mode: &readOnlyMode,
		},
	}
	certificateVolumeMounts := []corev1.VolumeMount{
		{
			Name:      volumeName,
			ReadOnly:  true,
			MountPath: "/etc/webhook/cert/server.crt",
			SubPath:   "server.crt",
		},
	}
	if typ == "server" {
		items = append(items, corev1.KeyToPath{
			Key:  "tls.key",
			Path: "server.key",
			Mode: &readOnlyMode,
		})
		certificateVolumeMounts = append(certificateVolumeMounts, corev1.VolumeMount{
			Name:      volumeName,
			ReadOnly:  true,
			MountPath: "/etc/webhook/cert/server.key",
			SubPath:   "server.key",
		})
	}

	namespace := dep.GetNamespace()
	certificateVolumes := corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: namespace + "-certificate",
				Items:      items,
			},
		},
	}
	certificateEnvs := []corev1.EnvVar{
		{
			Name:  "TLS_CERT",
			Value: "/etc/webhook/cert/server.crt",
		},
		{
			Name:  "TLS_KEY",
			Value: "/etc/webhook/cert/server.key",
		},
	}
	dep.Spec.Template.Spec.Volumes = append(dep.Spec.Template.Spec.Volumes, certificateVolumes)
	containers := dep.Spec.Template.Spec.Containers
	for idx := range containers {
		containers[idx].VolumeMounts = append(containers[idx].VolumeMounts, certificateVolumeMounts...)
		containers[idx].Env = append(containers[idx].Env, certificateEnvs...)
	}
	return false, nil
}
