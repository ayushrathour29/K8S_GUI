# 🚀 Kubernetes GUI - Modern Web-Based Kubernetes Management Dashboard

<div align="center">

![Kubernetes GUI Dashboard](image.png)

[![Go Version](https://img.shields.io/badge/Go-1.20+-blue.svg)](https://golang.org/)
[![React Version](https://img.shields.io/badge/React-18.2.0-blue.svg)](https://reactjs.org/)
[![Material-UI](https://img.shields.io/badge/Material--UI-5.15.20-blue.svg)](https://mui.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

**A modern, responsive web-based dashboard for managing Kubernetes clusters with real-time metrics, authentication, and comprehensive resource management.**

[🚀 Live Demo](#live-demo) • [📋 Features](#features) • [🛠️ Installation](#installation) • [📖 Usage](#usage) • [🏗️ Architecture](#architecture)

</div>

---

## 📋 Table of Contents

- [🎯 Overview](#-overview)
- [✨ Features](#-features)
- [🖥️ Screenshots & Features](#️-screenshots--features)
- [🛠️ Installation](#️-installation)
- [📖 Usage Guide](#-usage-guide)
- [🏗️ Architecture](#️-architecture)
- [🔧 Configuration](#-configuration)
- [🤝 Contributing](#-contributing)
- [📄 License](#-license)

---

## 🎯 Overview

The Kubernetes GUI is a comprehensive web-based dashboard that provides an intuitive interface for managing Kubernetes clusters. Built with React frontend and Go backend, it offers real-time monitoring, resource management, and secure authentication.

### Key Highlights

- 🔐 **Secure Authentication** - JWT-based login system
- 📊 **Real-time Monitoring** - Live metrics and resource tracking
- 🎨 **Modern UI** - Material-UI based responsive design
- 🔍 **Global Search** - Search across all Kubernetes resources
- 📱 **Responsive Design** - Works on desktop, tablet, and mobile
- ⚡ **Performance** - Fast and efficient resource management

---

## ✨ Features

### 🔐 Authentication & Security
- JWT-based authentication system
- Secure token management with automatic expiration
- Session validation and automatic logout
- Role-based access control ready

### 📊 Dashboard & Monitoring
- Real-time cluster overview and health status
- Resource metrics (CPU, Memory, Network)
- Interactive charts and visualizations
- Performance monitoring at node and pod levels

### 🏗️ Resource Management
- **Pods**: View, create, delete, and monitor pods
- **Deployments**: Manage deployment configurations and scaling
- **Services**: Service discovery and load balancing
- **Nodes**: Cluster node monitoring and management
- **Namespaces**: Multi-tenant resource organization
- **Events**: Real-time cluster event monitoring

### 🎨 User Experience
- Responsive design for all devices
- Material-UI components for modern look
- Real-time updates and notifications
- Intuitive sidebar navigation
- Global search functionality

---

## 🖥️ Screenshots & Features

### 🔐 1. Authentication System

**Login Interface**
![Login Screen](https://via.placeholder.com/800x400/1976d2/ffffff?text=Login+Screen)

**Features:**
- Clean, modern login form
- Default credentials: `admin` / `password`
- Secure JWT token generation
- Automatic session management
- Error handling for invalid credentials

**Session Management:**
- Automatic token validation
- Session expiration handling
- Secure logout functionality
- Visual session status indicator

### 📊 2. Dashboard Overview

**Main Dashboard**
![Dashboard](https://via.placeholder.com/800x400/4caf50/ffffff?text=Dashboard+Overview)

**Features:**
- Cluster health overview
- Resource utilization charts
- Quick access to key metrics
- Real-time status updates
- Navigation to all sections

### 🔍 3. Global Search

**Search Functionality**
![Global Search](https://via.placeholder.com/800x400/ff9800/ffffff?text=Global+Search)

**Features:**
- Search across all Kubernetes resources
- Real-time search results
- Resource type indicators
- Quick navigation to resources
- Status chips for each result

**Search Capabilities:**
- Search by resource name
- Search by namespace
- Filter by resource type (Pods, Deployments, Services, Nodes)
- Click to navigate to resource details

### 📦 4. Pods Management

**Pods List View**
![Pods Management](https://via.placeholder.com/800x400/2196f3/ffffff?text=Pods+Management)

**Features:**
- Complete pod listing with status
- Pod details and logs viewing
- Pod deletion functionality
- Status indicators (Running, Pending, Failed, etc.)
- Age and restart count tracking

**Pod Operations:**
- View pod details
- Access pod logs
- Delete pods
- Monitor pod status
- View container information

### 🚀 5. Deployments Management

**Deployments View**
![Deployments](https://via.placeholder.com/800x400/9c27b0/ffffff?text=Deployments+Management)

**Features:**
- Deployment listing and status
- Replica management
- Deployment scaling
- Configuration updates
- Health monitoring

### 🌐 6. Services Management

**Services Overview**
![Services](https://via.placeholder.com/800x400/607d8b/ffffff?text=Services+Management)

**Features:**
- Service discovery
- Load balancing configuration
- Port and protocol management
- Service type indicators
- Endpoint monitoring

### 🖥️ 7. Nodes Management

**Nodes Overview**
![Nodes](https://via.placeholder.com/800x400/795548/ffffff?text=Nodes+Management)

**Features:**
- Node health monitoring
- Resource capacity tracking
- Node metrics and performance
- OS and version information
- Allocatable resources

### 📁 8. Namespaces Management

**Namespaces View**
![Namespaces](https://via.placeholder.com/800x400/ff5722/ffffff?text=Namespaces+Management)

**Features:**
- Namespace organization
- Resource isolation
- Multi-tenant support
- Namespace creation and deletion
- Resource quotas

### 📈 9. Metrics & Monitoring

**Metrics Dashboard**
![Metrics](https://via.placeholder.com/800x400/00bcd4/ffffff?text=Metrics+%26+Monitoring)

**Features:**
- Real-time performance metrics
- CPU and memory utilization
- Network statistics
- Interactive charts
- Historical data tracking

### 🔔 10. Events & Notifications

**Events Monitoring**
![Events](https://via.placeholder.com/800x400/8bc34a/ffffff?text=Events+%26+Notifications)

**Features:**
- Real-time event monitoring
- Event filtering and search
- Notification system
- Event history tracking
- Alert management

---

## 🛠️ Installation

### Prerequisites

- **Go 1.20+** - Backend runtime
- **Node.js 16+** - Frontend development
- **Kubernetes Cluster** - Target cluster for management
- **kubectl** - Kubernetes command-line tool (optional)

### Quick Start

1. **Clone the Repository**

   ```bash
   git clone https://github.com/yourusername/kubernetes-gui.git
   cd kubernetes-gui
   ```

2. **Backend Setup**

   ```bash
   cd server

   # Create environment file
   cat > .env << EOF
   PORT=8081
   JWT_SECRET=your_super_secret_jwt_key_here
   EOF

   # Install dependencies
   go mod tidy

   # Run the backend
   go run cmd/api/main.go
   ```

3. **Frontend Setup**

   ```bash
   cd k8s_frontend

   # Install dependencies
   npm install

   # Start development server
   npm run dev
   ```

4. **Access the Application**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8081
   - Default credentials: `admin` / `password`

### Production Deployment

1. **Build Frontend**

   ```bash
   cd k8s_frontend
   npm run build
   ```

2. **Build Backend**

   ```bash
   cd server
   go build -o k8s-gui cmd/api/main.go
   ```

3. **Environment Configuration**
   ```bash
   # Production .env
   PORT=8081
   JWT_SECRET=your_production_jwt_secret_key
   ```

---

## 📖 Usage Guide

### 🔐 Authentication

1. **Login Process**
   - Navigate to the application URL
   - Enter credentials: `admin` / `password`
   - Click "Sign In" to authenticate
   - JWT token is automatically stored and managed

2. **Session Management**
   - Session automatically validates every 5 minutes
   - Automatic logout on session expiration
   - Visual session status indicator in header
   - Manual logout available via header button

### 📊 Dashboard Navigation

1. **Main Dashboard**
   - View cluster overview and health
   - Access quick metrics and statistics
   - Navigate to different resource sections

2. **Sidebar Navigation**
   - **Dashboard**: Main overview page
   - **Nodes**: Cluster node management
   - **Pods**: Pod monitoring and management
   - **Deployments**: Application deployment management
   - **Services**: Service networking
   - **Namespaces**: Resource organization
   - **Events**: Cluster event monitoring
   - **Metrics**: Performance metrics

### 🔍 Global Search

1. **Using Search**
   - Click the search bar in the header
   - Type resource name or namespace
   - View real-time search results
   - Click on result to navigate

2. **Search Features**
   - Searches across all resource types
   - Shows resource status and type
   - Provides quick navigation
   - Filters results by type

### 📦 Resource Management

1. **Pods Management**
   - View all pods in the cluster
   - Check pod status and health
   - Access pod logs and details
   - Delete pods when needed

2. **Deployments Management**
   - Monitor deployment status
   - Scale deployments up/down
   - Update deployment configurations
   - View deployment history

3. **Services Management**
   - Configure service networking
   - Manage load balancing
   - Monitor service endpoints
   - Update service configurations

4. **Nodes Management**
   - Monitor node health
   - View resource utilization
   - Check node capacity
   - Monitor node performance

### 📈 Monitoring & Metrics

1. **Real-time Metrics**
   - Monitor CPU and memory usage
   - Track network statistics
   - View performance trends
   - Analyze resource utilization

2. **Event Monitoring**
   - Track cluster events
   - Monitor resource changes
   - View error notifications
   - Analyze event patterns

---

## 🏗️ Architecture

### Frontend Architecture

```
k8s_frontend/
├── src/
│   ├── components/          # Reusable UI components
│   │   ├── Header.jsx      # Navigation and search
│   │   ├── Sidebar.jsx     # Main navigation
│   │   └── Login.jsx       # Authentication
│   ├── pages/              # Main application pages
│   │   ├── Dashboard.jsx   # Main dashboard
│   │   ├── Pods.jsx        # Pods management
│   │   ├── Deployments.jsx # Deployments management
│   │   └── ...             # Other resource pages
│   ├── hooks/              # Custom React hooks
│   │   └── useSession.js   # Session management
│   └── App.jsx             # Main application component
```

### Backend Architecture

```
server/
├── cmd/api/
│   └── main.go             # Application entry point
├── internal/
│   ├── api/                # API handlers
│   │   └── handlers.go     # Resource management endpoints
│   ├── auth/               # Authentication
│   │   └── auth.go         # JWT authentication
│   ├── k8s/                # Kubernetes client
│   │   └── client.go       # K8s client initialization
│   └── server/             # Server configuration
│       └── router.go       # Route definitions
```

### Technology Stack

- **Frontend**: React 18, Material-UI, Vite
- **Backend**: Go, Gorilla Mux, JWT
- **Database**: Kubernetes API (stateless)
- **Authentication**: JWT tokens
- **Styling**: Material-UI components

---

## 🔧 Configuration

### Environment Variables

```bash
# Server Configuration
PORT=8081                    # Server port
JWT_SECRET=your_secret_key   # JWT signing secret

# Kubernetes Configuration
KUBECONFIG=/path/to/kubeconfig  # K8s config path (optional)
```

### Frontend Configuration

```javascript
// vite.config.js
export default defineConfig({
  server: {
    proxy: {
      '/api': 'http://localhost:8081'
    }
  }
})
```

### Security Configuration

1. **JWT Secret**: Use a strong, unique secret key
2. **HTTPS**: Enable HTTPS in production
3. **CORS**: Configure CORS settings for your domain
4. **Rate Limiting**: Implement rate limiting for API endpoints

---

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Commit your changes**: `git commit -m 'Add amazing feature'`
4. **Push to the branch**: `git push origin feature/amazing-feature`
5. **Open a Pull Request**

### Development Guidelines

- Follow Go and React best practices
- Add tests for new features
- Update documentation for changes
- Ensure code quality and formatting

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- [Kubernetes](https://kubernetes.io/) - Container orchestration platform
- [React](https://reactjs.org/) - Frontend framework
- [Material-UI](https://mui.com/) - UI component library
- [Go](https://golang.org/) - Backend programming language

---

<div align="center">

**Made with ❤️ for the Kubernetes community**

[⭐ Star this repo](https://github.com/yourusername/kubernetes-gui) • [🐛 Report issues](https://github.com/yourusername/kubernetes-gui/issues) • [📖 Documentation](https://github.com/yourusername/kubernetes-gui/wiki)

</div> 