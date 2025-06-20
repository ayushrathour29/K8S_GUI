import React, { useState, useEffect } from 'react';
import { Grid, Card, CardContent, Typography, CircularProgress, Box, Paper } from '@mui/material';
import { Dns, Apps, Layers, NetworkWifi, Folder, Event } from '@mui/icons-material';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import { toast } from 'sonner';

const StatCard = ({ title, value, icon, color }) => (
  <Card component={Paper} elevation={4}>
    <CardContent sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
      <Box>
        <Typography color="text.secondary">{title}</Typography>
        <Typography variant="h4">{value}</Typography>
      </Box>
      <Box sx={{ backgroundColor: color, borderRadius: '50%', padding: 2, display: 'flex' }}>
        {icon}
      </Box>
    </CardContent>
  </Card>
);

const Dashboard = () => {
  const [stats, setStats] = useState({
    pods: 0, deployments: 0, services: 0, nodes: 0, namespaces: 0, events: 0
  });
  const [loading, setLoading] = useState(true);

  const getAuthHeaders = () => ({ 'Authorization': `Bearer ${localStorage.getItem('k8s-token')}` });

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        const endpoints = ['pods', 'deployments', 'services', 'nodes', 'namespaces', 'events'];
        const requests = endpoints.map(ep => fetch(`/api/${ep}`, { headers: getAuthHeaders() }));
        const responses = await Promise.all(requests);
        const data = await Promise.all(responses.map(res => res.ok ? res.json() : { items: [] }));
        
        setStats({
          pods: data[0].items?.length || 0,
          deployments: data[1].items?.length || 0,
          services: data[2].items?.length || 0,
          nodes: data[3].items?.length || 0,
          namespaces: data[4].items?.length || 0,
          events: data[5].items?.length || 0,
        });
      } catch (error) {
        toast.error('Failed to load dashboard data');
      } finally {
        setLoading(false);
      }
    };
    fetchDashboardData();
  }, []);

  const statCards = [
    { title: 'Nodes', value: stats.nodes, icon: <Dns sx={{ color: '#fff' }} />, color: '#1976d2' },
    { title: 'Pods', value: stats.pods, icon: <Apps sx={{ color: '#fff' }} />, color: '#388e3c' },
    { title: 'Deployments', value: stats.deployments, icon: <Layers sx={{ color: '#fff' }} />, color: '#f57c00' },
    { title: 'Services', value: stats.services, icon: <NetworkWifi sx={{ color: '#fff' }} />, color: '#d32f2f' },
    { title: 'Namespaces', value: stats.namespaces, icon: <Folder sx={{ color: '#fff' }} />, color: '#7b1fa2' },
    { title: 'Events', value: stats.events, icon: <Event sx={{ color: '#fff' }} />, color: '#00796b' },
  ];

  const chartData = [
    { name: 'Pods', value: stats.pods },
    { name: 'Deployments', value: stats.deployments },
    { name: 'Services', value: stats.services },
  ];

  const COLORS = ['#388e3c', '#f57c00', '#d32f2f'];

  if (loading) {
    return <CircularProgress sx={{ display: 'block', margin: 'auto' }} />;
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Dashboard</Typography>
      <Grid container spacing={3}>
        {statCards.map(card => (
          <Grid item xs={12} sm={6} md={4} key={card.title}>
            <StatCard {...card} />
          </Grid>
        ))}
        <Grid item xs={12} lg={6}>
          <Card component={Paper} elevation={4}>
            <CardContent>
              <Typography variant="h6">Resource Distribution</Typography>
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="name" />
                  <YAxis />
                  <Tooltip />
                  <Bar dataKey="value" fill="#1976d2" />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} lg={6}>
          <Card component={Paper} elevation={4}>
            <CardContent>
              <Typography variant="h6">Resource Overview</Typography>
              <ResponsiveContainer width="100%" height={300}>
                <PieChart>
                  <Pie data={chartData} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={100} label>
                    {chartData.map((entry, index) => <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />)}
                  </Pie>
                  <Tooltip />
                </PieChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};

export default Dashboard; 