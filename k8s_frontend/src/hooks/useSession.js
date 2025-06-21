import { useState, useEffect, useCallback } from 'react';
import { toast } from 'sonner';

const useSession = () => {
  const [token, setToken] = useState(localStorage.getItem('k8s-token'));
  const [isValidating, setIsValidating] = useState(false);

  // Check if token is expired (JWT tokens have expiration)
  const isTokenExpired = useCallback((token) => {
    if (!token) return true;
    
    try {
      // For JWT tokens, decode the payload to check expiration
      const payload = JSON.parse(atob(token.split('.')[1]));
      const currentTime = Date.now() / 1000;
      const isExpired = payload.exp < currentTime;
      console.log('Token expiration check:', { 
        expiration: new Date(payload.exp * 1000), 
        current: new Date(), 
        isExpired 
      });
      return isExpired;
    } catch (error) {
      console.error('Error parsing JWT token:', error);
      // If token is not a valid JWT, consider it expired
      return true;
    }
  }, []);

  // Validate token with server (optional validation)
  const validateToken = useCallback(async (token) => {
    if (!token) {
      console.log('No token provided for validation');
      return false;
    }
    
    try {
      console.log('Validating token with server...');
      const response = await fetch('/api/validate-token', {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });
      
      console.log('Token validation response:', response.status, response.statusText);
      
      if (response.ok) {
        const data = await response.json();
        console.log('Token validation successful:', data);
        return true;
      } else {
        console.log('Token validation failed:', response.status, response.statusText);
        return false;
      }
    } catch (error) {
      console.error('Token validation error:', error);
      // If the endpoint doesn't exist or fails, we'll rely on local validation only
      return false;
    }
  }, []);

  // Login function
  const login = useCallback((newToken) => {
    console.log('Login successful, storing token:', newToken ? 'Token received' : 'No token');
    localStorage.setItem('k8s-token', newToken);
    setToken(newToken);
  }, []);

  // Logout function
  const logout = useCallback((reason = 'Session expired') => {
    console.log('Logging out:', reason);
    localStorage.removeItem('k8s-token');
    setToken(null);
    
    if (reason !== 'manual') {
      toast.error(reason);
    }
  }, []);

  // Check session validity (simplified)
  const checkSession = useCallback(async () => {
    if (!token) {
      console.log('No token available for session check');
      return false;
    }
    
    console.log('Starting session validation...');
    setIsValidating(true);
    
    try {
      // First check if token is expired locally
      if (isTokenExpired(token)) {
        console.log('Token is expired locally');
        logout('Your session has expired. Please login again.');
        return false;
      }
      
      console.log('Token is not expired locally - session is valid');
      
      // Server validation is optional - if it fails, we still consider the session valid
      // as long as the token is not expired locally
      try {
        await validateToken(token);
        console.log('Server validation completed');
      } catch (error) {
        console.log('Server validation not available or failed, but local validation passed');
      }
      
      console.log('Session validation completed successfully');
      return true;
    } catch (error) {
      console.error('Session check error:', error);
      // Only logout if there's a critical error, not just server validation failure
      return true; // Keep session active if local validation passes
    } finally {
      setIsValidating(false);
    }
  }, [token, isTokenExpired, validateToken, logout]);

  // Check session on mount and periodically
  useEffect(() => {
    if (token) {
      console.log('Token found, starting session validation...');
      // Add a longer delay after login to ensure everything is settled
      const timeoutId = setTimeout(() => {
        checkSession();
      }, 500);
      
      // Check session every 5 minutes
      const interval = setInterval(checkSession, 5 * 60 * 1000);
      
      return () => {
        clearTimeout(timeoutId);
        clearInterval(interval);
      };
    } else {
      console.log('No token found, skipping session validation');
    }
  }, [token, checkSession]);

  // Create authenticated fetch function
  const authenticatedFetch = useCallback(async (url, options = {}) => {
    if (!token) {
      throw new Error('No authentication token available');
    }

    const response = await fetch(url, {
      ...options,
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });

    // Handle 401 Unauthorized responses
    if (response.status === 401) {
      logout('Your session has expired. Please login again.');
      throw new Error('Session expired');
    }

    return response;
  }, [token, logout]);

  return {
    token,
    isAuthenticated: !!token && !isTokenExpired(token),
    isValidating,
    login,
    logout,
    checkSession,
    authenticatedFetch,
  };
};

export default useSession; 