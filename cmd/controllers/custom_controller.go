package controllers

import (
	"fmt"

	"github.com/go-logr/zapr"
	"github.com/n-creativesystem/kubernetes-extensions/pkg/controllers/interfaces"
	"github.com/n-creativesystem/kubernetes-extensions/pkg/logger"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
}

func controllers(leaderElectionID string, reconcilers ...interfaces.ReconcilerConstructor) error {
	ctrl.SetLogger(zapr.NewLogger(logger.GetZapLogger()))
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     config.MetricsBindAddress,
		Port:                   9443,
		HealthProbeBindAddress: config.HealthProbeBindAddress,
		LeaderElection:         config.LeaderElect,
		LeaderElectionID:       leaderElectionID,
	})
	if err != nil {
		return fmt.Errorf("unable to start manager: %s", err)
	}

	for _, reconciler := range reconcilers {
		if err := reconciler(mgr.GetClient(), mgr.GetScheme()).SetupWithManager(mgr); err != nil {
			return fmt.Errorf("unable to create controller: %s", err)
		}
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return fmt.Errorf("unable to set up health check: %s", err)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return fmt.Errorf("unable to set up ready check: %s", err)
	}

	logger.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		return fmt.Errorf("problem running manager: %s", err)
	}
	return nil
}
