package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DClawAppSpec defines the desired state of DClawApp
type DClawAppSpec struct {
	// DisplayName is the human-readable app name shown in DPanel
	DisplayName string `json:"displayName"`

	// Description is a short summary of the app
	Description string `json:"description,omitempty"`

	// Version is the semantic version of the app
	Version string `json:"version"`

	// Icon is a URL or data URI for the app icon
	Icon string `json:"icon,omitempty"`

	// Category groups the app in the app store (e.g., "productivity", "clinical", "dev")
	Category string `json:"category,omitempty"`

	// Frontend configures the Next.js frontend deployment
	Frontend AppFrontend `json:"frontend"`

	// Backend configures the FastAPI backend deployment
	Backend AppBackend `json:"backend"`

	// Database configures the PostgreSQL database provisioned by CloudNativePG
	Database AppDatabase `json:"database,omitempty"`

	// Ingress configures public access to the app
	Ingress AppIngress `json:"ingress,omitempty"`

	// Resources defines compute resources for the app workloads
	Resources AppResources `json:"resources,omitempty"`

	// Enabled determines if the app should be running
	Enabled bool `json:"enabled"`
}

// AppFrontend defines the frontend deployment settings
type AppFrontend struct {
	Image      string            `json:"image"`
	Replicas   int32             `json:"replicas,omitempty"`
	Port       int32             `json:"port,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
	EnvSecrets []string          `json:"envSecrets,omitempty"`
}

// AppBackend defines the backend deployment settings
type AppBackend struct {
	Image      string            `json:"image"`
	Replicas   int32             `json:"replicas,omitempty"`
	Port       int32             `json:"port,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
	EnvSecrets []string          `json:"envSecrets,omitempty"`
}

// AppDatabase defines the database provisioning settings
type AppDatabase struct {
	Enabled    bool   `json:"enabled,omitempty"`
	Size       string `json:"size,omitempty"`
	Storage    string `json:"storage,omitempty"`
	BackupEnabled bool `json:"backupEnabled,omitempty"`
}

// AppIngress defines ingress configuration
type AppIngress struct {
	Enabled   bool     `json:"enabled,omitempty"`
	Host      string   `json:"host,omitempty"`
	Path      string   `json:"path,omitempty"`
	TLS       bool     `json:"tls,omitempty"`
	Whitelist []string `json:"whitelist,omitempty"`
}

// AppResources defines compute resource constraints
type AppResources struct {
	CPULimit      string `json:"cpuLimit,omitempty"`
	MemoryLimit   string `json:"memoryLimit,omitempty"`
	CPURequest    string `json:"cpuRequest,omitempty"`
	MemoryRequest string `json:"memoryRequest,omitempty"`
}

// DClawAppStatus defines the observed state of DClawApp
type DClawAppStatus struct {
	// Phase is the current lifecycle phase (Pending, Provisioning, Ready, Failed, Terminating)
	Phase string `json:"phase,omitempty"`

	// Conditions represent the latest available observations of the app's state
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Namespace is the isolated namespace allocated for this app
	Namespace string `json:"namespace,omitempty"`

	// URL is the public endpoint for the app
	URL string `json:"url,omitempty"`

	// DatabaseRef points to the provisioned CloudNativePG cluster
	DatabaseRef string `json:"databaseRef,omitempty"`

	// LastUpdated is the timestamp of the last successful reconciliation
	LastUpdated *metav1.Time `json:"lastUpdated,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=dca
// +kubebuilder:printcolumn:name="Display",type=string,JSONPath=`.spec.displayName`
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// DClawApp is the Schema for the dclawapps API
type DClawApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DClawAppSpec   `json:"spec,omitempty"`
	Status DClawAppStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DClawAppList contains a list of DClawApp
type DClawAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DClawApp `json:"items"`
}
