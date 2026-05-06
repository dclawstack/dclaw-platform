# CRD Reference

## DClawApp

```yaml
apiVersion: platform.dclaw.io/v1
kind: DClawApp
metadata:
  name: string           # Required. App identifier
  namespace: string      # Optional. Defaults to dclaw-system
spec:
  appId: string          # Required. Unique app ID
  appName: string        # Required. Human-readable name
  version: string        # Required. Semver version
  category: string       # Required. App category
  enabled: boolean       # Required. Whether app is active
  
  frontend:
    image: string        # Required. Container image
    replicas: integer    # Optional. Default: 2
    port: integer        # Optional. Default: 3000
    
  backend:
    image: string        # Required. Container image
    replicas: integer    # Optional. Default: 2
    port: integer        # Optional. Default: 8000
    
  database:
    enabled: boolean     # Optional. Default: true
    storage: string      # Optional. Default: "10Gi"
    backups: boolean     # Optional. Default: true
    
  ingress:
    enabled: boolean     # Optional. Default: true
    host: string         # Required if ingress enabled
    tls: boolean         # Optional. Default: true
    
  resources:
    limits:
      cpu: string        # Optional. Default: "1000m"
      memory: string     # Optional. Default: "2Gi"
    requests:
      cpu: string        # Optional. Default: "250m"
      memory: string     # Optional. Default: "512Mi"
      
  branding:
    primaryColor: string # Optional. Hex color code
    
  billing:
    tier: string         # Optional. "free", "pro", "team", "enterprise"
    
status:
  phase: string          # "Pending", "Provisioning", "Running", "Failed", "Terminating"
  conditions:            # Array of conditions
    - type: string
      status: string
      reason: string
      message: string
      lastTransitionTime: string
  frontendUrl: string
  backendUrl: string
  databaseHost: string
```

## Example

```yaml
apiVersion: platform.dclaw.io/v1
kind: DClawApp
metadata:
  name: chat
spec:
  appId: chat
  appName: DClaw Chat
  version: 0.2.0
  category: communication
  enabled: true
  frontend:
    image: ghcr.io/dclawstack/dclaw-chat:latest
    replicas: 2
  backend:
    image: ghcr.io/dclawstack/dclaw-chat-backend:latest
    replicas: 2
  database:
    enabled: true
    storage: 10Gi
  ingress:
    enabled: true
    host: chat.dclawstack.io
    tls: true
  resources:
    limits: { cpu: 1000m, memory: 2Gi }
    requests: { cpu: 250m, memory: 512Mi }
  branding:
    primaryColor: "#3B82F6"
  billing:
    tier: pro
```
