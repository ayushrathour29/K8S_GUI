import React, { useState } from 'react';
import { Button, TextField, Container, Typography, Box, Paper, Avatar, CircularProgress } from '@mui/material';
import { LockOutlined } from '@mui/icons-material';
import { toast } from 'sonner';

const Login = ({ onLogin }) => {
  const [credentials, setCredentials] = useState({ username: '', password: '' });
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      console.log('Attempting login with credentials:', { username: credentials.username, password: '***' });
      const response = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(credentials)
      });

      console.log('Login response status:', response.status);
      
      if (response.ok) {
        const data = await response.json();
        console.log('Login successful, token received:', data.token ? 'Token present' : 'No token');
        console.log('Token length:', data.token ? data.token.length : 0);
        
        if (data.token) {
          onLogin(data.token);
          toast.success('Login successful!');
        } else {
          console.error('No token in response');
          toast.error('Login failed: No token received');
        }
      } else {
        console.log('Login failed with status:', response.status);
        const errorText = await response.text();
        console.log('Error response:', errorText);
        toast.error('Invalid credentials');
      }
    } catch (error) {
      console.error('Login error:', error);
      toast.error('Connection failed. Please check if the server is running.');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e) => {
    setCredentials({ ...credentials, [e.target.name]: e.target.value });
  };

  return (
    <Container component="main" maxWidth="xs">
      <Paper elevation={6} sx={{
        marginTop: 8,
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        padding: 4,
      }}>
        <Avatar sx={{ m: 1, bgcolor: 'secondary.main' }}>
          <LockOutlined />
        </Avatar>
        <Typography component="h1" variant="h5">
          Sign in
        </Typography>
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 1 }}>
          <TextField
            margin="normal"
            required
            fullWidth
            id="username"
            label="Username"
            name="username"
            autoComplete="username"
            autoFocus
            value={credentials.username}
            onChange={handleChange}
          />
          <TextField
            margin="normal"
            required
            fullWidth
            name="password"
            label="Password"
            type="password"
            id="password"
            autoComplete="current-password"
            value={credentials.password}
            onChange={handleChange}
          />
          <Button
            type="submit"
            fullWidth
            variant="contained"
            sx={{ mt: 3, mb: 2 }}
            disabled={loading}
          >
            {loading ? <CircularProgress size={24} /> : 'Sign In'}
          </Button>
          <Typography variant="body2" color="text.secondary" align="center">
            Default: admin / password
          </Typography>
        </Box>
      </Paper>
    </Container>
  );
};

export default Login; 