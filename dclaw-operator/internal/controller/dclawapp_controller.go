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
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dclawv1 "github.com/dclawstack/dclaw-operator/api/v1"
)

// DClawAppReconciler reconciles a DClawApp object
type DClawAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=platform.dclaw.io,resources=dclawapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=platform.dclaw.io,resources=dclawapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=platform.dclaw.io,resources=dclawapps/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=resourcequotas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=postgresql.cnpg.io,resources=clusters,verbs=get;list;watch;create;update;patch;delete

// Reconcile is the main reconciliation loop
func (r *DClawAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling DClawApp", "name", req.Name, "namespace", req.Namespace)

	// Fetch the DClawApp instance
	dclawApp := &dclawv1.DClawApp{}
	if err := r.Get(ctx, req.NamespacedName, dclawApp); err != nil {
		if errors.IsNotFound(err) {
			log.Info("DClawApp resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get DClawApp")
		return ctrl.Result{}, err
	}

	// Initialize status if needed
	if dclawApp.Status.Phase == "" {
		dclawApp.Status.Phase = "Pending"
		if err := r.Status().Update(ctx, dclawApp); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Step 1: Create Namespace
	if result, err := r.reconcileNamespace(ctx, dclawApp); err != nil {
		return result, err
	}

	// Step 2: Create ResourceQuota
	if result, err := r.reconcileResourceQuota(ctx, dclawApp); err != nil {
		return result, err
	}

	// Step 3: Create NetworkPolicy
	if result, err := r.reconcileNetworkPolicy(ctx, dclawApp); err != nil {
		return result, err
	}

	// Step 4: Create PostgreSQL Cluster
	if result, err := r.reconcileDatabase(ctx, dclawApp); err != nil {
		return result, err
	}

	// Step 5: Create Frontend Deployment
	if result, err := r.reconcileFrontend(ctx, dclawApp); err != nil {
		return result, err
	}

	// Step 6: Create Backend Deployment
	if result, err := r.reconcileBackend(ctx, dclawApp); err != nil {
		return result, err
	}

	// Step 7: Create Ingress
	if result, err := r.reconcileIngress(ctx, dclawApp); err != nil {
		return result, err
	}

	// Step 8: Register in DPanel
	if result, err := r.reconcileDPanelRegistration(ctx, dclawApp); err != nil {
		return result, err
	}

	// Step 9: Update Status to Ready
	if err := r.updateStatus(ctx, dclawApp, "Ready", "All components provisioned"); err != nil {
		return ctrl.Result{}, err
	}

	log.Info("DClawApp reconciliation complete", "name", dclawApp.Name, "phase", dclawApp.Status.Phase)
	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

// reconcileNamespace creates the isolated namespace for the app
func (r *DClawAppReconciler) reconcileNamespace(ctx context.Context, app *dclawv1.DClawApp) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	nsName := fmt.Sprintf("dclaw-%s", app.Spec.AppId)

	ns := &corev1.Namespace{}
	err := r.Get(ctx, types.NamespacedName{Name: nsName}, ns)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating namespace", "name", nsName)
		ns = &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: nsName,
				Labels: map[string]string{
					"managed-by":   "dclaw-operator",
					"appId":        app.Spec.AppId,
					"appName":      app.Spec.AppName,
					"dclaw.io/app": "true",
				},
			},
		}
		if err := controllerutil.SetControllerReference(app, ns, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, ns); err != nil {
			log.Error(err, "Failed to create namespace")
			return ctrl.Result{RequeueAfter: 10 * time.Second}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// reconcileResourceQuota creates resource limits for the app namespace
func (r *DClawAppReconciler) reconcileResourceQuota(ctx context.Context, app *dclawv1.DClawApp) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	nsName := fmt.Sprintf("dclaw-%s", app.Spec.AppId)
	quotaName := fmt.Sprintf("dclaw-%s-quota", app.Spec.AppId)

	quota := &corev1.ResourceQuota{}
	err := r.Get(ctx, types.NamespacedName{Name: quotaName, Namespace: nsName}, quota)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating ResourceQuota", "name", quotaName, "namespace", nsName)

		cpuLimit := app.Spec.Resources.Limits.CPU
		if cpuLimit == "" {
			cpuLimit = "1000m"
		}
		memLimit := app.Spec.Resources.Limits.Memory
		if memLimit == "" {
			memLimit = "2Gi"
		}

		quota = &corev1.ResourceQuota{
			ObjectMeta: metav1.ObjectMeta{
				Name:      quotaName,
				Namespace: nsName,
				Labels: map[string]string{
					"managed-by": "dclaw-operator",
					"appId":      app.Spec.AppId,
				},
			},
			Spec: corev1.ResourceQuotaSpec{
				Hard: corev1.ResourceList{
					corev1.ResourceRequestsCPU:            resource.MustParse(cpuLimit),
					corev1.ResourceRequestsMemory:         resource.MustParse(memLimit),
					corev1.ResourceLimitsCPU:              resource.MustParse(cpuLimit),
					corev1.ResourceLimitsMemory:           resource.MustParse(memLimit),
					corev1.ResourcePersistentVolumeClaims: resource.MustParse("2"),
					corev1.ResourceServicesLoadBalancers:  resource.MustParse("0"),
				},
			},
		}
		if err := controllerutil.SetControllerReference(app, quota, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, quota); err != nil {
			log.Error(err, "Failed to create ResourceQuota")
			return ctrl.Result{RequeueAfter: 10 * time.Second}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// reconcileNetworkPolicy isolates the app namespace
func (r *DClawAppReconciler) reconcileNetworkPolicy(ctx context.Context, app *dclawv1.DClawApp) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	nsName := fmt.Sprintf("dclaw-%s", app.Spec.AppId)
	npName := fmt.Sprintf("dclaw-%s-isolation", app.Spec.AppId)

	np := &networkingv1.NetworkPolicy{}
	err := r.Get(ctx, types.NamespacedName{Name: npName, Namespace: nsName}, np)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating NetworkPolicy", "name", npName, "namespace", nsName)

		tcpProtocol := corev1.ProtocolTCP

		np = &networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      npName,
				Namespace: nsName,
				Labels: map[string]string{
					"managed-by": "dclaw-operator",
					"appId":      app.Spec.AppId,
				},
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
										"name": "dclaw-core",
									},
								},
							},
							{
								NamespaceSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"name": "ingress-nginx",
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
										"name": "dclaw-core",
									},
								},
							},
						},
						Ports: []networkingv1.NetworkPolicyPort{
							{
								Protocol: &tcpProtocol,
								Port:     &intstr.IntOrString{Type: intstr.Int, IntVal: 443},
							},
						},
					},
				},
			},
		}
		if err := controllerutil.SetControllerReference(app, np, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, np); err != nil {
			log.Error(err, "Failed to create NetworkPolicy")
			return ctrl.Result{RequeueAfter: 10 * time.Second}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// reconcileDatabase creates PostgreSQL cluster via CloudNativePG
func (r *DClawAppReconciler) reconcileDatabase(ctx context.Context, app *dclawv1.DClawApp) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	nsName := fmt.Sprintf("dclaw-%s", app.Spec.AppId)
	clusterName := fmt.Sprintf("dclaw-%s-db", app.Spec.AppId)

	log.Info("Database provisioning requested", "cluster", clusterName, "namespace", nsName)
	log.Info("Note: Ensure CloudNativePG operator is installed in the cluster")

	// TODO: Import cnpg types and create actual Cluster CR
	// For now, we log and mark as pending
	// In production, this would create a postgresql.cnpg.io/v1 Cluster resource

	return ctrl.Result{}, nil
}

// reconcileFrontend creates the Next.js frontend deployment
func (r *DClawAppReconciler) reconcileFrontend(ctx context.Context, app *dclawv1.DClawApp) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	nsName := fmt.Sprintf("dclaw-%s", app.Spec.AppId)
	deploymentName := fmt.Sprintf("dclaw-%s-frontend", app.Spec.AppId)

	replicas := int32(2)
	cpuRequest := app.Spec.Resources.Requests.CPU
	if cpuRequest == "" {
		cpuRequest = "250m"
	}
	memRequest := app.Spec.Resources.Requests.Memory
	if memRequest == "" {
		memRequest = "512Mi"
	}
	cpuLimit := app.Spec.Resources.Limits.CPU
	if cpuLimit == "" {
		cpuLimit = "1000m"
	}
	memLimit := app.Spec.Resources.Limits.Memory
	if memLimit == "" {
		memLimit = "2Gi"
	}

	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: deploymentName, Namespace: nsName}, deployment)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating frontend deployment", "name", deploymentName, "namespace", nsName)

		deployment = &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploymentName,
				Namespace: nsName,
				Labels: map[string]string{
					"managed-by":  "dclaw-operator",
					"appId":       app.Spec.AppId,
					"component":   "frontend",
					"app":         fmt.Sprintf("dclaw-%s", app.Spec.AppId),
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app":       fmt.Sprintf("dclaw-%s", app.Spec.AppId),
						"component": "frontend",
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app":       fmt.Sprintf("dclaw-%s", app.Spec.AppId),
							"component": "frontend",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "frontend",
								Image: "ghcr.io/dclawstack/dclaw-" + app.Spec.AppId + ":latest",
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: 3000,
										Name:          "http",
									},
								},
								Resources: corev1.ResourceRequirements{
									Requests: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse(cpuRequest),
										corev1.ResourceMemory: resource.MustParse(memRequest),
									},
									Limits: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse(cpuLimit),
										corev1.ResourceMemory: resource.MustParse(memLimit),
									},
								},
								Env: []corev1.EnvVar{
									{
										Name:  "APP_ID",
										Value: app.Spec.AppId,
									},
									{
										Name:  "NEXT_PUBLIC_APP_NAME",
										Value: app.Spec.AppName,
									},
								},
							},
						},
					},
				},
			},
		}
		if err := controllerutil.SetControllerReference(app, deployment, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, deployment); err != nil {
			log.Error(err, "Failed to create frontend deployment")
			return ctrl.Result{RequeueAfter: 10 * time.Second}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// Create Service for frontend
	serviceName := fmt.Sprintf("dclaw-%s-frontend", app.Spec.AppId)
	svc := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: serviceName, Namespace: nsName}, svc)
	if err != nil && errors.IsNotFound(err) {
		svc = &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceName,
				Namespace: nsName,
				Labels: map[string]string{
					"managed-by": "dclaw-operator",
					"appId":      app.Spec.AppId,
					"component":  "frontend",
				},
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"app":       fmt.Sprintf("dclaw-%s", app.Spec.AppId),
					"component": "frontend",
				},
				Ports: []corev1.ServicePort{
					{
						Port:       80,
						TargetPort: intstr.FromInt(3000),
						Protocol:   corev1.ProtocolTCP,
					},
				},
				Type: corev1.ServiceTypeClusterIP,
			},
		}
		if err := controllerutil.SetControllerReference(app, svc, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, svc); err != nil {
			log.Error(err, "Failed to create frontend service")
			return ctrl.Result{RequeueAfter: 10 * time.Second}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// reconcileBackend creates the FastAPI backend deployment
func (r *DClawAppReconciler) reconcileBackend(ctx context.Context, app *dclawv1.DClawApp) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	nsName := fmt.Sprintf("dclaw-%s", app.Spec.AppId)
	deploymentName := fmt.Sprintf("dclaw-%s-backend", app.Spec.AppId)

	replicas := int32(2)
	cpuRequest := app.Spec.Resources.Requests.CPU
	if cpuRequest == "" {
		cpuRequest = "250m"
	}
	memRequest := app.Spec.Resources.Requests.Memory
	if memRequest == "" {
		memRequest = "512Mi"
	}
	cpuLimit := app.Spec.Resources.Limits.CPU
	if cpuLimit == "" {
		cpuLimit = "1000m"
	}
	memLimit := app.Spec.Resources.Limits.Memory
	if memLimit == "" {
		memLimit = "2Gi"
	}

	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: deploymentName, Namespace: nsName}, deployment)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating backend deployment", "name", deploymentName, "namespace", nsName)

		deployment = &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploymentName,
				Namespace: nsName,
				Labels: map[string]string{
					"managed-by":  "dclaw-operator",
					"appId":       app.Spec.AppId,
					"component":   "backend",
					"app":         fmt.Sprintf("dclaw-%s", app.Spec.AppId),
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app":       fmt.Sprintf("dclaw-%s", app.Spec.AppId),
						"component": "backend",
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app":       fmt.Sprintf("dclaw-%s", app.Spec.AppId),
							"component": "backend",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "backend",
								Image: "ghcr.io/dclawstack/dclaw-" + app.Spec.AppId + "-backend:latest",
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: 8000,
										Name:          "http",
									},
								},
								Resources: corev1.ResourceRequirements{
									Requests: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse(cpuRequest),
										corev1.ResourceMemory: resource.MustParse(memRequest),
									},
									Limits: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse(cpuLimit),
										corev1.ResourceMemory: resource.MustParse(memLimit),
									},
								},
								Env: []corev1.EnvVar{
									{
										Name:  "APP_ID",
										Value: app.Spec.AppId,
									},
									{
										Name:  "DB_HOST",
										Value: fmt.Sprintf("dclaw-%s-db-rw", app.Spec.AppId),
									},
									{
										Name:  "DB_NAME",
										Value: app.Spec.AppId,
									},
									{
										Name:  "SHIELD_ENDPOINT",
										Value: "http://dclaw-shield.dclaw-core.svc.cluster.local:8080",
									},
									{
										Name:  "VOICE_ENDPOINT",
										Value: "http://dclaw-voice.dclaw-core.svc.cluster.local:8080",
									},
								},
							},
						},
					},
				},
			},
		}
		if err := controllerutil.SetControllerReference(app, deployment, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, deployment); err != nil {
			log.Error(err, "Failed to create backend deployment")
			return ctrl.Result{RequeueAfter: 10 * time.Second}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// Create Service for backend
	serviceName := fmt.Sprintf("dclaw-%s-backend", app.Spec.AppId)
	svc := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: serviceName, Namespace: nsName}, svc)
	if err != nil && errors.IsNotFound(err) {
		svc = &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceName,
				Namespace: nsName,
				Labels: map[string]string{
					"managed-by": "dclaw-operator",
					"appId":      app.Spec.AppId,
					"component":  "backend",
				},
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"app":       fmt.Sprintf("dclaw-%s", app.Spec.AppId),
					"component": "backend",
				},
				Ports: []corev1.ServicePort{
					{
						Port:       8000,
						TargetPort: intstr.FromInt(8000),
						Protocol:   corev1.ProtocolTCP,
					},
				},
				Type: corev1.ServiceTypeClusterIP,
			},
		}
		if err := controllerutil.SetControllerReference(app, svc, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, svc); err != nil {
			log.Error(err, "Failed to create backend service")
			return ctrl.Result{RequeueAfter: 10 * time.Second}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// reconcileIngress creates the Ingress resource
func (r *DClawAppReconciler) reconcileIngress(ctx context.Context, app *dclawv1.DClawApp) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	nsName := fmt.Sprintf("dclaw-%s", app.Spec.AppId)
	ingressName := fmt.Sprintf("dclaw-%s", app.Spec.AppId)

	pathType := networkingv1.PathTypePrefix

	ingress := &networkingv1.Ingress{}
	err := r.Get(ctx, types.NamespacedName{Name: ingressName, Namespace: nsName}, ingress)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Ingress", "name", ingressName, "namespace", nsName)

		ingress = &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      ingressName,
				Namespace: nsName,
				Labels: map[string]string{
					"managed-by": "dclaw-operator",
					"appId":      app.Spec.AppId,
				},
				Annotations: map[string]string{
					"cert-manager.io/cluster-issuer":           "letsencrypt",
					"nginx.ingress.kubernetes.io/ssl-redirect": "true",
				},
			},
			Spec: networkingv1.IngressSpec{
				TLS: []networkingv1.IngressTLS{
					{
						Hosts:      []string{app.Spec.Ingress.Host},
						SecretName: fmt.Sprintf("dclaw-%s-tls", app.Spec.AppId),
					},
				},
				Rules: []networkingv1.IngressRule{
					{
						Host: app.Spec.Ingress.Host,
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path:     "/api",
										PathType: &pathType,
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: fmt.Sprintf("dclaw-%s-backend", app.Spec.AppId),
												Port: networkingv1.ServiceBackendPort{
													Number: 8000,
												},
											},
										},
									},
									{
										Path:     "/ws",
										PathType: &pathType,
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: fmt.Sprintf("dclaw-%s-backend", app.Spec.AppId),
												Port: networkingv1.ServiceBackendPort{
													Number: 8000,
												},
											},
										},
									},
									{
										Path:     "/",
										PathType: &pathType,
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: fmt.Sprintf("dclaw-%s-frontend", app.Spec.AppId),
												Port: networkingv1.ServiceBackendPort{
													Number: 80,
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
		if err := controllerutil.SetControllerReference(app, ingress, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, ingress); err != nil {
			log.Error(err, "Failed to create Ingress")
			return ctrl.Result{RequeueAfter: 10 * time.Second}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// reconcileDPanelRegistration registers the app in DPanel
func (r *DClawAppReconciler) reconcileDPanelRegistration(ctx context.Context, app *dclawv1.DClawApp) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	configMapName := "dclaw-apps-registry"
	configMapNamespace := "dclaw-core"

	cm := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: configMapName, Namespace: configMapNamespace}, cm)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating DPanel registry ConfigMap", "name", configMapName)
		cm = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configMapName,
				Namespace: configMapNamespace,
				Labels: map[string]string{
					"managed-by": "dclaw-operator",
					"component":  "dpanel-registry",
				},
			},
			Data: map[string]string{},
		}
		if err := r.Create(ctx, cm); err != nil {
			log.Error(err, "Failed to create DPanel registry")
			return ctrl.Result{RequeueAfter: 10 * time.Second}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// Update the registry with app metadata
	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}

	appJSON := fmt.Sprintf(`{
		"appId": "%s",
		"appName": "%s",
		"appIcon": "%s",
		"category": "%s",
		"version": "%s",
		"primaryColor": "%s",
		"path": "%s",
		"tier": "%s",
		"status": "ready",
		"url": "https://%s"
	}`, app.Spec.AppId, app.Spec.AppName, app.Spec.AppIcon, app.Spec.Category,
		app.Spec.Version, app.Spec.Branding.PrimaryColor, app.Spec.Ingress.Path,
		app.Spec.Billing.Tier, app.Spec.Ingress.Host)

	cm.Data[app.Spec.AppId] = appJSON

	if err := r.Update(ctx, cm); err != nil {
		log.Error(err, "Failed to update DPanel registry")
		return ctrl.Result{RequeueAfter: 10 * time.Second}, err
	}

	log.Info("App registered in DPanel", "appId", app.Spec.AppId)
	return ctrl.Result{}, nil
}

// updateStatus updates the DClawApp status
func (r *DClawAppReconciler) updateStatus(ctx context.Context, app *dclawv1.DClawApp, phase string, message string) error {
	log := log.FromContext(ctx)

	app.Status.Phase = phase
	app.Status.Message = message
	app.Status.URL = fmt.Sprintf("https://%s", app.Spec.Ingress.Host)

	now := metav1.Now()
	condition := dclawv1.DClawAppCondition{
		Type:               "Ready",
		Status:             metav1.ConditionTrue,
		LastTransitionTime: now,
		Reason:             "ReconciliationComplete",
		Message:            message,
	}

	// Update or append condition
	found := false
	for i, c := range app.Status.Conditions {
		if c.Type == "Ready" {
			app.Status.Conditions[i] = condition
			found = true
			break
		}
	}
	if !found {
		app.Status.Conditions = append(app.Status.Conditions, condition)
	}

	if err := r.Status().Update(ctx, app); err != nil {
		log.Error(err, "Failed to update DClawApp status")
		return err
	}

	log.Info("Status updated", "phase", phase, "message", message)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DClawAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dclawv1.DClawApp{}).
		Owns(&corev1.Namespace{}).
		Owns(&corev1.ResourceQuota{}).
		Owns(&networkingv1.NetworkPolicy{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&networkingv1.Ingress{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
