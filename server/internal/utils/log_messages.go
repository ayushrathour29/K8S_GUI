package utils

// Log message constants for backend logging
const (
	LogWarnInitK8sClient    = "Warning: Error initializing Kubernetes client: %v"
	LogNoEnvFile            = "No .env file found"
	LogLimitedServer        = "Starting server with limited functionality (auth endpoints will still work)"
	LogK8sClientInitSuccess = "Kubernetes client initialized successfully"

	LogFailedListPods       = "Failed to list pods: %v"
	LogFailedEncodePodsList = "Failed to encode pods list: %v"
	LogFailedGetPod         = "Failed to get pod: %v"
	LogFailedEncodePod      = "Failed to encode pod: %v"
	LogFailedDeletePod      = "Failed to delete pod: %v"
	LogFailedGetPodLogs     = "Failed to get pod logs: %v"
	LogFailedWritePodLogs   = "Failed to write pod logs to response: %v"

	LogFailedListNodes       = "Failed to list nodes: %v"
	LogFailedEncodeNodesList = "Failed to encode nodes list: %v"
	LogFailedGetNode         = "Failed to get node: %v"
	LogFailedEncodeNode      = "Failed to encode node: %v"

	LogFailedListNamespaces         = "Failed to list namespaces: %v"
	LogFailedEncodeNamespacesList   = "Failed to encode namespaces list: %v"
	LogFailedGetNamespace           = "Failed to get namespace: %v"
	LogFailedEncodeNamespace        = "Failed to encode namespace: %v"
	LogFailedCreateNamespace        = "Failed to create namespace: %v"
	LogFailedEncodeCreatedNamespace = "Failed to encode created namespace: %v"
	LogFailedDeleteNamespace        = "Failed to delete namespace: %v"

	LogFailedListDeployments         = "Failed to list deployments: %v"
	LogFailedEncodeDeploymentsList   = "Failed to encode deployments list: %v"
	LogFailedGetDeployment           = "Failed to get deployment: %v"
	LogFailedEncodeDeployment        = "Failed to encode deployment: %v"
	LogFailedCreateDeployment        = "Failed to create deployment: %v"
	LogFailedEncodeCreatedDeployment = "Failed to encode created deployment: %v"
	LogFailedUpdateDeployment        = "Failed to update deployment: %v"
	LogFailedEncodeUpdatedDeployment = "Failed to encode updated deployment: %v"
	LogFailedDeleteDeployment        = "Failed to delete deployment: %v"

	LogFailedListEvents                = "Failed to list events: %v"
	LogFailedEncodeEventsList          = "Failed to encode events list: %v"
	LogFailedListEventsNamespace       = "Failed to list events for namespace %s: %v"
	LogFailedEncodeEventsListNamespace = "Failed to encode events list for namespace: %v"

	LogFailedListServices         = "Failed to list services: %v"
	LogFailedEncodeServicesList   = "Failed to encode services list: %v"
	LogFailedGetService           = "Failed to get service: %v"
	LogFailedEncodeService        = "Failed to encode service: %v"
	LogFailedCreateService        = "Failed to create service: %v"
	LogFailedEncodeCreatedService = "Failed to encode created service: %v"
	LogFailedDeleteService        = "Failed to delete service: %v"

	LogFailedGetNodeMetrics                = "Failed to get node %s: %v"
	LogFailedGetNodeMetricsAPI             = "Failed to get metrics for node %s: %v"
	LogFailedEncodeNodeMetricsError        = "Failed to encode node metrics error response: %v"
	LogFailedEncodeNodeMetrics             = "Failed to encode node metrics: %v"
	LogFailedGetNodeMetricsList            = "Failed to get node metrics: %v"
	LogFailedEncodeNodeMetricsList         = "Failed to encode node metrics list: %v"
	LogFailedGetPodMetrics                 = "Failed to get pod metrics: %v"
	LogFailedEncodePodMetricsList          = "Failed to encode pod metrics list: %v"
	LogFailedGetPodMetricsNamespace        = "Failed to get pod metrics for namespace %s: %v"
	LogFailedEncodePodMetricsListNamespace = "Failed to encode pod metrics list for namespace: %v"

	LogFailedGetServerVersion     = "Failed to get server version: %v"
	LogFailedEncodeClusterInfo    = "Failed to encode cluster info: %v"
	LogFailedEncodeClusterHealth  = "Failed to encode cluster health: %v"
	LogFailedEncodeClusterVersion = "Failed to encode cluster version: %v"
)

// User-facing error and status messages
const (
	MsgInvalidRequestBody          = "Invalid request body"
	MsgInvalidCredentials          = "Invalid credentials"
	MsgFailedGenerateToken         = "Failed to generate token"
	MsgAuthorizationHeaderRequired = "Authorization header required"
	MsgTokenRequired               = "Token required"
	MsgInvalidToken                = "Invalid token"
	MsgMissingAuthorizationHeader  = "Missing Authorization header"
	MsgInvalidOrExpiredToken       = "Invalid or expired token"

	MsgFailedListServices  = "Failed to list services"
	MsgServiceNotFound     = "Service not found"
	MsgFailedCreateService = "Failed to create service"
	MsgFailedDeleteService = "Failed to delete service"

	MsgFailedListPods       = "Failed to list pods"
	MsgPodNotFound          = "Pod not found"
	MsgFailedDeletePod      = "Failed to delete pod"
	MsgInvalidTailParameter = "Invalid 'tail' parameter"
	MsgFailedGetPodLogs     = "Failed to get pod logs"

	MsgFailedListNodes = "Failed to list nodes"
	MsgNodeNotFound    = "Node not found"

	MsgFailedListNamespaces  = "Failed to list namespaces"
	MsgNamespaceNotFound     = "Namespace not found"
	MsgFailedCreateNamespace = "Failed to create namespace"
	MsgFailedDeleteNamespace = "Failed to delete namespace"

	MsgFailedListDeployments  = "Failed to list deployments"
	MsgDeploymentNotFound     = "Deployment not found"
	MsgFailedCreateDeployment = "Failed to create deployment"
	MsgFailedUpdateDeployment = "Failed to update deployment"
	MsgFailedDeleteDeployment = "Failed to delete deployment"

	MsgFailedListEvents = "Failed to list events"

	MsgFailedGetClusterInfo    = "Failed to get cluster info"
	MsgFailedGetClusterHealth  = "Failed to get cluster health"
	MsgFailedGetClusterVersion = "Failed to get cluster version"

	MsgFailedGetNodeMetrics = "Failed to get node metrics"
	MsgFailedGetPodMetrics  = "Failed to get pod metrics"

	MsgServerRunningLimited = `{"status": "ok", "message": "Server running (Kubernetes not available)"}`
)
