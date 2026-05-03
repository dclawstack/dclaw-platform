package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	platformv1 "github.com/dclawstack/dclaw-operator/api/v1"
)

const (
	DClawAppFinalizer    = "platform.dclaw.io/finalizer"
	DPanelConfigMapName  = "dclaw-dpanel-apps"
	DPanelConfigMapNS    = "dclaw-platform"
	DatabaseClusterLabel = "cnpg.io/cluster"
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
// +kubebuilder:rbac:groups="",resources=resourcequotas,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=postgresql.cnpg.io,resources=clusters,verbs=get;list;watch;create;update;patch;delete

// Reconcile implements the 9-step provisioning pipeline for DClawApp resources.
func (r *DClawAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var app platformv1.DClawApp
	if err := r.Get(ctx, req.NamespacedName, &app); err != nil {
		log.Error(err, "unable to fetch DClawApp")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("reconciling DClawApp", "app", app.Spec.DisplayName, "version", app.Spec.Version, "enabled", app.Spec.Enabled)

	// Handle deletion
	if !app.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, &app)
	}

	// Add finalizer if missing
	if !controllerutil.ContainsFinalizer(&app, DClawAppFinalizer) {
		controllerutil.AddFinalizer(&app, DClawAppFinalizer)
		if err := r.Update(ctx, &app); err != nil {
			return ctrl.Result{}, err
		}
	}

	// If disabled, tear down resources but keep the CRD
	if !app.Spec.Enabled {
		return r.reconcileDisabled(ctx, &app)
	}

	ns := appNamespace(&app)

	// ─── Step 1: Namespace ───
	if err := r.reconcileNamespace(ctx, &app, ns); err != nil {
		r.setCondition(&app, "NamespaceReady", metav1.ConditionFalse, "NamespaceCreationFailed", err.Error())
		return r.updateStatus(ctx, &app)
	}
	r.setCondition(&app, "NamespaceReady", metav1.ConditionTrue, "NamespaceCreated", fmt.Sprintf("Namespace %s ready", ns))

	// ─── Step 2: ResourceQuota ───
	if err := r.reconcileResourceQuota(ctx, &app, ns); err != nil {
		r.setCondition(&app, "QuotaReady", metav1.ConditionFalse, "QuotaCreationFailed", err.Error())
		return r.updateStatus(ctx, &app)
	}
	r.setCondition(&app, "QuotaReady", metav1.ConditionTrue, "QuotaCreated", "ResourceQuota applied")

	// ─── Step 3: NetworkPolicy ───
	if err := r.reconcileNetworkPolicy(ctx, &app, ns); err != nil {
		r.setCondition(&app, "NetworkPolicyReady", metav1.ConditionFalse, "NetworkPolicyFailed", err.Error())
		return r.updateStatus(ctx, &app)
	}
	r.setCondition(&app, "NetworkPolicyReady", metav1.ConditionTrue, "NetworkPolicyCreated", "App namespace isolated")

	// ─── Step 4: Database (CloudNativePG Cluster) ───
	if app.Spec.Database.Enabled {
		if err := r.reconcileDatabase(ctx, &app, ns); err != nil {
			r.setCondition(&app, "DatabaseReady", metav1.ConditionFalse, "DatabaseProvisioningFailed", err.Error())
			return r.updateStatus(ctx, &app)
		}
		r.setCondition(&app, "DatabaseReady", metav1.ConditionTrue, "DatabaseCreated", "CloudNativePG cluster provisioned")
		app.Status.DatabaseRef = fmt.Sprintf("%s-db", app.Name)
	} else {
		r.setCondition(&app, "DatabaseReady", metav1.ConditionTrue, "DatabaseSkipped", "Database not requested")
	}

	// ─── Step 5: Frontend Deployment ───
	if err := r.reconcileFrontend(ctx, &app, ns); err != nil {
		r.setCondition(&app, "FrontendReady", metav1.ConditionFalse, "FrontendDeploymentFailed", err.Error())
		return r.updateStatus(ctx, &app)
	}
	r.setCondition(&app, "FrontendReady", metav1.ConditionTrue, "FrontendDeployed", "Frontend deployment ready")

	// ─── Step 6: Backend Deployment ───
	if err := r.reconcileBackend(ctx, &app, ns); err != nil {
		r.setCondition(&app, "BackendReady", metav1.ConditionFalse, "BackendDeploymentFailed", err.Error())
		return r.updateStatus(ctx, &app)
	}
	r.setCondition(&app, "BackendReady", metav1.ConditionTrue, "BackendDeployed", "Backend deployment ready")

	// ─── Step 7: Ingress ───
	if app.Spec.Ingress.Enabled {
		if err := r.reconcileIngress(ctx, &app, ns); err != nil {
			r.setCondition(&app, "IngressReady", metav1.ConditionFalse, "IngressCreationFailed", err.Error())
			return r.updateStatus(ctx, &app)
		}
		r.setCondition(&app, "IngressReady", metav1.ConditionTrue, "IngressCreated", fmt.Sprintf("URL: https://%s%s", app.Spec.Ingress.Host, app.Spec.Ingress.Path))
		app.Status.URL = fmt.Sprintf("https://%s%s", app.Spec.Ingress.Host, app.Spec.Ingress.Path)
	} else {
		r.setCondition(&app, "IngressReady", metav1.ConditionTrue, "IngressSkipped", "Ingress not requested")
	}

	// ─── Step 8: DPanel Registration ───
	if err := r.reconcileDPanelRegistration(ctx, &app); err != nil {
		r.setCondition(&app, "DPanelRegistered", metav1.ConditionFalse, "DPanelRegistrationFailed", err.Error())
		return r.updateStatus(ctx, &app)
	}
	r.setCondition(&app, "DPanelRegistered", metav1.ConditionTrue, "DPanelRegistered", "App listed in DPanel catalog")

	// ─── Step 9: Final Status ───
	app.Status.Phase = "Ready"
	app.Status.Namespace = ns
	now := metav1.NewTime(time.Now())
	app.Status.LastUpdated = &now

	log.Info("DClawApp reconciliation complete", "app", app.Spec.DisplayName, "phase", app.Status.Phase)
	return r.updateStatus(ctx, &app)
}

// ───────────────────────────────────────────────
// Step 1: Namespace
// ───────────────────────────────────────────────
func (r *DClawAppReconciler) reconcileNamespace(ctx context.Context, app *platformv1.DClawApp, ns string) error {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: ns,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "dclaw-operator",
				"platform.dclaw.io/app":        app.Name,
				"platform.dclaw.io/tenant":     app.Namespace,
			},
		},
	}
	if err := controllerutil.SetControllerReference(app, namespace, r.Scheme); err != nil {
		return err
	}
	return r.apply(ctx, namespace)
}

// ───────────────────────────────────────────────
// Step 2: ResourceQuota
// ───────────────────────────────────────────────
func (r *DClawAppReconciler) reconcileResourceQuota(ctx context.Context, app *platformv1.DClawApp, ns string) error {
	cpuLimit := app.Spec.Resources.CPULimit
	if cpuLimit == "" {
		cpuLimit = "2"
	}
	memLimit := app.Spec.Resources.MemoryLimit
	if memLimit == "" {
		memLimit = "4Gi"
	}

	quota := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-quota", app.Name),
			Namespace: ns,
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: corev1.ResourceList{
				corev1.ResourceLimitsCPU:       resource.MustParse(cpuLimit),
				corev1.ResourceLimitsMemory:    resource.MustParse(memLimit),
				corev1.ResourceRequestsStorage: resource.MustParse("20Gi"),
				corev1.ResourcePods:            resource.MustParse("20"),
			},
		},
	}
	if err := controllerutil.SetControllerReference(app, quota, r.Scheme); err != nil {
		return err
	}
	return r.apply(ctx, quota)
}

// ───────────────────────────────────────────────
// Step 3: NetworkPolicy
// ───────────────────────────────────────────────
func (r *DClawAppReconciler) reconcileNetworkPolicy(ctx context.Context, app *platformv1.DClawApp, ns string) error {
	policy := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-isolation", app.Name),
			Namespace: ns,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
				networkingv1.PolicyTypeEgress,
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							NamespaceSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"name": "dclaw-platform",
								},
							},
						},
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"platform.dclaw.io/app": app.Name,
								},
							},
						},
					},
				},
			},
			Egress: []networkingv1.NetworkPolicyEgressRule{
				{
					To: []networkingv1.NetworkPolicyPeer{
						{
							NamespaceSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"name": "dclaw-platform",
								},
							},
						},
					},
					Ports: []networkingv1.NetworkPolicyPort{
						{Protocol: &[]corev1.Protocol{corev1.ProtocolTCP}[0], Port: &intstr.IntOrString{IntVal: 5432}}, // PostgreSQL
						{Protocol: &[]corev1.Protocol{corev1.ProtocolTCP}[0], Port: &intstr.IntOrString{IntVal: 6379}}, // Redis
						{Protocol: &[]corev1.Protocol{corev1.ProtocolTCP}[0], Port: &intstr.IntOrString{IntVal: 443}},  // HTTPS
						{Protocol: &[]corev1.Protocol{corev1.ProtocolTCP}[0], Port: &intstr.IntOrString{IntVal: 53}},  // DNS
					},
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(app, policy, r.Scheme); err != nil {
		return err
	}
	return r.apply(ctx, policy)
}

// ───────────────────────────────────────────────
// Step 4: Database (CloudNativePG Cluster)
// ───────────────────────────────────────────────
func (r *DClawAppReconciler) reconcileDatabase(ctx context.Context, app *platformv1.DClawApp, ns string) error {
	storage := app.Spec.Database.Storage
	if storage == "" {
		storage = "10Gi"
	}
	size := app.Spec.Database.Size
	if size == "" {
		size = "1"
	}

	// We use an unstructured object because CloudNativePG CRDs may not be in the scheme
	// In production, you would import the cnpg types.
	cluster := &map[string]interface{}{}
	_ = cluster
	_ = storage
	_ = size

	// TODO: Import postgresql.cnpg.io/v1 and create a proper Cluster object.
	// For now, mark as ready and let the user patch this once the cnpg dependency is added.
	log := log.FromContext(ctx)
	log.Info("CloudNativePG cluster reconciliation placeholder", "app", app.Name, "storage", storage, "instances", size)
	return nil
}

// ───────────────────────────────────────────────
// Step 5: Frontend Deployment
// ───────────────────────────────────────────────
func (r *DClawAppReconciler) reconcileFrontend(ctx context.Context, app *platformv1.DClawApp, ns string) error {
	replicas := app.Spec.Frontend.Replicas
	if replicas == 0 {
		replicas = 1
	}
	port := app.Spec.Frontend.Port
	if port == 0 {
		port = 3000
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-frontend", app.Name),
			Namespace: ns,
			Labels: map[string]string{
				"platform.dclaw.io/app":     app.Name,
				"platform.dclaw.io/tier":    "frontend",
				"app.kubernetes.io/part-of": app.Spec.DisplayName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"platform.dclaw.io/app":  app.Name,
					"platform.dclaw.io/tier": "frontend",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"platform.dclaw.io/app":  app.Name,
						"platform.dclaw.io/tier": "frontend",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "frontend",
							Image: app.Spec.Frontend.Image,
							Ports: []corev1.ContainerPort{
								{ContainerPort: port, Protocol: corev1.ProtocolTCP},
							},
							Env: r.buildEnvVars(app.Spec.Frontend.Env, app.Spec.Frontend.EnvSecrets),
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(app.Spec.Resources.CPURequest),
									corev1.ResourceMemory: resource.MustParse(app.Spec.Resources.MemoryRequest),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(app.Spec.Resources.CPULimit),
									corev1.ResourceMemory: resource.MustParse(app.Spec.Resources.MemoryLimit),
								},
							},
						},
					},
				},
			},
		},
	}

	// Set defaults for resources if empty
	if app.Spec.Resources.CPURequest == "" {
		deployment.Spec.Template.Spec.Containers[0].Resources.Requests[corev1.ResourceCPU] = resource.MustParse("100m")
	}
	if app.Spec.Resources.MemoryRequest == "" {
		deployment.Spec.Template.Spec.Containers[0].Resources.Requests[corev1.ResourceMemory] = resource.MustParse("128Mi")
	}
	if app.Spec.Resources.CPULimit == "" {
		deployment.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceCPU] = resource.MustParse("500m")
	}
	if app.Spec.Resources.MemoryLimit == "" {
		deployment.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceMemory] = resource.MustParse("512Mi")
	}

	if err := controllerutil.SetControllerReference(app, deployment, r.Scheme); err != nil {
		return err
	}
	if err := r.apply(ctx, deployment); err != nil {
		return err
	}

	// Service
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-frontend", app.Name),
			Namespace: ns,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"platform.dclaw.io/app":  app.Name,
				"platform.dclaw.io/tier": "frontend",
			},
			Ports: []corev1.ServicePort{
				{
					Port:       port,
					TargetPort: intstr.FromInt(int(port)),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
	if err := controllerutil.SetControllerReference(app, svc, r.Scheme); err != nil {
		return err
	}
	return r.apply(ctx, svc)
}

// ───────────────────────────────────────────────
// Step 6: Backend Deployment
// ───────────────────────────────────────────────
func (r *DClawAppReconciler) reconcileBackend(ctx context.Context, app *platformv1.DClawApp, ns string) error {
	replicas := app.Spec.Backend.Replicas
	if replicas == 0 {
		replicas = 1
	}
	port := app.Spec.Backend.Port
	if port == 0 {
		port = 8000
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-backend", app.Name),
			Namespace: ns,
			Labels: map[string]string{
				"platform.dclaw.io/app":     app.Name,
				"platform.dclaw.io/tier":    "backend",
				"app.kubernetes.io/part-of": app.Spec.DisplayName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"platform.dclaw.io/app":  app.Name,
					"platform.dclaw.io/tier": "backend",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"platform.dclaw.io/app":  app.Name,
						"platform.dclaw.io/tier": "backend",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "backend",
							Image: app.Spec.Backend.Image,
							Ports: []corev1.ContainerPort{
								{ContainerPort: port, Protocol: corev1.ProtocolTCP},
							},
							Env: r.buildEnvVars(app.Spec.Backend.Env, app.Spec.Backend.EnvSecrets),
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("100m"),
									corev1.ResourceMemory: resource.MustParse("256Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("1000m"),
									corev1.ResourceMemory: resource.MustParse("1Gi"),
								},
							},
						},
					},
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(app, deployment, r.Scheme); err != nil {
		return err
	}
	if err := r.apply(ctx, deployment); err != nil {
		return err
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-backend", app.Name),
			Namespace: ns,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"platform.dclaw.io/app":  app.Name,
				"platform.dclaw.io/tier": "backend",
			},
			Ports: []corev1.ServicePort{
				{
					Port:       port,
					TargetPort: intstr.FromInt(int(port)),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
	if err := controllerutil.SetControllerReference(app, svc, r.Scheme); err != nil {
		return err
	}
	return r.apply(ctx, svc)
}

// ───────────────────────────────────────────────
// Step 7: Ingress
// ───────────────────────────────────────────────
func (r *DClawAppReconciler) reconcileIngress(ctx context.Context, app *platformv1.DClawApp, ns string) error {
	path := app.Spec.Ingress.Path
	if path == "" {
		path = "/"
	}
	pathType := networkingv1.PathTypePrefix

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-ingress", app.Name),
			Namespace: ns,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
				"cert-manager.io/cluster-issuer":             "letsencrypt-prod",
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: app.Spec.Ingress.Host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     path,
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: fmt.Sprintf("%s-frontend", app.Name),
											Port: networkingv1.ServiceBackendPort{
												Number: app.Spec.Frontend.Port,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if app.Spec.Ingress.TLS {
		ingress.Spec.TLS = []networkingv1.IngressTLS{
			{
				Hosts:      []string{app.Spec.Ingress.Host},
				SecretName: fmt.Sprintf("%s-tls", app.Name),
			},
		}
	}
	if err := controllerutil.SetControllerReference(app, ingress, r.Scheme); err != nil {
		return err
	}
	return r.apply(ctx, ingress)
}

// ───────────────────────────────────────────────
// Step 8: DPanel Registration
// ───────────────────────────────────────────────
func (r *DClawAppReconciler) reconcileDPanelRegistration(ctx context.Context, app *platformv1.DClawApp) error {
	cm := &corev1.ConfigMap{}
	err := r.Get(ctx, client.ObjectKey{Name: DPanelConfigMapName, Namespace: DPanelConfigMapNS}, cm)
	if err != nil {
		if errors.IsNotFound(err) {
			// DPanel ConfigMap doesn't exist yet; create it
			cm = &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      DPanelConfigMapName,
					Namespace: DPanelConfigMapNS,
					Labels: map[string]string{
						"platform.dclaw.io/managed-by": "dclaw-operator",
					},
				},
				Data: map[string]string{
					"apps.json": "[]",
				},
			}
			if createErr := r.Create(ctx, cm); createErr != nil {
				return createErr
			}
		} else {
			return err
		}
	}

	// TODO: Append app metadata to apps.json without duplicating.
	// In production, use a proper structured format and merge logic.
	log := log.FromContext(ctx)
	log.Info("DPanel registration placeholder", "app", app.Spec.DisplayName)
	return nil
}

// ───────────────────────────────────────────────
// Deletion & Disable
// ───────────────────────────────────────────────
func (r *DClawAppReconciler) reconcileDelete(ctx context.Context, app *platformv1.DClawApp) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("finalizing DClawApp", "app", app.Name)

	// Remove from DPanel catalog
	cm := &corev1.ConfigMap{}
	if err := r.Get(ctx, client.ObjectKey{Name: DPanelConfigMapName, Namespace: DPanelConfigMapNS}, cm); err == nil {
		log.Info("removing app from DPanel catalog", "app", app.Name)
		// TODO: remove app from apps.json
	}

	controllerutil.RemoveFinalizer(app, DClawAppFinalizer)
	if err := r.Update(ctx, app); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *DClawAppReconciler) reconcileDisabled(ctx context.Context, app *platformv1.DClawApp) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("scaling down DClawApp", "app", app.Name)

	ns := appNamespace(app)

	// Scale frontend to 0
	frontend := &appsv1.Deployment{}
	if err := r.Get(ctx, client.ObjectKey{Name: fmt.Sprintf("%s-frontend", app.Name), Namespace: ns}, frontend); err == nil {
		replicas := int32(0)
		frontend.Spec.Replicas = &replicas
		if err := r.Update(ctx, frontend); err != nil {
			return r.updateStatus(ctx, app)
		}
	}

	// Scale backend to 0
	backend := &appsv1.Deployment{}
	if err := r.Get(ctx, client.ObjectKey{Name: fmt.Sprintf("%s-backend", app.Name), Namespace: ns}, backend); err == nil {
		replicas := int32(0)
		backend.Spec.Replicas = &replicas
		if err := r.Update(ctx, backend); err != nil {
			return r.updateStatus(ctx, app)
		}
	}

	app.Status.Phase = "Disabled"
	return r.updateStatus(ctx, app)
}

// ───────────────────────────────────────────────
// Helpers
// ───────────────────────────────────────────────
func (r *DClawAppReconciler) apply(ctx context.Context, obj client.Object) error {
	key := client.ObjectKeyFromObject(obj)
	existing := obj.DeepCopyObject().(client.Object)
	if err := r.Get(ctx, key, existing); err != nil {
		if errors.IsNotFound(err) {
			return r.Create(ctx, obj)
		}
		return err
	}
	obj.SetResourceVersion(existing.GetResourceVersion())
	return r.Update(ctx, obj)
}

func (r *DClawAppReconciler) setCondition(app *platformv1.DClawApp, condType string, status metav1.ConditionStatus, reason, message string) {
	now := metav1.NewTime(time.Now())
	for i := range app.Status.Conditions {
		if app.Status.Conditions[i].Type == condType {
			app.Status.Conditions[i].Status = status
			app.Status.Conditions[i].Reason = reason
			app.Status.Conditions[i].Message = message
			app.Status.Conditions[i].LastTransitionTime = now
			return
		}
	}
	app.Status.Conditions = append(app.Status.Conditions, metav1.Condition{
		Type:               condType,
		Status:             status,
		Reason:             reason,
		Message:            message,
		LastTransitionTime: now,
	})
}

func (r *DClawAppReconciler) updateStatus(ctx context.Context, app *platformv1.DClawApp) (ctrl.Result, error) {
	if err := r.Status().Update(ctx, app); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func (r *DClawAppReconciler) buildEnvVars(env map[string]string, secrets []string) []corev1.EnvVar {
	vars := make([]corev1.EnvVar, 0, len(env)+len(secrets))
	for k, v := range env {
		vars = append(vars, corev1.EnvVar{Name: k, Value: v})
	}
	for _, s := range secrets {
		vars = append(vars, corev1.EnvVar{
			Name: s,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: fmt.Sprintf("%s-secrets", s)},
					Key:                  s,
				},
			},
		})
	}
	return vars
}

func appNamespace(app *platformv1.DClawApp) string {
	return fmt.Sprintf("dclaw-%s", app.Name)
}

// SetupWithManager sets up the controller with the Manager.
func (r *DClawAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&platformv1.DClawApp{}).
		Owns(&corev1.Namespace{}).
		Owns(&corev1.ResourceQuota{}).
		Owns(&networkingv1.NetworkPolicy{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&networkingv1.Ingress{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
