package namespaces

import (
	"github.com/n-creativesystem/kubernetes-extensions/pkg/controllers/interfaces"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewNamespaceReconciler(client client.Client, schema *runtime.Scheme) interfaces.Reconciler {
	return &NamespaceReconciler{
		Client: client,
		Scheme: schema,
	}
}
