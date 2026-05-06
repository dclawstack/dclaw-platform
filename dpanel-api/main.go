package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	registryConfigMapName = "dclaw-apps-registry"
	registryNamespace     = "dclaw-core"
)

// RegistryApp represents a single app entry from the operator's ConfigMap registry
type RegistryApp struct {
	AppID        string `json:"appId"`
	AppName      string `json:"appName"`
	AppIcon      string `json:"appIcon"`
	Category     string `json:"category"`
	Version      string `json:"version"`
	PrimaryColor string `json:"primaryColor"`
	Path         string `json:"path"`
	Tier         string `json:"tier"`
	Status       string `json:"status"`
	URL          string `json:"url"`
}

// AppListResponse is the response for GET /api/v1/apps
type AppListResponse struct {
	Apps []RegistryApp `json:"apps"`
}

// AppDetailResponse is the response for GET /api/v1/apps/:id
type AppDetailResponse struct {
	App RegistryApp `json:"app"`
}

// ErrorResponse is the error response shape
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type Server struct {
	client kubernetes.Interface
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8088" // NOT 8080 — taken by kubectl port-forward on dev machine. See PORT_REGISTRY.md
	}

	var config *rest.Config
	var err error

	// Try in-cluster config first, then fallback to kubeconfig
	config, err = rest.InClusterConfig()
	if err != nil {
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("Failed to create k8s config: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create k8s client: %v", err)
	}

	server := &Server{client: clientset}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", server.handleHealth)
	mux.HandleFunc("/api/v1/apps", server.handleListApps)
	mux.HandleFunc("/api/v1/apps/", server.handleGetApp)

	// Docs API
	mux.HandleFunc("/api/v1/docs/apps", server.handleListDocApps)
	mux.HandleFunc("/api/v1/docs/apps/", server.handleGetDocApp)
	mux.HandleFunc("/api/v1/docs/ecosystem/", server.handleGetEcosystemDoc)
	mux.HandleFunc("/api/v1/docs/search", server.handleSearchDocs)

	handler := withCORS(mux)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("dpanel-api starting on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleListApps(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	apps, err := s.fetchApps(ctx)
	if err != nil {
		log.Printf("Failed to fetch apps: %v", err)
		writeError(w, http.StatusInternalServerError, "Failed to read app registry")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AppListResponse{Apps: apps})
}

func (s *Server) handleGetApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Path is /api/v1/apps/:id
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/apps/")
	appID := strings.TrimSpace(path)
	if appID == "" {
		writeError(w, http.StatusBadRequest, "App ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	apps, err := s.fetchApps(ctx)
	if err != nil {
		log.Printf("Failed to fetch apps: %v", err)
		writeError(w, http.StatusInternalServerError, "Failed to read app registry")
		return
	}

	for _, app := range apps {
		if app.AppID == appID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(AppDetailResponse{App: app})
			return
		}
	}

	writeError(w, http.StatusNotFound, "App not found")
}

func (s *Server) fetchApps(ctx context.Context) ([]RegistryApp, error) {
	cm, err := s.client.CoreV1().ConfigMaps(registryNamespace).Get(ctx, registryConfigMapName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get configmap: %w", err)
	}

	var apps []RegistryApp
	for _, value := range cm.Data {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		var app RegistryApp
		if err := json.Unmarshal([]byte(value), &app); err != nil {
			log.Printf("Warning: failed to parse app entry: %v", err)
			continue
		}
		apps = append(apps, app)
	}

	return apps, nil
}

// DocAppMeta represents the metadata for an app's documentation
type DocAppMeta struct {
	AppID       string   `json:"appId"`
	Title       string   `json:"title"`
	Version     string   `json:"version"`
	HasDocs     bool     `json:"hasDocs"`
	Nav         []DocNav `json:"nav,omitempty"`
}

// DocNav represents a navigation group in app docs
type DocNav struct {
	Group string   `json:"group"`
	Pages []string `json:"pages"`
}

// DocContentResponse represents a single doc page
type DocContentResponse struct {
	AppID   string `json:"appId"`
	Path    string `json:"path"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// SearchResult represents a single search result
type SearchResult struct {
	AppID   string `json:"appId"`
	AppName string `json:"appName"`
	Path    string `json:"path"`
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
}

func (s *Server) handleListDocApps(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	apps, err := s.fetchApps(ctx)
	if err != nil {
		log.Printf("Failed to fetch apps for docs: %v", err)
		writeError(w, http.StatusInternalServerError, "Failed to read app registry")
		return
	}

	var docApps []DocAppMeta
	for _, app := range apps {
		docApps = append(docApps, DocAppMeta{
			AppID:   app.AppID,
			Title:   app.AppName,
			Version: app.Version,
			HasDocs: true,
			Nav: []DocNav{
				{Group: "Getting Started", Pages: []string{"getting-started/index", "getting-started/installation", "getting-started/quickstart", "getting-started/configuration"}},
				{Group: "Guides", Pages: []string{"guides/index", "guides/use-cases", "guides/best-practices"}},
				{Group: "Reference", Pages: []string{"reference/index", "reference/architecture", "reference/stack", "reference/api"}},
				{Group: "Troubleshooting", Pages: []string{"troubleshooting/index", "troubleshooting/common-issues", "troubleshooting/faq"}},
				{Group: "Releases", Pages: []string{"releases/index", "releases/changelog", "releases/roadmap"}},
			},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"apps": docApps})
}

func (s *Server) handleGetDocApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Path is /api/v1/docs/apps/:id or /api/v1/docs/apps/:id/:path...
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/docs/apps/")
	parts := strings.SplitN(path, "/", 2)
	appID := strings.TrimSpace(parts[0])
	if appID == "" {
		writeError(w, http.StatusBadRequest, "App ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	apps, err := s.fetchApps(ctx)
	if err != nil {
		log.Printf("Failed to fetch apps: %v", err)
		writeError(w, http.StatusInternalServerError, "Failed to read app registry")
		return
	}

	var foundApp *RegistryApp
	for _, app := range apps {
		if app.AppID == appID {
			foundApp = &app
			break
		}
	}
	if foundApp == nil {
		writeError(w, http.StatusNotFound, "App not found")
		return
	}

	// If no doc path specified, return metadata
	if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		meta := DocAppMeta{
			AppID:   foundApp.AppID,
			Title:   foundApp.AppName,
			Version: foundApp.Version,
			HasDocs: true,
			Nav: []DocNav{
				{Group: "Getting Started", Pages: []string{"getting-started/index", "getting-started/installation", "getting-started/quickstart", "getting-started/configuration"}},
				{Group: "Guides", Pages: []string{"guides/index", "guides/use-cases", "guides/best-practices"}},
				{Group: "Reference", Pages: []string{"reference/index", "reference/architecture", "reference/stack", "reference/api"}},
				{Group: "Troubleshooting", Pages: []string{"troubleshooting/index", "troubleshooting/common-issues", "troubleshooting/faq"}},
				{Group: "Releases", Pages: []string{"releases/index", "releases/changelog", "releases/roadmap"}},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(meta)
		return
	}

	// Doc content — placeholder: fetch from GitHub raw in production
	docPath := strings.TrimSpace(parts[1])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DocContentResponse{
		AppID:   foundApp.AppID,
		Path:    docPath,
		Title:   foundApp.AppName + " Documentation",
		Content: "# " + foundApp.AppName + "\n\nDocumentation content for path: " + docPath + "\n\n*Content is fetched from the app's repository at build time.*",
	})
}

func (s *Server) handleGetEcosystemDoc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/docs/ecosystem/")
	if path == "" {
		// List ecosystem docs
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"sections": []map[string]interface{}{
				{"path": "getting-started", "title": "Getting Started"},
				{"path": "architecture", "title": "Architecture"},
				{"path": "reference", "title": "Reference"},
				{"path": "troubleshooting", "title": "Troubleshooting"},
				{"path": "releases", "title": "Releases"},
			},
		})
		return
	}

	// Placeholder: in production, read from bundled docs or ConfigMap
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DocContentResponse{
		AppID:   "ecosystem",
		Path:    path,
		Title:   "Ecosystem Documentation",
		Content: "# Ecosystem Doc\n\nPath: " + path + "\n\n*Content is served from bundled markdown in production.*",
	})
}

func (s *Server) handleSearchDocs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		writeError(w, http.StatusBadRequest, "Query parameter 'q' is required")
		return
	}

	// Placeholder: in production, use full-text search index
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"query":   query,
		"results": []SearchResult{},
		"total":   0,
	})
}

func writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
	})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
