package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

// ClusterInfo represents cluster information
type ClusterInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Nodes   int    `json:"nodes"`
	Healthy bool   `json:"healthy"`
	Status  string `json:"status"`
}

// Pod represents a simplified pod view
type Pod struct {
	Name         string            `json:"name"`
	Namespace    string            `json:"namespace"`
	Status       string            `json:"status"`
	RestartCount int32             `json:"restartCount"`
	CreatedAt    string            `json:"createdAt"`
	NodeName     string            `json:"nodeName"`
	PodIP        string            `json:"podIP"`
	Containers   []string          `json:"containers"`
	Labels       map[string]string `json:"labels,omitempty"`
}

// Deployment represents a simplified deployment view
type Deployment struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	Replicas          int32             `json:"replicas"`
	AvailableReplicas int32             `json:"availableReplicas"`
	CreatedAt         string            `json:"createdAt"`
	Strategy          string            `json:"strategy"`
	Labels            map[string]string `json:"labels,omitempty"`
}

// Service represents a simplified service view
type Service struct {
	Name      string        `json:"name"`
	Namespace string        `json:"namespace"`
	Type      string        `json:"type"`
	ClusterIP string        `json:"clusterIP"`
	Ports     []ServicePort `json:"ports"`
	CreatedAt string        `json:"createdAt"`
}

// ServicePort represents service port information
type ServicePort struct {
	Port     int32  `json:"port"`
	Protocol string `json:"protocol"`
}

// PodListResponse represents pod list response
type PodListResponse struct {
	Items []Pod `json:"items"`
}

// DeploymentListResponse represents deployment list response
type DeploymentListResponse struct {
	Items []Deployment `json:"items"`
}

// ServiceListResponse represents service list response
type ServiceListResponse struct {
	Items []Service `json:"items"`
}

// HandleLogin authenticates a user
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Implementation moved to auth package
}

// VerifyToken checks token validity
func VerifyToken(w http.ResponseWriter, r *http.Request) {
	// Implementation moved to auth package
}

// GetClusters returns cluster info
func GetClusters(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		version, err := clientset.ServerVersion()
		if err != nil {
			log.Printf("Failed to get server version: %v", err)
			http.Error(w, "Failed to get cluster info", http.StatusInternalServerError)
			return
		}

		nodes, err := clientset.CoreV1().Nodes().List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list nodes: %v", err)
			http.Error(w, "Failed to get cluster info", http.StatusInternalServerError)
			return
		}

		clusterInfo := ClusterInfo{
			Name:    "default-cluster",
			Version: version.GitVersion,
			Nodes:   len(nodes.Items),
			Healthy: true,
			Status:  "Healthy",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clusterInfo)
	}
}

// ListPods returns all pods
func ListPods(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pods, err := clientset.CoreV1().Pods("").List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list pods: %v", err)
			http.Error(w, "Failed to list pods", http.StatusInternalServerError)
			return
		}

		response := PodListResponse{Items: make([]Pod, 0, len(pods.Items))}
		for _, p := range pods.Items {
			restartCount := int32(0)
			for _, cs := range p.Status.ContainerStatuses {
				restartCount += cs.RestartCount
			}

			containers := make([]string, len(p.Spec.Containers))
			for i, c := range p.Spec.Containers {
				containers[i] = c.Name
			}

			response.Items = append(response.Items, Pod{
				Name:         p.Name,
				Namespace:    p.Namespace,
				Status:       string(p.Status.Phase),
				RestartCount: restartCount,
				CreatedAt:    p.CreationTimestamp.Time.Format(time.RFC3339),
				NodeName:     p.Spec.NodeName,
				PodIP:        p.Status.PodIP,
				Containers:   containers,
				Labels:       p.Labels,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetPod returns pod details
func GetPod(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		name := vars["name"]

		pod, err := clientset.CoreV1().Pods(namespace).Get(r.Context(), name, metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get pod: %v", err)
			http.Error(w, "Pod not found", http.StatusNotFound)
			return
		}

		restartCount := int32(0)
		for _, cs := range pod.Status.ContainerStatuses {
			restartCount += cs.RestartCount
		}

		containers := make([]string, len(pod.Spec.Containers))
		for i, c := range pod.Spec.Containers {
			containers[i] = c.Name
		}

		response := Pod{
			Name:         pod.Name,
			Namespace:    pod.Namespace,
			Status:       string(pod.Status.Phase),
			RestartCount: restartCount,
			CreatedAt:    pod.CreationTimestamp.Time.Format(time.RFC3339),
			NodeName:     pod.Spec.NodeName,
			PodIP:        pod.Status.PodIP,
			Containers:   containers,
			Labels:       pod.Labels,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// DeletePod deletes a pod
func DeletePod(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		name := vars["name"]

		err := clientset.CoreV1().Pods(namespace).Delete(r.Context(), name, metav1.DeleteOptions{})
		if err != nil {
			log.Printf("Failed to delete pod: %v", err)
			http.Error(w, "Failed to delete pod", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// ListDeployments returns all deployments
func ListDeployments(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deployments, err := clientset.AppsV1().Deployments("").List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list deployments: %v", err)
			http.Error(w, "Failed to list deployments", http.StatusInternalServerError)
			return
		}

		response := DeploymentListResponse{Items: make([]Deployment, 0, len(deployments.Items))}
		for _, d := range deployments.Items {
			strategy := "RollingUpdate"
			if d.Spec.Strategy.Type == appsv1.RecreateDeploymentStrategyType {
				strategy = "Recreate"
			}

			response.Items = append(response.Items, Deployment{
				Name:              d.Name,
				Namespace:         d.Namespace,
				Replicas:          *d.Spec.Replicas,
				AvailableReplicas: d.Status.AvailableReplicas,
				CreatedAt:         d.CreationTimestamp.Time.Format(time.RFC3339),
				Strategy:          strategy,
				Labels:            d.Labels,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetDeployment returns deployment details
func GetDeployment(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		name := vars["name"]

		deployment, err := clientset.AppsV1().Deployments(namespace).Get(r.Context(), name, metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get deployment: %v", err)
			http.Error(w, "Deployment not found", http.StatusNotFound)
			return
		}

		strategy := "RollingUpdate"
		if deployment.Spec.Strategy.Type == appsv1.RecreateDeploymentStrategyType {
			strategy = "Recreate"
		}

		response := Deployment{
			Name:              deployment.Name,
			Namespace:         deployment.Namespace,
			Replicas:          *deployment.Spec.Replicas,
			AvailableReplicas: deployment.Status.AvailableReplicas,
			CreatedAt:         deployment.CreationTimestamp.Time.Format(time.RFC3339),
			Strategy:          strategy,
			Labels:            deployment.Labels,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// CreateDeployment creates a new deployment
func CreateDeployment(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
			Image     string `json:"image"`
			Replicas  int32  `json:"replicas"`
			Port      int32  `json:"port"`
		}
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      req.Name,
				Namespace: req.Namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &req.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": req.Name},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": req.Name},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  req.Name,
								Image: req.Image,
								Ports: []corev1.ContainerPort{
									{ContainerPort: req.Port},
								},
							},
						},
					},
				},
			},
		}

		created, err := clientset.AppsV1().Deployments(req.Namespace).Create(r.Context(), deployment, metav1.CreateOptions{})
		if err != nil {
			http.Error(w, "Failed to create deployment: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message":   "Deployment created successfully",
			"name":      created.Name,
			"namespace": created.Namespace,
		})
	}
}

// UpdateDeployment updates an existing deployment
func UpdateDeployment(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		name := vars["name"]

		type Request struct {
			Image    string `json:"image"`
			Replicas int32  `json:"replicas"`
		}
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		deploy, err := clientset.AppsV1().Deployments(namespace).Get(r.Context(), name, metav1.GetOptions{})
		if err != nil {
			http.Error(w, "Deployment not found: "+err.Error(), http.StatusNotFound)
			return
		}

		deploy.Spec.Replicas = &req.Replicas
		if len(deploy.Spec.Template.Spec.Containers) > 0 {
			deploy.Spec.Template.Spec.Containers[0].Image = req.Image
		}

		_, err = clientset.AppsV1().Deployments(namespace).Update(r.Context(), deploy, metav1.UpdateOptions{})
		if err != nil {
			http.Error(w, "Failed to update deployment: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message":   "Deployment updated successfully",
			"name":      name,
			"namespace": namespace,
		})
	}
}

// DeleteDeployment deletes a deployment
func DeleteDeployment(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		name := vars["name"]

		err := clientset.AppsV1().Deployments(namespace).Delete(r.Context(), name, metav1.DeleteOptions{})
		if err != nil {
			log.Printf("Failed to delete deployment: %v", err)
			http.Error(w, "Failed to delete deployment", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// ListServices returns all services
func ListServices(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services, err := clientset.CoreV1().Services("").List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list services: %v", err)
			http.Error(w, "Failed to list services", http.StatusInternalServerError)
			return
		}

		response := ServiceListResponse{Items: make([]Service, 0, len(services.Items))}
		for _, s := range services.Items {
			ports := make([]ServicePort, len(s.Spec.Ports))
			for i, p := range s.Spec.Ports {
				ports[i] = ServicePort{
					Port:     p.Port,
					Protocol: string(p.Protocol),
				}
			}

			response.Items = append(response.Items, Service{
				Name:      s.Name,
				Namespace: s.Namespace,
				Type:      string(s.Spec.Type),
				ClusterIP: s.Spec.ClusterIP,
				Ports:     ports,
				CreatedAt: s.CreationTimestamp.Time.Format(time.RFC3339),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetService returns service details
func GetService(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		name := vars["name"]

		service, err := clientset.CoreV1().Services(namespace).Get(r.Context(), name, metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get service: %v", err)
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}

		ports := make([]ServicePort, len(service.Spec.Ports))
		for i, p := range service.Spec.Ports {
			ports[i] = ServicePort{
				Port:     p.Port,
				Protocol: string(p.Protocol),
			}
		}

		response := Service{
			Name:      service.Name,
			Namespace: service.Namespace,
			Type:      string(service.Spec.Type),
			ClusterIP: service.Spec.ClusterIP,
			Ports:     ports,
			CreatedAt: service.CreationTimestamp.Time.Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// CreateService creates a new service
func CreateService(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Name       string `json:"name"`
			Namespace  string `json:"namespace"`
			Port       int32  `json:"port"`
			TargetPort int32  `json:"targetPort"`
		}
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      req.Name,
				Namespace: req.Namespace,
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{"app": req.Name},
				Ports: []corev1.ServicePort{
					{
						Port:       req.Port,
						TargetPort: intstr.FromInt(int(req.TargetPort)),
					},
				},
				Type: corev1.ServiceTypeClusterIP,
			},
		}

		created, err := clientset.CoreV1().Services(req.Namespace).Create(r.Context(), service, metav1.CreateOptions{})
		if err != nil {
			http.Error(w, "Failed to create service: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// ✅ JSON Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message":   "Service created successfully",
			"name":      created.Name,
			"namespace": created.Namespace,
			"clusterIP": created.Spec.ClusterIP,
		})
	}
}

func intstrFromInt(i int32) intstr.IntOrString {
	return intstr.IntOrString{Type: intstr.Int, IntVal: i}
}

// DeleteService deletes a service
func DeleteService(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		name := vars["name"]

		err := clientset.CoreV1().Services(namespace).Delete(r.Context(), name, metav1.DeleteOptions{})
		if err != nil {
			log.Printf("Failed to delete service: %v", err)
			http.Error(w, "Failed to delete service", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetPodLogs returns pod logs
func GetPodLogs(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		name := vars["name"]

		logOpts := &corev1.PodLogOptions{}
		req := clientset.CoreV1().Pods(namespace).GetLogs(name, logOpts)
		logs, err := req.Stream(r.Context())
		if err != nil {
			log.Printf("Failed to get pod logs: %v", err)
			http.Error(w, "Failed to get pod logs", http.StatusInternalServerError)
			return
		}
		defer logs.Close()

		buf := make([]byte, 4096)
		var logContent []byte
		for {
			n, err := logs.Read(buf)
			if n > 0 {
				logContent = append(logContent, buf[:n]...)
			}
			if err != nil {
				break
			}
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write(logContent)
	}
}

// Add these struct definitions to your existing structs

// Namespace represents a simplified namespace view
type Namespace struct {
	Name      string            `json:"name"`
	Status    string            `json:"status"`
	CreatedAt string            `json:"createdAt"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// Node represents a simplified node view
type Node struct {
	Name        string            `json:"name"`
	Status      string            `json:"status"`
	Version     string            `json:"version"`
	OSImage     string            `json:"osImage"`
	Capacity    map[string]string `json:"capacity"`
	Allocatable map[string]string `json:"allocatable"`
	CreatedAt   string            `json:"createdAt"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// Event represents a simplified event view
type Event struct {
	Name           string `json:"name"`
	Namespace      string `json:"namespace"`
	Reason         string `json:"reason"`
	Message        string `json:"message"`
	Type           string `json:"type"`
	InvolvedObject string `json:"involvedObject"`
	FirstTimestamp string `json:"firstTimestamp"`
	LastTimestamp  string `json:"lastTimestamp"`
	Count          int32  `json:"count"`
}

// Response types
type NamespaceListResponse struct {
	Items []Namespace `json:"items"`
}

type NodeListResponse struct {
	Items []Node `json:"items"`
}

type EventListResponse struct {
	Items []Event `json:"items"`
}

// NAMESPACE ENDPOINTS

// ListNamespaces returns all namespaces
func ListNamespaces(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		namespaces, err := clientset.CoreV1().Namespaces().List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list namespaces: %v", err)
			http.Error(w, "Failed to list namespaces", http.StatusInternalServerError)
			return
		}

		response := NamespaceListResponse{Items: make([]Namespace, 0, len(namespaces.Items))}
		for _, ns := range namespaces.Items {
			response.Items = append(response.Items, Namespace{
				Name:      ns.Name,
				Status:    string(ns.Status.Phase),
				CreatedAt: ns.CreationTimestamp.Time.Format(time.RFC3339),
				Labels:    ns.Labels,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetNamespace returns namespace details
func GetNamespace(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		namespace, err := clientset.CoreV1().Namespaces().Get(r.Context(), name, metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get namespace: %v", err)
			http.Error(w, "Namespace not found", http.StatusNotFound)
			return
		}

		response := Namespace{
			Name:      namespace.Name,
			Status:    string(namespace.Status.Phase),
			CreatedAt: namespace.CreationTimestamp.Time.Format(time.RFC3339),
			Labels:    namespace.Labels,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// CreateNamespace creates a new namespace
func CreateNamespace(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Name   string            `json:"name"`
			Labels map[string]string `json:"labels,omitempty"`
		}
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:   req.Name,
				Labels: req.Labels,
			},
		}

		created, err := clientset.CoreV1().Namespaces().Create(r.Context(), namespace, metav1.CreateOptions{})
		if err != nil {
			http.Error(w, "Failed to create namespace: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Namespace created successfully",
			"name":    created.Name,
		})
	}
}

// DeleteNamespace deletes a namespace
func DeleteNamespace(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		err := clientset.CoreV1().Namespaces().Delete(r.Context(), name, metav1.DeleteOptions{})
		if err != nil {
			log.Printf("Failed to delete namespace: %v", err)
			http.Error(w, "Failed to delete namespace", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// NODE ENDPOINTS

// ListNodes returns all nodes
func ListNodes(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodes, err := clientset.CoreV1().Nodes().List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list nodes: %v", err)
			http.Error(w, "Failed to list nodes", http.StatusInternalServerError)
			return
		}

		response := NodeListResponse{Items: make([]Node, 0, len(nodes.Items))}
		for _, node := range nodes.Items {
			status := "Unknown"
			for _, condition := range node.Status.Conditions {
				if condition.Type == corev1.NodeReady {
					if condition.Status == corev1.ConditionTrue {
						status = "Ready"
					} else {
						status = "NotReady"
					}
					break
				}
			}

			capacity := make(map[string]string)
			for k, v := range node.Status.Capacity {
				capacity[string(k)] = v.String()
			}

			allocatable := make(map[string]string)
			for k, v := range node.Status.Allocatable {
				allocatable[string(k)] = v.String()
			}

			response.Items = append(response.Items, Node{
				Name:        node.Name,
				Status:      status,
				Version:     node.Status.NodeInfo.KubeletVersion,
				OSImage:     node.Status.NodeInfo.OSImage,
				Capacity:    capacity,
				Allocatable: allocatable,
				CreatedAt:   node.CreationTimestamp.Time.Format(time.RFC3339),
				Labels:      node.Labels,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetNode returns node details
func GetNode(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		node, err := clientset.CoreV1().Nodes().Get(r.Context(), name, metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get node: %v", err)
			http.Error(w, "Node not found", http.StatusNotFound)
			return
		}

		status := "Unknown"
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady {
				if condition.Status == corev1.ConditionTrue {
					status = "Ready"
				} else {
					status = "NotReady"
				}
				break
			}
		}

		capacity := make(map[string]string)
		for k, v := range node.Status.Capacity {
			capacity[string(k)] = v.String()
		}

		allocatable := make(map[string]string)
		for k, v := range node.Status.Allocatable {
			allocatable[string(k)] = v.String()
		}

		response := Node{
			Name:        node.Name,
			Status:      status,
			Version:     node.Status.NodeInfo.KubeletVersion,
			OSImage:     node.Status.NodeInfo.OSImage,
			Capacity:    capacity,
			Allocatable: allocatable,
			CreatedAt:   node.CreationTimestamp.Time.Format(time.RFC3339),
			Labels:      node.Labels,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetNodeMetrics returns node metrics (placeholder - requires metrics-server)
// func GetNodeMetrics(clientset *kubernetes.Clientset) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		name := vars["name"]

// 		// This is a placeholder implementation
// 		// In a real scenario, you would use metrics-server client
// 		// For now, returning basic node info
// 		node, err := clientset.CoreV1().Nodes().Get(r.Context(), name, metav1.GetOptions{})
// 		if err != nil {
// 			log.Printf("Failed to get node: %v", err)
// 			http.Error(w, "Node not found", http.StatusNotFound)
// 			return
// 		}

// 		// Placeholder metrics response
// 		metrics := map[string]interface{}{
// 			"nodeName": node.Name,
// 			"message":  "Metrics endpoint requires metrics-server to be installed",
// 			"cpu":      "N/A",
// 			"memory":   "N/A",
// 			"timestamp": time.Now(),
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(metrics)
// 	}
// }

// EVENT ENDPOINTS

// ListEvents returns all events across all namespaces
func ListEvents(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		events, err := clientset.CoreV1().Events("").List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list events: %v", err)
			http.Error(w, "Failed to list events", http.StatusInternalServerError)
			return
		}

		log.Printf("Found %d events", len(events.Items))

		response := EventListResponse{Items: make([]Event, 0, len(events.Items))}
		for _, event := range events.Items {
			involvedObject := ""
			if event.InvolvedObject.Kind != "" && event.InvolvedObject.Name != "" {
				involvedObject = event.InvolvedObject.Kind + "/" + event.InvolvedObject.Name
			}

			// Handle timestamps properly - check if they're zero values
			var firstTimestamp, lastTimestamp string

			if !event.FirstTimestamp.Time.IsZero() {
				firstTimestamp = event.FirstTimestamp.Time.Format(time.RFC3339)
			}

			if !event.LastTimestamp.Time.IsZero() {
				lastTimestamp = event.LastTimestamp.Time.Format(time.RFC3339)
			} else if !event.FirstTimestamp.Time.IsZero() {
				// If LastTimestamp is empty but FirstTimestamp exists, use FirstTimestamp
				lastTimestamp = event.FirstTimestamp.Time.Format(time.RFC3339)
			}

			log.Printf("Event: %s, FirstTimestamp: %s, LastTimestamp: %s",
				event.Name, firstTimestamp, lastTimestamp)

			response.Items = append(response.Items, Event{
				Name:           event.Name,
				Namespace:      event.Namespace,
				Reason:         event.Reason,
				Message:        event.Message,
				Type:           event.Type,
				InvolvedObject: involvedObject,
				FirstTimestamp: firstTimestamp,
				LastTimestamp:  lastTimestamp,
				Count:          event.Count,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// ListEventsByNamespace returns events for a specific namespace
func ListEventsByNamespace(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]

		events, err := clientset.CoreV1().Events(namespace).List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list events for namespace %s: %v", namespace, err)
			http.Error(w, "Failed to list events", http.StatusInternalServerError)
			return
		}

		response := EventListResponse{Items: make([]Event, 0, len(events.Items))}
		for _, event := range events.Items {
			involvedObject := ""
			if event.InvolvedObject.Kind != "" && event.InvolvedObject.Name != "" {
				involvedObject = event.InvolvedObject.Kind + "/" + event.InvolvedObject.Name
			}

			// Handle timestamps properly
			var firstTimestamp, lastTimestamp string

			if !event.FirstTimestamp.Time.IsZero() {
				firstTimestamp = event.FirstTimestamp.Time.Format(time.RFC3339)
			}

			if !event.LastTimestamp.Time.IsZero() {
				lastTimestamp = event.LastTimestamp.Time.Format(time.RFC3339)
			} else if !event.FirstTimestamp.Time.IsZero() {
				lastTimestamp = event.FirstTimestamp.Time.Format(time.RFC3339)
			}

			response.Items = append(response.Items, Event{
				Name:           event.Name,
				Namespace:      event.Namespace,
				Reason:         event.Reason,
				Message:        event.Message,
				Type:           event.Type,
				InvolvedObject: involvedObject,
				FirstTimestamp: firstTimestamp,
				LastTimestamp:  lastTimestamp,
				Count:          event.Count,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

type NodeMetrics struct {
	NodeName  string             `json:"nodeName"`
	CPU       NodeResourceMetric `json:"cpu"`
	Memory    NodeResourceMetric `json:"memory"`
	Timestamp time.Time          `json:"timestamp"`
	Window    string             `json:"window"`
	Available bool               `json:"available"`
	Message   string             `json:"message,omitempty"`
}

// NodeResourceMetric represents resource usage details
type NodeResourceMetric struct {
	Value      string  `json:"value"`
	Quantity   int64   `json:"quantity"`
	Percentage float64 `json:"percentage,omitempty"`
	Unit       string  `json:"unit"`
}

// PodMetrics represents pod resource usage
type PodMetrics struct {
	PodName    string             `json:"podName"`
	Namespace  string             `json:"namespace"`
	Containers []ContainerMetrics `json:"containers"`
	CPU        NodeResourceMetric `json:"cpu"`
	Memory     NodeResourceMetric `json:"memory"`
	Timestamp  time.Time          `json:"timestamp"`
	Window     string             `json:"window"`
}

// ContainerMetrics represents container resource usage
type ContainerMetrics struct {
	Name   string             `json:"name"`
	CPU    NodeResourceMetric `json:"cpu"`
	Memory NodeResourceMetric `json:"memory"`
}

// MetricsListResponse for multiple metrics
type NodeMetricsListResponse struct {
	Items []NodeMetrics `json:"items"`
}

type PodMetricsListResponse struct {
	Items []PodMetrics `json:"items"`
}

// GetNodeMetrics returns real node metrics if metrics-server is available
func GetNodeMetrics(clientset *kubernetes.Clientset, metricsClient metricsclientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		// First check if the node exists
		node, err := clientset.CoreV1().Nodes().Get(r.Context(), name, metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get node %s: %v", name, err)
			http.Error(w, "Node not found", http.StatusNotFound)
			return
		}

		// Try to get metrics from metrics-server
		nodeMetrics, err := metricsClient.MetricsV1beta1().NodeMetricses().Get(r.Context(), name, metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get node metrics (metrics-server may not be installed): %v", err)
			// Return placeholder response
			response := NodeMetrics{
				NodeName:  name,
				Available: false,
				Message:   "Metrics server is not available. Please install metrics-server in your cluster.",
				CPU: NodeResourceMetric{
					Value: "N/A",
					Unit:  "cores",
				},
				Memory: NodeResourceMetric{
					Value: "N/A",
					Unit:  "bytes",
				},
				Timestamp: time.Now(),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// Calculate percentages based on node capacity
		cpuCapacity := node.Status.Capacity.Cpu().MilliValue()
		memoryCapacity := node.Status.Capacity.Memory().Value()

		cpuUsage := nodeMetrics.Usage.Cpu().MilliValue()
		memoryUsage := nodeMetrics.Usage.Memory().Value()

		cpuPercentage := float64(cpuUsage) / float64(cpuCapacity) * 100
		memoryPercentage := float64(memoryUsage) / float64(memoryCapacity) * 100

		// Real metrics response
		response := NodeMetrics{
			NodeName:  nodeMetrics.Name,
			Available: true,
			CPU: NodeResourceMetric{
				Value:      nodeMetrics.Usage.Cpu().String(),
				Quantity:   cpuUsage,
				Percentage: cpuPercentage,
				Unit:       "millicores",
			},
			Memory: NodeResourceMetric{
				Value:      nodeMetrics.Usage.Memory().String(),
				Quantity:   memoryUsage,
				Percentage: memoryPercentage,
				Unit:       "bytes",
			},
			Timestamp: nodeMetrics.Timestamp.Time,
			Window:    nodeMetrics.Window.Duration.String(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetNodesMetrics returns metrics for all nodes
func GetNodesMetrics(metricsClient metricsclientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodeMetricsList, err := metricsClient.MetricsV1beta1().NodeMetricses().List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list node metrics: %v", err)
			http.Error(w, "Failed to get node metrics", http.StatusInternalServerError)
			return
		}

		response := NodeMetricsListResponse{Items: make([]NodeMetrics, 0, len(nodeMetricsList.Items))}

		for _, nodeMetrics := range nodeMetricsList.Items {
			response.Items = append(response.Items, NodeMetrics{
				NodeName:  nodeMetrics.Name,
				Available: true,
				CPU: NodeResourceMetric{
					Value:    nodeMetrics.Usage.Cpu().String(),
					Quantity: nodeMetrics.Usage.Cpu().MilliValue(),
					Unit:     "millicores",
				},
				Memory: NodeResourceMetric{
					Value:    nodeMetrics.Usage.Memory().String(),
					Quantity: nodeMetrics.Usage.Memory().Value(),
					Unit:     "bytes",
				},
				Timestamp: nodeMetrics.Timestamp.Time,
				Window:    nodeMetrics.Window.Duration.String(),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetPodsMetrics returns metrics for all pods
func GetPodsMetrics(metricsClient metricsclientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		podMetricsList, err := metricsClient.MetricsV1beta1().PodMetricses("").List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list pod metrics: %v", err)
			http.Error(w, "Failed to get pod metrics", http.StatusInternalServerError)
			return
		}

		response := PodMetricsListResponse{Items: make([]PodMetrics, 0, len(podMetricsList.Items))}

		for _, podMetrics := range podMetricsList.Items {
			containers := make([]ContainerMetrics, 0, len(podMetrics.Containers))

			var totalCPU, totalMemory int64

			for _, container := range podMetrics.Containers {
				cpuUsage := container.Usage.Cpu().MilliValue()
				memoryUsage := container.Usage.Memory().Value()

				totalCPU += cpuUsage
				totalMemory += memoryUsage

				containers = append(containers, ContainerMetrics{
					Name: container.Name,
					CPU: NodeResourceMetric{
						Value:    container.Usage.Cpu().String(),
						Quantity: cpuUsage,
						Unit:     "millicores",
					},
					Memory: NodeResourceMetric{
						Value:    container.Usage.Memory().String(),
						Quantity: memoryUsage,
						Unit:     "bytes",
					},
				})
			}

			response.Items = append(response.Items, PodMetrics{
				PodName:    podMetrics.Name,
				Namespace:  podMetrics.Namespace,
				Containers: containers,
				CPU: NodeResourceMetric{
					Value:    fmt.Sprintf("%dm", totalCPU),
					Quantity: totalCPU,
					Unit:     "millicores",
				},
				Memory: NodeResourceMetric{
					Value:    fmt.Sprintf("%d", totalMemory),
					Quantity: totalMemory,
					Unit:     "bytes",
				},
				Timestamp: podMetrics.Timestamp.Time,
				Window:    podMetrics.Window.Duration.String(),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetPodMetricsByNamespace returns metrics for pods in a specific namespace
func GetPodMetricsByNamespace(metricsClient metricsclientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]

		podMetricsList, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list pod metrics for namespace %s: %v", namespace, err)
			http.Error(w, "Failed to get pod metrics", http.StatusInternalServerError)
			return
		}

		response := PodMetricsListResponse{Items: make([]PodMetrics, 0, len(podMetricsList.Items))}

		for _, podMetrics := range podMetricsList.Items {
			containers := make([]ContainerMetrics, 0, len(podMetrics.Containers))

			var totalCPU, totalMemory int64

			for _, container := range podMetrics.Containers {
				cpuUsage := container.Usage.Cpu().MilliValue()
				memoryUsage := container.Usage.Memory().Value()

				totalCPU += cpuUsage
				totalMemory += memoryUsage

				containers = append(containers, ContainerMetrics{
					Name: container.Name,
					CPU: NodeResourceMetric{
						Value:    container.Usage.Cpu().String(),
						Quantity: cpuUsage,
						Unit:     "millicores",
					},
					Memory: NodeResourceMetric{
						Value:    container.Usage.Memory().String(),
						Quantity: memoryUsage,
						Unit:     "bytes",
					},
				})
			}

			response.Items = append(response.Items, PodMetrics{
				PodName:    podMetrics.Name,
				Namespace:  podMetrics.Namespace,
				Containers: containers,
				CPU: NodeResourceMetric{
					Value:    fmt.Sprintf("%dm", totalCPU),
					Quantity: totalCPU,
					Unit:     "millicores",
				},
				Memory: NodeResourceMetric{
					Value:    fmt.Sprintf("%d", totalMemory),
					Quantity: totalMemory,
					Unit:     "bytes",
				},
				Timestamp: podMetrics.Timestamp.Time,
				Window:    podMetrics.Window.Duration.String(),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetClusterHealth provides overall cluster health status
func GetClusterHealth(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Check nodes
		nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list nodes for health check: %v", err)
			http.Error(w, "Failed to check cluster health", http.StatusInternalServerError)
			return
		}

		readyNodes := 0
		totalNodes := len(nodes.Items)

		for _, node := range nodes.Items {
			for _, condition := range node.Status.Conditions {
				if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
					readyNodes++
					break
				}
			}
		}

		// Check system pods
		systemPods, err := clientset.CoreV1().Pods("kube-system").List(ctx, metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list system pods: %v", err)
		}

		runningSystemPods := 0
		totalSystemPods := len(systemPods.Items)

		if err == nil {
			for _, pod := range systemPods.Items {
				if pod.Status.Phase == corev1.PodRunning {
					runningSystemPods++
				}
			}
		}

		// Determine overall health
		healthy := readyNodes == totalNodes && (totalSystemPods == 0 || runningSystemPods > totalSystemPods/2)
		status := "Healthy"
		if !healthy {
			status = "Degraded"
		}

		health := map[string]interface{}{
			"status":  status,
			"healthy": healthy,
			"nodes": map[string]interface{}{
				"ready": readyNodes,
				"total": totalNodes,
			},
			"systemPods": map[string]interface{}{
				"running": runningSystemPods,
				"total":   totalSystemPods,
			},
			"timestamp": time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(health)
	}
}

// GetClusterVersion returns Kubernetes cluster version info
func GetClusterVersion(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		version, err := clientset.ServerVersion()
		if err != nil {
			log.Printf("Failed to get server version: %v", err)
			http.Error(w, "Failed to get cluster version", http.StatusInternalServerError)
			return
		}

		versionInfo := map[string]interface{}{
			"major":        version.Major,
			"minor":        version.Minor,
			"gitVersion":   version.GitVersion,
			"gitCommit":    version.GitCommit,
			"gitTreeState": version.GitTreeState,
			"buildDate":    version.BuildDate,
			"goVersion":    version.GoVersion,
			"compiler":     version.Compiler,
			"platform":     version.Platform,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(versionInfo)
	}
}
