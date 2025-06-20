import React, { useEffect, useState } from 'react';
import { DataGrid } from '@mui/x-data-grid';
import { Box, Typography, CircularProgress, Card, CardContent, Grid } from '@mui/material';
import { toast } from 'sonner';

const MetricsNodes = () => {
  const [metrics, setMetrics] = useState([]);
  const [loading, setLoading] = useState(true);

  const fetchMetrics = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem('k8s-token');
      const response = await fetch('/api/metrics/nodes', { headers: { 'Authorization': `Bearer ${token}` } });
      if (response.ok) {
        const data = await response.json();
        setMetrics(data.items || []);
      } else {
        toast.error('Failed to fetch node metrics');
      }
    } catch (error) {
      toast.error('Error fetching node metrics: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchMetrics();
  }, []);

  // Summary cards
  const totalNodes = metrics.length;
  const avgCPU = metrics.length ? (metrics.reduce((sum, m) => sum + (m.cpu.percentage || 0), 0) / metrics.length).toFixed(1) : 0;
  const avgMem = metrics.length ? (metrics.reduce((sum, m) => sum + (m.memory.percentage || 0), 0) / metrics.length).toFixed(1) : 0;

  const columns = [
    { field: 'nodeName', headerName: 'Node Name', flex: 1.5 },
    { field: 'cpu', headerName: 'CPU Usage', flex: 1, valueGetter: (params) => params?.row?.cpu && typeof params.row.cpu.value !== 'undefined' ? params.row.cpu.value : 'N/A' },
    { field: 'cpuPercent', headerName: 'CPU %', flex: 0.7, valueGetter: (params) => params?.row?.cpu && typeof params.row.cpu.percentage === 'number' ? params.row.cpu.percentage.toFixed(1) + '%' : 'N/A' },
    { field: 'memory', headerName: 'Memory Usage', flex: 1, valueGetter: (params) => params?.row?.memory && typeof params.row.memory.value !== 'undefined' ? params.row.memory.value : 'N/A' },
    { field: 'memoryPercent', headerName: 'Memory %', flex: 0.7, valueGetter: (params) => params?.row?.memory && typeof params.row.memory.percentage === 'number' ? params.row.memory.percentage.toFixed(1) + '%' : 'N/A' },
    { field: 'timestamp', headerName: 'Timestamp', flex: 1, valueGetter: (params) => params?.row?.timestamp ? new Date(params.row.timestamp).toLocaleString() : 'N/A' },
    { field: 'status', headerName: 'Status', flex: 1, valueGetter: (params) => typeof params?.row?.available === 'boolean' ? (params.row.available ? 'Available' : 'Unavailable') : 'N/A' },
    { field: 'message', headerName: 'Message', flex: 1.5, valueGetter: (params) => params?.row?.message || '' },
  ];

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>Node Metrics</Typography>
      <Grid container spacing={2} sx={{ mb: 3 }}>
        <Grid item xs={12} sm={4}>
          <Card sx={{ bgcolor: '#1976d2', color: '#fff' }}>
            <CardContent>
              <Typography variant="h6">Total Nodes</Typography>
              <Typography variant="h4">{totalNodes}</Typography>
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
          rows={metrics.map((m, i) => ({ ...m, id: m.nodeName || i }))}
          columns={columns}
          pageSize={10}
          rowsPerPageOptions={[10, 20, 50]}
          autoHeight
        />
      )}
    </Box>
  );
};

export default MetricsNodes; 