# K8s GUI Backend Server

A Go-based backend server for the Kubernetes GUI application, providing RESTful APIs to manage Kubernetes resources.

## Features

- **Pod Management**: List, get, delete pods and retrieve logs
- **Deployment Management**: CRUD operations for deployments
- **Service Management**: Create, list, get, and delete services
- **Namespace Management**: CRUD operations for namespaces
- **Node Management**: List and get node information
- **Event Monitoring**: List events across namespaces
- **Metrics**: Node and pod resource metrics (requires metrics-server)
- **Cluster Health**: Monitor cluster health and version information

## Project Structure

```
server/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/
│   │   ├── cluster.go           # Cluster-related handlers
│   │   ├── deployments.go       # Deployment-related handlers
│   │   ├── events.go            # Event-related handlers
│   │   ├── metrics.go           # Metrics-related handlers
│   │   ├── namespaces.go        # Namespace-related handlers
│   │   ├── nodes.go             # Node-related handlers
│   │   ├── pods.go              # Pod-related handlers
│   │   └── services.go          # Service-related handlers
│   ├── auth/
│   │   └── auth.go              # Authentication logic
│   ├── k8s/
│   │   └── client.go            # Kubernetes client setup
│   ├── models/
│   │   ├── cluster.go           # Cluster data models
│   │   ├── deployment.go        # Deployment data models
│   │   ├── event.go             # Event data models
│   │   ├── metrics.go           # Metrics data models
│   │   ├── namespace.go         # Namespace data models
│   │   ├── node.go              # Node data models
│   │   ├── pod.go               # Pod data models
│   │   └── service.go           # Service data models
│   ├── server/
│   │   └── router.go            # HTTP router configuration
│   └── utils/
│       └── utils.go             # Utility functions
├── go.mod                       # Go module file
├── go.sum                       # Go module checksums
├── Makefile                     # Build and development tasks
├── air.toml                     # Hot reload configuration
├── .gitignore                   # Git ignore rules
└── README.md                    # This file
```

## Prerequisites

- Go 1.19 or later
- Kubernetes cluster access
- kubectl configured with your cluster

## Installation

1. Clone the repository:

```bash
git clone https://github.com/ayushrathour29/K8S_GUI.git
cd k8s_GUI/server
```

2. Install dependencies:

```bash
go mod download
go mod tidy
```

3. Set up environment variables:

```bash
cp .env.example .env
# Edit .env with your configuration
```

## Development

### Using Makefile

The project includes a comprehensive Makefile with common development tasks:

```bash
# Build the application
make build

# Run the application
make run

# Run with hot reload (requires air)
make dev

# Install air for hot reload
make install-air

# Run tests
make test

# Format code
make fmt

# Run linter
make lint

# Install linter
make install-lint

# Show all available commands
make help
```

### Manual Commands

```bash
# Build
go build -o build/k8s-gui-server cmd/api/main.go

# Run
go run cmd/api/main.go

# Test
go test ./...

# Format
go fmt ./...
```

## API Endpoints

### Authentication

- `POST /api/login` - User authentication
- `GET /api/verify` - Token verification

### Cluster

- `GET /api/clusters` - Get cluster information
- `GET /api/cluster/health` - Get cluster health status
- `GET /api/cluster/version` - Get cluster version

### Pods

- `GET /api/pods` - List all pods
- `GET /api/pods/{namespace}/{name}` - Get specific pod
- `DELETE /api/pods/{namespace}/{name}` - Delete pod
- `GET /api/pods/{namespace}/{name}/logs` - Get pod logs

### Deployments

- `GET /api/deployments` - List all deployments
- `GET /api/deployments/{namespace}/{name}` - Get specific deployment
- `POST /api/deployments` - Create deployment
- `PUT /api/deployments/{namespace}/{name}` - Update deployment
- `DELETE /api/deployments/{namespace}/{name}` - Delete deployment

### Services

- `GET /api/services` - List all services
- `GET /api/services/{namespace}/{name}` - Get specific service
- `POST /api/services` - Create service
- `DELETE /api/services/{namespace}/{name}` - Delete service

### Namespaces

- `GET /api/namespaces` - List all namespaces
- `GET /api/namespaces/{name}` - Get specific namespace
- `POST /api/namespaces` - Create namespace
- `DELETE /api/namespaces/{name}` - Delete namespace

### Nodes

- `GET /api/nodes` - List all nodes
- `GET /api/nodes/{name}` - Get specific node

### Events

- `GET /api/events` - List all events
- `GET /api/events/{namespace}` - List events by namespace

### Metrics

- `GET /api/metrics/nodes` - Get all node metrics
- `GET /api/metrics/nodes/{name}` - Get specific node metrics
- `GET /api/metrics/pods` - Get all pod metrics
- `GET /api/metrics/pods/{namespace}` - Get pod metrics by namespace

## Configuration

The server can be configured using environment variables:

- `PORT`: Server port (default: 8080)
- `HOST`: Server host (default: localhost)
- `KUBECONFIG_PATH`: Path to kubeconfig file
- `JWT_SECRET`: Secret key for JWT tokens
- `JWT_EXPIRY`: JWT token expiry time
- `CORS_ALLOWED_ORIGINS`: Allowed CORS origins
- `LOG_LEVEL`: Logging level
- `DEBUG`: Debug mode

## Docker

Build and run with Docker:

```bash
# Build image
make docker-build

# Run container
make docker-run
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Run tests and linting
6. Submit a pull request

## License

This project is licensed under the MIT License.
