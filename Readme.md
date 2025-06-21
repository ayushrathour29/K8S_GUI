
The Dashboard of Kubernetes Task Manager
![alt text](image.png)




## The file structure for backend:

1. cmd/backend/main.go: main entry point
2. internal/k8s/client.go: Kubernetes client initialization
3. internal/api/handlers.go: HTTP handlers for auth, pods, deployments, services, clusters
4. internal/auth/auth.go: authentication functions
5. internal/server/router.go: router and middleware setup