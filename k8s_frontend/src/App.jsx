import React, { useState } from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import { createTheme, ThemeProvider, CssBaseline, Box, Toolbar, CircularProgress } from '@mui/material';
import Sidebar from './components/Sidebar';
import Header from './components/Header';
import Login from './components/Login';
import Dashboard from './pages/Dashboard';
import Nodes from './pages/Nodes';
import Pods from './pages/Pods';
import Deployments from './pages/Deployments';
import Services from './pages/Services';
import Namespaces from './pages/Namespaces';
import Events from './pages/Events';
import Metrics from './pages/Metrics';
import { Toaster } from 'sonner';
import useSession from './hooks/useSession';

const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
    background: {
      default: '#f4f6f8',
      paper: '#ffffff',
    },
  },
  typography: {
    fontFamily: 'Roboto, sans-serif',
  },
});

const App = () => {
  const { token, isAuthenticated, isValidating, login, logout } = useSession();
  const [sidebarOpen, setSidebarOpen] = useState(true);

  const handleLogout = () => {
    logout('manual');
  };

  const toggleSidebar = () => {
    setSidebarOpen(!sidebarOpen);
  };

  // Show loading while validating session
  if (isValidating) {
    return (
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <Toaster richColors position="top-right" />
        <Box 
          sx={{ 
            display: 'flex', 
            justifyContent: 'center', 
            alignItems: 'center', 
            height: '100vh' 
          }}
        >
          <CircularProgress size={60} />
        </Box>
      </ThemeProvider>
    );
  }

  // Show login if not authenticated
  if (!isAuthenticated) {
    return (
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <Toaster richColors position="top-right" />
        <Login onLogin={login} />
      </ThemeProvider>
    );
  }

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Toaster richColors position="top-right" />
      <Router future={{ v7_startTransition: true, v7_relativeSplatPath: true }}>
        <Box sx={{ display: 'flex' }}>
          <Header onLogout={handleLogout} onToggleSidebar={toggleSidebar} sidebarOpen={sidebarOpen} />
          <Sidebar open={sidebarOpen} />
          <Box component="main" sx={{ flexGrow: 1, p: 3, width: { sm: `calc(100% - 240px)` } }}>
            <Toolbar />
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/nodes" element={<Nodes />} />
              <Route path="/pods" element={<Pods />} />
              <Route path="/deployments" element={<Deployments />} />
              <Route path="/services" element={<Services />} />
              <Route path="/namespaces" element={<Namespaces />} />
              <Route path="/events" element={<Events />} />
              <Route path="/metrics" element={<Metrics />} />
              <Route path="*" element={<Navigate to="/" />} />
            </Routes>
          </Box>
        </Box>
      </Router>
    </ThemeProvider>
  );
};

export default App; 