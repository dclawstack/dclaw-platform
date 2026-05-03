package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DClawAppSpec defines the desired state of DClawApp
type DClawAppSpec struct {
	// AppId is the unique identifier for this app (e.g., "chat", "med", "code")
	AppId string `json:"appId"`

	// AppName is the human-readable app name shown in DPanel
	AppName string `json:"appName"`

	// Description is a short summary of the app
	Description string `json:"description,omitempty"`

	// Version is the semantic version of the app
	Version string `json:"version"`

	// AppIcon is a URL or data URI for the app icon
	AppIcon string `json:"appIcon,omitempty"`

	// Category groups the app in the app store (e.g., "productivity", "clinical", "dev")
	Category string `json:"category,omitempty"`

	// Frontend configures the Next.js frontend deployment
	Frontend AppFrontend `json:"frontend,omitempty"`

	// Backend configures the FastAPI backend deployment
	Backend AppBackend `json:"backend,omitempty"`

	// Database configures the PostgreSQL database provisioned by CloudNativePG
	Database AppDatabase `json:"database,omitempty"`

	// Ingress configures public access to the app
	Ingress AppIngress `json:"ingress,omitempty"`

	// Resources defines compute resources for the app workloads
	Resources AppResources `json:"resources,omitempty"`

	// Branding defines visual customization for the app
	Branding AppBranding `json:"branding,omitempty"`

	// Billing defines subscription tier and pricing
	Billing AppBilling `json:"billing,omitempty"`

	// Enabled determines if the app should be running
	Enabled bool `json:"enabled"`
}

// AppFrontend defines the frontend deployment settings
type AppFrontend struct {
	Image      string            `json:"image,omitempty"`
	Replicas   int32             `json:"replicas,omitempty"`
	Port       int32             `json:"port,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
	EnvSecrets []string          `json:"envSecrets,omitempty"`
}

// AppBackend defines the backend deployment settings
type AppBackend struct {
	Image      string            `json:"image,omitempty"`
	Replicas   int32             `json:"replicas,omitempty"`
	Port       int32             `json:"port,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
	EnvSecrets []string          `json:"envSecrets,omitempty"`
}

// AppDatabase defines the database provisioning settings
type AppDatabase struct {
	Enabled       bool   `json:"enabled,omitempty"`
	Size          string `json:"size,omitempty"`
	Storage       string `json:"storage,omitempty"`
	BackupEnabled bool   `json:"backupEnabled,omitempty"`
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
	Limits   AppResourceLimits   `json:"limits,omitempty"`
	Requests AppResourceRequests `json:"requests,omitempty"`
}

// AppResourceLimits defines hard resource caps
type AppResourceLimits struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// AppResourceRequests defines minimum resource guarantees
type AppResourceRequests struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// AppBranding defines visual customization
type AppBranding struct {
	PrimaryColor string `json:"primaryColor,omitempty"`
	LogoURL      string `json:"logoUrl,omitempty"`
}

// AppBilling defines subscription configuration
type AppBilling struct {
	Tier        string `json:"tier,omitempty"`
	PriceCents  int    `json:"priceCents,omitempty"`
	TrialDays   int    `json:"trialDays,omitempty"`
}

// DClawAppCondition defines a single condition for the app status
type DClawAppCondition struct {
	Type               string      `json:"type"`
	Status             metav1.ConditionStatus `json:"status"`
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
	Reason             string      `json:"reason,omitempty"`
	Message            string      `json:"message,omitempty"`
}

// DClawAppStatus defines the observed state of DClawApp
type DClawAppStatus struct {
	// Phase is the current lifecycle phase (Pending, Provisioning, Ready, Failed, Disabled)
	Phase string `json:"phase,omitempty"`

	// Message is a human-readable status message
	Message string `json:"message,omitempty"`

	// Conditions represent the latest available observations of the app's state
	Conditions []DClawAppCondition `json:"conditions,omitempty"`

	// Namespace is the isolated namespace allocated for this app
	Namespace string `json:"namespace,omitempty"`

	// URL is the public endpoint for the app
	URL string `json:"url,omitempty"`

	// DatabaseRef points to the provisioned CloudNativePG cluster
	DatabaseRef string `json:"databaseRef,omitempty"`

	// LastUpdated is the timestamp of the last successful reconciliation
	LastUpdated *metav1.Time `json:"lastUpdated,omitempty"`
}

// DeepCopyInto copies all properties of this object into another object of the same type.
func (in *DClawAppSpec) DeepCopyInto(out *DClawAppSpec) {
	*out = *in
	out.Frontend = in.Frontend
	out.Backend = in.Backend
	out.Database = in.Database
	out.Ingress = in.Ingress
	out.Resources = in.Resources
	out.Branding = in.Branding
	out.Billing = in.Billing
	if in.Frontend.Env != nil {
		in, out := &in.Frontend.Env, &out.Frontend.Env
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Backend.Env != nil {
		in, out := &in.Backend.Env, &out.Backend.Env
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy creates a new DClawAppSpec.
func (in *DClawAppSpec) DeepCopy() *DClawAppSpec {
	if in == nil {
		return nil
	}
	out := new(DClawAppSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies all properties of this object into another object of the same type.
func (in *DClawAppStatus) DeepCopyInto(out *DClawAppStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]DClawAppCondition, len(*in))
		for i := range *in {
			(*out)[i] = (*in)[i]
		}
	}
	if in.LastUpdated != nil {
		in, out := &in.LastUpdated, &out.LastUpdated
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy creates a new DClawAppStatus.
func (in *DClawAppStatus) DeepCopy() *DClawAppStatus {
	if in == nil {
		return nil
	}
	out := new(DClawAppStatus)
	in.DeepCopyInto(out)
	return out
}

// DClawApp is the Schema for the dclawapps API
type DClawApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DClawAppSpec   `json:"spec,omitempty"`
	Status DClawAppStatus `json:"status,omitempty"`
}

// DeepCopyInto copies all properties of this object into another object of the same type.
func (in *DClawApp) DeepCopyInto(out *DClawApp) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopyObject returns a generically typed copy of an object.
func (in *DClawApp) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopy creates a new DClawApp.
func (in *DClawApp) DeepCopy() *DClawApp {
	if in == nil {
		return nil
	}
	out := new(DClawApp)
	in.DeepCopyInto(out)
	return out
}

// DClawAppList contains a list of DClawApp
type DClawAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DClawApp `json:"items"`
}

// DeepCopyInto copies all properties of this object into another object of the same type.
func (in *DClawAppList) DeepCopyInto(out *DClawAppList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DClawApp, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopyObject returns a generically typed copy of an object.
func (in *DClawAppList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopy creates a new DClawAppList.
func (in *DClawAppList) DeepCopy() *DClawAppList {
	if in == nil {
		return nil
	}
	out := new(DClawAppList)
	in.DeepCopyInto(out)
	return out
}
