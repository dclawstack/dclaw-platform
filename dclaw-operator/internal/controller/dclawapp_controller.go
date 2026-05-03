package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	platformv1 "github.com/dclawstack/dclaw-operator/api/v1"
)

// DClawAppReconciler reconciles a DClawApp object
type DClawAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=platform.dclaw.io,resources=dclawapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=platform.dclaw.io,resources=dclawapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=platform.dclaw.io,resources=dclawapps/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete

// Reconcile is the main reconciliation loop for DClawApp resources.
// It ensures that the desired state (namespace, deployments, databases, ingress,
// network policies) matches the observed state in the cluster.
func (r *DClawAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var app platformv1.DClawApp
	if err := r.Get(ctx, req.NamespacedName, &app); err != nil {
		log.Error(err, "unable to fetch DClawApp")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("reconciling DClawApp", "app", app.Spec.DisplayName, "version", app.Spec.Version)

	// TODO: Implement reconciliation logic:
	// 1. Ensure isolated namespace exists (or use existing)
	// 2. Provision CloudNativePG cluster if Database.Enabled
	// 3. Create/update frontend Deployment + Service
	// 4. Create/update backend Deployment + Service
	// 5. Create/update Ingress
	// 6. Create/update NetworkPolicy for isolation
	// 7. Register app in DPanel (ConfigMap or API call)
	// 8. Update status (Phase, Conditions, URL, DatabaseRef)

	if !app.Spec.Enabled {
		log.Info("app is disabled, scaling down resources")
		// TODO: scale down or delete owned resources
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DClawAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&platformv1.DClawApp{}).
		Complete(r)
}

// appNamespace returns the target namespace for the app resources.
// Uses the app name as namespace name for isolation.
func appNamespace(app *platformv1.DClawApp) string {
	return fmt.Sprintf("dclaw-%s", app.Name)
}
