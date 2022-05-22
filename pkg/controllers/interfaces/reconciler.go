package interfaces

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type ReconcilerConstructor func(client client.Client, schema *runtime.Scheme) Reconciler

type Reconciler interface {
	SetupWithManager(mgr manager.Manager) error
}
