import React, { useEffect, useState } from 'react';
import { DataGrid } from '@mui/x-data-grid';
import { Box, Typography, CircularProgress, Card, CardContent, Grid, Chip } from '@mui/material';
import { toast } from 'sonner';

const MetricsPods = () => {
  const [metrics, setMetrics] = useState([]);
  const [loading, setLoading] = useState(true);

  const fetchMetrics = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem('k8s-token');
      const response = await fetch('/api/metrics/pods', { headers: { 'Authorization': `Bearer ${token}` } });
      if (response.ok) {
        const data = await response.json();
        setMetrics(data.items || []);
      } else {
        toast.error('Failed to fetch pod metrics');
      }
    } catch (error) {
      toast.error('Error fetching pod metrics: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchMetrics();
  }, []);

  // Summary cards
  const totalPods = metrics.length;
  const avgCPU = metrics.length ? (metrics.reduce((sum, m) => sum + (m.cpu.percentage || 0), 0) / metrics.length).toFixed(1) : 0;
  const avgMem = metrics.length ? (metrics.reduce((sum, m) => sum + (m.memory.percentage || 0), 0) / metrics.length).toFixed(1) : 0;

  const columns = [
    { field: 'podName', headerName: 'Pod Name', flex: 1.5 },
    { field: 'namespace', headerName: 'Namespace', flex: 1 },
    { field: 'cpu', headerName: 'CPU Usage', flex: 1, valueGetter: (params) => params.row.cpu?.value || 'N/A' },
    { field: 'cpuPercent', headerName: 'CPU %', flex: 0.7, valueGetter: (params) => params.row.cpu?.percentage ? params.row.cpu.percentage.toFixed(1) + '%' : 'N/A' },
    { field: 'memory', headerName: 'Memory Usage', flex: 1, valueGetter: (params) => params.row.memory?.value || 'N/A' },
    { field: 'memoryPercent', headerName: 'Memory %', flex: 0.7, valueGetter: (params) => params.row.memory?.percentage ? params.row.memory.percentage.toFixed(1) + '%' : 'N/A' },
    { field: 'timestamp', headerName: 'Timestamp', flex: 1, valueGetter: (params) => params.row.timestamp ? new Date(params.row.timestamp).toLocaleString() : 'N/A' },
    { field: 'window', headerName: 'Window', flex: 0.7 },
    { field: 'containers', headerName: 'Containers', flex: 2, renderCell: (params) => Array.isArray(params.value) ? params.value.map(c => <Chip key={c.name} label={c.name} size="small" sx={{ mr: 0.5 }} />) : 'N/A' },
  ];

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>Pod Metrics</Typography>
      <Grid container spacing={2} sx={{ mb: 3 }}>
        <Grid item xs={12} sm={4}>
          <Card sx={{ bgcolor: '#1976d2', color: '#fff' }}>
            <CardContent>
              <Typography variant="h6">Total Pods</Typography>
              <Typography variant="h4">{totalPods}</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={4}>
          <Card sx={{ bgcolor: '#43a047', color: '#fff' }}>
            <CardContent>
              <Typography variant="h6">Avg CPU Usage</Typography>
              <Typography variant="h4">{avgCPU}%</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={4}>
          <Card sx={{ bgcolor: '#fbc02d', color: '#fff' }}>
            <CardContent>
              <Typography variant="h6">Avg Memory Usage</Typography>
              <Typography variant="h4">{avgMem}%</Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
      {loading ? (
        <CircularProgress sx={{ display: 'block', margin: 'auto' }} />
      ) : (
        <DataGrid
          rows={metrics.map((m, i) => ({ ...m, id: m.podName + '-' + m.namespace }))}
          columns={columns}
          pageSize={10}
          rowsPerPageOptions={[10, 20, 50]}
          autoHeight
        />
      )}
    </Box>
  );
};

export default MetricsPods; 