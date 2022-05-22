/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package namespaces

import (
	"context"
	"fmt"

	"github.com/n-creativesystem/kubernetes-extensions/pkg/cert"
	"github.com/n-creativesystem/kubernetes-extensions/pkg/helper"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// NamespaceReconciler reconciles a Namespace object
type NamespaceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=namespaces/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=namespaces/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Namespace object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *NamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var resource corev1.Namespace
	err := r.Get(ctx, req.NamespacedName, &resource)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	if err := r.reconcile(ctx, resource); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *NamespaceReconciler) reconcile(ctx context.Context, resource corev1.Namespace) error {
	log := log.FromContext(ctx)
	var (
		namespace      = resource.GetName()
		caName         = namespace + "-ca"
		isCreateCA     = true
		serverName     = namespace + "-server"
		isCreateServer = true
	)
	log.Info("Namespace Request", "namespace", namespace)
	if resource.Status.Phase != corev1.NamespaceActive {
		return nil
	}

	annotation := helper.NewAnnotationHelper(resource.GetAnnotations(), "tls", "generate")
	if annotation.GetValue("enable") != "true" {
		return nil
	}
	k8sClient, err := helper.NewK8SClient()
	if err != nil {
		return err
	}
	// var secrets corev1.SecretList
	secrets, err := k8sClient.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, secret := range secrets.Items {
		if secret.Name == caName && secret.Type == corev1.SecretTypeTLS {
			isCreateCA = false
		}
		if secret.Name == serverName && secret.Type == corev1.SecretTypeTLS {
			isCreateServer = false
		}
	}
	createFailures := []error{}
	if isCreateCA && isCreateServer {
		log.Info("Generate TLS Certificate")
		root, server, err := cert.GenerateCertificate([]string{fmt.Sprintf("*.%s.svc", namespace)})
		if err != nil {
			return err
		}
		rootTLS := helper.NewSecret(caName, namespace, "")
		rootTLS.Data[corev1.TLSCertKey] = root.Certificate
		rootTLS.Data[corev1.TLSPrivateKeyKey] = root.Private
		serverTLS := helper.NewSecret(serverName, namespace, "")
		serverTLS.Data[corev1.TLSCertKey] = server.Certificate
		serverTLS.Data[corev1.TLSPrivateKeyKey] = server.Private
		secrets := []*corev1.Secret{
			rootTLS,
			serverTLS,
		}
		for _, secret := range secrets {
			secret := secret
			if _, err := k8sClient.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{}); err != nil && !apierrors.IsAlreadyExists(err) {
				if !apierrors.HasStatusCause(err, corev1.NamespaceTerminatingCause) {
					createFailures = append(createFailures, err)
				}
			}
		}
		log.Info("Complete TLS Certificate")
	}
	return utilerrors.Flatten(utilerrors.NewAggregate(createFailures))
}

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Namespace{}).
		Complete(r)
}
