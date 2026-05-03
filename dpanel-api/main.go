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
