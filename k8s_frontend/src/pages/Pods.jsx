import React, { useState, useEffect } from 'react';
import { DataGrid } from '@mui/x-data-grid';
import { Button, CircularProgress, Box, Typography, Chip, IconButton, Tooltip, Dialog, DialogTitle, DialogContent, DialogActions } from '@mui/material';
import { Delete, Description, Terminal, Visibility } from '@mui/icons-material';
import { toast } from 'sonner';
import { formatDistanceToNow } from 'date-fns';

const getStatusChip = (params) => {
  const { status } = params.row;
  let color = 'default';
  switch (status) {
    case 'Running':
      color = 'success';
      break;
    case 'Succeeded':
      color = 'primary';
      break;
    case 'Pending':
      color = 'warning';
      break;
    case 'Failed':
      color = 'error';
      break;
    default:
      break;
  }
  return <Chip label={status} color={color} size="small" />;
};

const PodLogsModal = ({ open, onClose, namespace, name }) => {
  const [logs, setLogs] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (!open || !namespace || !name) return;
    setLoading(true);
    setError(null);
    setLogs('');
    const token = localStorage.getItem('k8s-token');
    fetch(`/api/pods/${namespace}/${name}/logs`, { headers: { 'Authorization': `Bearer ${token}` } })
      .then(res => {
        if (!res.ok) throw new Error('Failed to fetch pod logs');
        return res.text();
      })
      .then(data => setLogs(data))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false));
  }, [open, namespace, name]);

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>Logs: {name}</DialogTitle>
      <DialogContent dividers>
        {loading && <CircularProgress />}
        {error && <Typography color="error">{error}</Typography>}
        {!loading && !error && (
          <pre style={{ background: '#222', color: '#fff', padding: 12, borderRadius: 4, maxHeight: 400, overflow: 'auto' }}>{logs}</pre>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Close</Button>
      </DialogActions>
    </Dialog>
  );
};

const PodDetailsModal = ({ open, onClose, namespace, name }) => {
  const [pod, setPod] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (!open || !namespace || !name) return;
    setLoading(true);
    setError(null);
    setPod(null);
    const token = localStorage.getItem('k8s-token');
    fetch(`/api/pods/${namespace}/${name}`, { headers: { 'Authorization': `Bearer ${token}` } })
      .then(res => {
        if (!res.ok) throw new Error('Failed to fetch pod details');
        return res.json();
      })
      .then(data => setPod(data))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false));
  }, [open, namespace, name]);

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>Pod Details: {name}</DialogTitle>
      <DialogContent dividers>
        {loading && <CircularProgress />}
        {error && <Typography color="error">{error}</Typography>}
        {pod && (
          <Box>
            <Typography variant="subtitle1">Namespace: {pod.namespace}</Typography>
            <Typography variant="subtitle1">Status: {pod.status}</Typography>
            <Typography variant="subtitle1">Node: {pod.nodeName}</Typography>
            <Typography variant="subtitle1">Pod IP: {pod.podIP}</Typography>
            <Typography variant="subtitle1">Created At: {pod.createdAt}</Typography>
            <Typography variant="subtitle1">Containers:</Typography>
            <pre style={{ background: '#f5f5f5', padding: 8, borderRadius: 4 }}>{JSON.stringify(pod.containers, null, 2)}</pre>
            <Typography variant="subtitle1">Labels:</Typography>
            <pre style={{ background: '#f5f5f5', padding: 8, borderRadius: 4 }}>{JSON.stringify(pod.labels, null, 2)}</pre>
          </Box>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Close</Button>
      </DialogActions>
    </Dialog>
  );
};

const Pods = () => {
  const [pods, setPods] = useState([]);
  const [loading, setLoading] = useState(true);
  const [logsModal, setLogsModal] = useState({ open: false, namespace: '', name: '' });
  const [detailsModal, setDetailsModal] = useState({ open: false, namespace: '', name: '' });

  const fetchPods = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem('k8s-token');
      const response = await fetch('/api/pods', { headers: { 'Authorization': `Bearer ${token}` } });
      if (response.ok) {
        const data = await response.json();
        setPods(data.items || []);
      } else {
        toast.error('Failed to fetch pods');
      }
    } catch (error) {
      toast.error('Error fetching pods: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPods();
  }, []);

  const handleDelete = async (namespace, name) => {
    try {
      const token = localStorage.getItem('k8s-token');
      const response = await fetch(`/api/pods/${namespace}/${name}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
      });
      if (response.ok) {
        toast.success(`Pod ${name} deleted successfully`);
        fetchPods();
      } else {
        toast.error('Failed to delete pod');
      }
    } catch (error) {
      toast.error('Error deleting pod: ' + error.message);
    }
  };

  const columns = [
    { field: 'name', headerName: 'Name', flex: 1.5 },
    { field: 'namespace', headerName: 'Namespace', flex: 1 },
    { field: 'status', headerName: 'Status', flex: 1, renderCell: getStatusChip },
    { field: 'restartCount', headerName: 'Restarts', flex: 0.5 },
    { field: 'nodeName', headerName: 'Node', flex: 1 },
    {
      field: 'age',
      headerName: 'Age',
      flex: 1,
      sortable: false,
    },
    {
      field: 'actions',
      headerName: 'Actions',
      sortable: false,
      flex: 1,
      renderCell: (params) => (
        <>
          <Tooltip title="View Details">
            <IconButton onClick={() => setDetailsModal({ open: true, namespace: params.row.namespace, name: params.row.name })}>
              <Visibility />
            </IconButton>
          </Tooltip>
          <Tooltip title="View Logs">
            <IconButton onClick={() => setLogsModal({ open: true, namespace: params.row.namespace, name: params.row.name })}>
              <Description />
            </IconButton>
          </Tooltip>
          <Tooltip title="Delete Pod">
            <IconButton color="error" onClick={() => handleDelete(params.row.namespace, params.row.name)}>
              <Delete />
            </IconButton>
          </Tooltip>
        </>
      ),
    },
  ];

  return (
    <Box sx={{ height: 'calc(100vh - 128px)', width: '100%' }}>
      <Typography variant="h4" gutterBottom>Pods</Typography>
      <Button onClick={fetchPods} variant="contained" sx={{ mb: 2 }}>Refresh</Button>
      {loading ? (
        <CircularProgress sx={{ display: 'block', margin: 'auto' }} />
      ) : (
        <DataGrid
          rows={pods.map(pod => {
            let age = 'No date';
            if (pod.createdAt) {
              try {
                const date = new Date(pod.createdAt);
                if (!isNaN(date.getTime())) {
                  age = `${formatDistanceToNow(date)} ago`;
                }
              } catch (error) {
                age = 'Invalid date';
              }
            }
            return {
              ...pod,
              id: `${pod.namespace}-${pod.name}`,
              age: age
            };
          })}
          columns={columns}
          pageSize={10}
          rowsPerPageOptions={[10, 20, 50]}
          autoHeight
        />
      )}
      <PodLogsModal open={logsModal.open} onClose={() => setLogsModal({ open: false, namespace: '', name: '' })} namespace={logsModal.namespace} name={logsModal.name} />
      <PodDetailsModal open={detailsModal.open} onClose={() => setDetailsModal({ open: false, namespace: '', name: '' })} namespace={detailsModal.namespace} name={detailsModal.name} />
    </Box>
  );
};

export default Pods;