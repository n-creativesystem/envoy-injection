package helper

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// corev1apply "k8s.io/client-go/applyconfigurations/core/v1"
	// metav1apply "k8s.io/client-go/applyconfigurations/meta/v1"
	"k8s.io/client-go/kubernetes"
	kubernetesConfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	defaultLabels = map[string]string{"ncs.extensions.k8s.io" + "/components": ComponentName}
)

func GetDefaultLabels() map[string]string {
	labels := make(map[string]string, len(defaultLabels))
	// copy
	for key, value := range defaultLabels {
		labels[key] = value
	}
	return labels
}

func NewK8SClient() (kubernetes.Interface, error) {
	kubeConfig, err := kubernetesConfig.GetConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(kubeConfig)
}

func NewSecret(name, namespace string, typ corev1.SecretType) *corev1.Secret {
	return &corev1.Secret{
		// TypeMetaApplyConfiguration:   *metav1apply.TypeMeta().WithKind("Secret").WithAPIVersion("v1"),
		// ObjectMetaApplyConfiguration: metav1apply.ObjectMeta().WithName(name).WithNamespace(namespace),
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    GetDefaultLabels(),
		},
		StringData: make(map[string]string),
		Data:       make(map[string][]byte),
		Type:       typ,
	}
}
