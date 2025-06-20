import React, { useState, useEffect } from 'react';
import {
  Grid, Card, CardContent, Typography, CircularProgress, Box, Paper
} from '@mui/material';
import {
  LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip,
  ResponsiveContainer, BarChart, Bar, PieChart, Pie, Cell, Legend
} from 'recharts';
import { toast } from 'sonner';
import { DataGrid } from '@mui/x-data-grid';

const nanoToMilli = (nano) => nano / 1e6;
const bytesToMebibytes = (bytes) => bytes / (1024 * 1024);

const Metrics = () => {
  const [nodeMetrics, setNodeMetrics] = useState([]);
  const [podMetrics, setPodMetrics] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchMetrics = async () => {
      setLoading(true);
      try {
        const token = localStorage.getItem('k8s-token');
        const [nodeRes, podRes] = await Promise.all([
          fetch('/api/metrics/nodes', { headers: { Authorization: `Bearer ${token}` } }),
          fetch('/api/metrics/pods', { headers: { Authorization: `Bearer ${token}` } }),
        ]);

        if (nodeRes.ok) {
          const nodeData = await nodeRes.json();
          setNodeMetrics(nodeData.items || []);
        } else {
          toast.error('Failed to fetch node metrics');
        }

        if (podRes.ok) {
          const podData = await podRes.json();
          setPodMetrics(podData.items || []);
        } else {
          toast.error('Failed to fetch pod metrics');
        }
      } catch (error) {
        toast.error('Error fetching metrics: ' + error.message);
      } finally {
        setLoading(false);
      }
    };
    fetchMetrics();
  }, []);

  // Patch: Aggregate container-level CPU/memory if top-level data is missing
  const fixedPodMetrics = podMetrics.map(pod => {
    if (!pod.cpu || !pod.memory) {
      let totalCPU = 0, totalMemory = 0;
      (pod.containers || []).forEach(c => {
        totalCPU += c.cpu?.quantity || 0;
        totalMemory += c.memory?.quantity || 0;
      });
      return {
        ...pod,
        cpu: { value: `${totalCPU}m`, quantity: totalCPU, unit: 'millicores' },
        memory: { value: `${totalMemory}`, quantity: totalMemory, unit: 'bytes' },
      };
    }
    return pod;
  });

  // Charts: use fixed pod metrics
  const podCpuData = fixedPodMetrics.map(m => ({
    name: m.podName,
    value: typeof m.cpu?.quantity === 'number' ? m.cpu.quantity : 0,
  }));

  const podMemData = fixedPodMetrics.map(m => ({
    name: m.podName,
    value: typeof m.memory?.quantity === 'number' ? bytesToMebibytes(m.memory.quantity) : 0,
  }));

  const nodeCpuData = nodeMetrics.map(m => ({
    name: m.nodeName,
    value: typeof m.cpu?.quantity === 'number' ? m.cpu.quantity : 0,
  }));

  // Summaries
  const totalNodes = nodeMetrics.length;
  const totalNodeCpu = nodeMetrics.reduce((sum, m) => sum + (typeof m.cpu?.quantity === 'number' ? m.cpu.quantity : 0), 0);
  const avgNodeCPU = totalNodes ? (totalNodeCpu / totalNodes).toFixed(1) : '0.0';

  const totalPods = fixedPodMetrics.length;
  const totalPodCpu = fixedPodMetrics.reduce((sum, m) => sum + (typeof m.cpu?.quantity === 'number' ? m.cpu.quantity : 0), 0);
  const avgPodCPU = totalPods ? (totalPodCpu / totalPods).toFixed(1) : '0.0';

  const COLORS = ['#1976d2', '#43a047', '#fbc02d', '#ff7043', '#7b1fa2', '#d32f2f', '#388e3c', '#f57c00', '#ff8042', '#82ca9d'];

  // Debug logging for API data
  console.log('Node Metrics:', nodeMetrics);
  console.log('Pod Metrics:', fixedPodMetrics);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" height="100vh">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>Metrics Dashboard</Typography>

      {/* Summary Cards */}
      <Grid container spacing={3} sx={{ mb: 3 }}>
        <Grid item xs={12} sm={6} md={3}>
          <Card sx={{ bgcolor: '#1976d2', color: '#fff' }}>
            <CardContent>
              <Typography variant="h6">Total Nodes</Typography>
              <Typography variant="h4">{totalNodes}</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card sx={{ bgcolor: '#43a047', color: '#fff' }}>
            <CardContent>
              <Typography variant="h6">Avg Node CPU</Typography>
              <Typography variant="h4">{avgNodeCPU} m</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card sx={{ bgcolor: '#fbc02d', color: '#fff' }}>
            <CardContent>
              <Typography variant="h6">Total Pods</Typography>
              <Typography variant="h4">{totalPods}</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card sx={{ bgcolor: '#ff7043', color: '#fff' }}>
            <CardContent>
              <Typography variant="h6">Avg Pod CPU</Typography>
              <Typography variant="h4">{avgPodCPU} m</Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Charts */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} md={6}>
          <Card component={Paper} elevation={4} sx={{ p: 2 }}>
            <Typography variant="h6" align="center" gutterBottom>Node CPU Usage (millicores)</Typography>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={nodeCpuData}
                  dataKey="value"
                  nameKey="name"
                  cx="50%"
                  cy="50%"
                  outerRadius={100}
                  label={({ name, value }) => `${name}: ${value} m`}
                >
                  {nodeCpuData.map((entry, index) => (
                    <Cell key={`node-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip formatter={(v) => [`${v} millicores`, 'CPU']} />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          </Card>
        </Grid>

        <Grid item xs={12} md={6}>
          <Card component={Paper} elevation={4} sx={{ p: 2 }}>
            <Typography variant="h6" align="center" gutterBottom>Top Pod Memory Usage (MiB)</Typography>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={[...podMemData].sort((a, b) => b.value - a.value).slice(0, 10)}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" angle={-45} textAnchor="end" height={80} />
                <YAxis />
                <Tooltip formatter={(v) => [`${v.toFixed(1)} MiB`, 'Memory']} />
                <Bar dataKey="value" fill="#8884d8">
                  {podMemData.slice(0, 10).map((entry, index) => (
                    <Cell key={`mem-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </Card>
        </Grid>
      </Grid>

      {/* Node Table */}
      <Typography variant="h5" sx={{ mb: 2 }}>Node Metrics</Typography>
      {nodeMetrics.length === 0 && (
        <Typography color="error" sx={{ mb: 2 }}>No node metrics data available.</Typography>
      )}
      <Box sx={{ height: 400, mb: 4 }}>
        <DataGrid
          rows={nodeMetrics.filter(m => m.nodeName)}
          columns={[
            { field: 'nodeName', headerName: 'Node Name', flex: 1 },
            {
              field: 'cpuMilli', headerName: 'CPU (m)', flex: 1,
              renderCell: (params) => {
                if (!params || !params.row || typeof params.row.cpu?.quantity !== 'number') return 'N/A';
                return `${params.row.cpu.quantity} m`;
              }
            },
            {
              field: 'memoryMiB', headerName: 'Memory (MiB)', flex: 1,
              renderCell: (params) => {
                if (!params || !params.row || typeof params.row.memory?.quantity !== 'number') return 'N/A';
                return `${(params.row.memory.quantity / (1024 * 1024)).toFixed(1)} MiB`;
              }
            },
            {
              field: 'timestamp', headerName: 'Timestamp', flex: 1,
              renderCell: (params) => {
                if (!params || !params.value) return 'N/A';
                const d = new Date(params.value);
                return isNaN(d.getTime()) ? 'N/A' : d.toLocaleString();
              }
            },
          ]}
          getRowId={(row) => row.nodeName}
          pageSize={10}
          rowsPerPageOptions={[10]}
        />
      </Box>

      {/* Pod Table */}
      <Typography variant="h5" sx={{ mb: 2 }}>Pod Metrics</Typography>
      {fixedPodMetrics.length === 0 && (
        <Typography color="error" sx={{ mb: 2 }}>No pod metrics data available.</Typography>
      )}
      <Box sx={{ height: 400 }}>
        <DataGrid
          rows={fixedPodMetrics.filter(m => m.podName && m.namespace)}
          columns={[
            { field: 'podName', headerName: 'Pod Name', flex: 1 },
            { field: 'namespace', headerName: 'Namespace', flex: 1 },
            {
              field: 'cpuMilli', headerName: 'CPU (m)', flex: 1,
              renderCell: (params) => {
                if (!params || !params.row || typeof params.row.cpu?.quantity !== 'number') return 'N/A';
                return `${params.row.cpu.quantity} m`;
              }
            },
            {
              field: 'memoryMiB', headerName: 'Memory (MiB)', flex: 1,
              renderCell: (params) => {
                if (!params || !params.row || typeof params.row.memory?.quantity !== 'number') return 'N/A';
                return `${(params.row.memory.quantity / (1024 * 1024)).toFixed(1)} MiB`;
              }
            },
            {
              field: 'timestamp', headerName: 'Timestamp', flex: 1,
              renderCell: (params) => {
                if (!params || !params.value) return 'N/A';
                const d = new Date(params.value);
                return isNaN(d.getTime()) ? 'N/A' : d.toLocaleString();
              }
            },
            {
              field: 'containers', headerName: 'Containers', flex: 2,
              renderCell: (params) => {
                if (!params || !params.row) return 'N/A';
                const containers = Array.isArray(params.row.containers) ? params.row.containers : [];
                return containers.map(c => c.name).join(', ') || 'N/A';
              }
            },
          ]}
          getRowId={(row) => `${row.podName}-${row.namespace}`}
          pageSize={10}
          rowsPerPageOptions={[10]}
        />
      </Box>
    </Box>
  );
};

export default Metrics;
