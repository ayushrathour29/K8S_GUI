import React, { useState, useEffect } from 'react';
import { DataGrid } from '@mui/x-data-grid';
import { Button, CircularProgress, Box, Typography, Chip, IconButton, Tooltip, Dialog, DialogTitle, DialogContent, DialogActions, TextField } from '@mui/material';
import { Delete, Edit, Visibility } from '@mui/icons-material';
import { toast } from 'sonner';
import { formatDistanceToNow } from 'date-fns';

const getHealthStatus = (params) => {
  const { availableReplicas, replicas } = params.row;
  if (availableReplicas === replicas) {
    return <Chip label="Healthy" color="success" size="small" />;
  }
  if (availableReplicas > 0) {
    return <Chip label="Degraded" color="warning" size="small" />;
  }
  return <Chip label="Unhealthy" color="error" size="small" />;
};

const DeploymentDetailsModal = ({ open, onClose, namespace, name }) => {
  const [deployment, setDeployment] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (!open || !namespace || !name) return;
    setLoading(true);
    setError(null);
    setDeployment(null);
    const token = localStorage.getItem('k8s-token');
    fetch(`/api/deployments/${namespace}/${name}`, { headers: { 'Authorization': `Bearer ${token}` } })
      .then(res => {
        if (!res.ok) throw new Error('Failed to fetch deployment details');
        return res.json();
      })
      .then(data => setDeployment(data))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false));
  }, [open, namespace, name]);

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>Deployment Details: {name}</DialogTitle>
      <DialogContent dividers>
        {loading && <CircularProgress />}
        {error && <Typography color="error">{error}</Typography>}
        {deployment && (
          <Box>
            <Typography variant="subtitle1">Namespace: {deployment.namespace}</Typography>
            <Typography variant="subtitle1">Replicas: {deployment.availableReplicas} / {deployment.replicas}</Typography>
            <Typography variant="subtitle1">Strategy: {deployment.strategy}</Typography>
            <Typography variant="subtitle1">Created At: {deployment.createdAt}</Typography>
            <Typography variant="subtitle1">Labels:</Typography>
            <pre style={{ background: '#f5f5f5', padding: 8, borderRadius: 4 }}>{JSON.stringify(deployment.labels, null, 2)}</pre>
          </Box>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Close</Button>
      </DialogActions>
    </Dialog>
  );
};

const CreateDeploymentModal = ({ open, onClose, onCreated }) => {
  const [form, setForm] = useState({
    name: '',
    namespace: '',
    image: '',
    replicas: 1,
    port: 80,
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      const token = localStorage.getItem('k8s-token');
      const response = await fetch('/api/deployments', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          ...form,
          replicas: Number(form.replicas),
          port: Number(form.port),
        }),
      });
      if (!response.ok) {
        const errText = await response.text();
        throw new Error(errText || 'Failed to create deployment');
      }
      onCreated();
      onClose();
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Create Deployment</DialogTitle>
      <DialogContent dividers>
        <form onSubmit={handleSubmit} id="create-deployment-form">
          <Box display="flex" flexDirection="column" gap={2}>
            <TextField label="Name" name="name" value={form.name} onChange={handleChange} required fullWidth />
            <TextField label="Namespace" name="namespace" value={form.namespace} onChange={handleChange} required fullWidth />
            <TextField label="Image" name="image" value={form.image} onChange={handleChange} required fullWidth />
            <TextField label="Replicas" name="replicas" type="number" value={form.replicas} onChange={handleChange} required fullWidth />
            <TextField label="Port" name="port" type="number" value={form.port} onChange={handleChange} required fullWidth />
            {error && <Typography color="error">{error}</Typography>}
          </Box>
        </form>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} disabled={loading}>Cancel</Button>
        <Button type="submit" form="create-deployment-form" variant="contained" disabled={loading}>Create</Button>
      </DialogActions>
    </Dialog>
  );
};

const EditDeploymentModal = ({ open, onClose, deployment, onUpdated }) => {
  const [form, setForm] = useState({
    image: '',
    replicas: 1,
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (deployment) {
      setForm({
        image: deployment.image || '',
        replicas: deployment.replicas || 1,
      });
    }
  }, [deployment]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      const token = localStorage.getItem('k8s-token');
      const response = await fetch(`/api/deployments/${deployment.namespace}/${deployment.name}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          image: form.image,
          replicas: Number(form.replicas),
        }),
      });
      if (!response.ok) {
        const errText = await response.text();
        throw new Error(errText || 'Failed to update deployment');
      }
      onUpdated();
      onClose();
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Edit Deployment: {deployment?.name}</DialogTitle>
      <DialogContent dividers>
        <form onSubmit={handleSubmit} id="edit-deployment-form">
          <Box display="flex" flexDirection="column" gap={2}>
            <TextField label="Image" name="image" value={form.image} onChange={handleChange} required fullWidth />
            <TextField label="Replicas" name="replicas" type="number" value={form.replicas} onChange={handleChange} required fullWidth />
            {error && <Typography color="error">{error}</Typography>}
          </Box>
        </form>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} disabled={loading}>Cancel</Button>
        <Button type="submit" form="edit-deployment-form" variant="contained" disabled={loading}>Update</Button>
      </DialogActions>
    </Dialog>
  );
};

const Deployments = () => {
  const [deployments, setDeployments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [detailsModal, setDetailsModal] = useState({ open: false, namespace: '', name: '' });
  const [createModalOpen, setCreateModalOpen] = useState(false);
  const [editModal, setEditModal] = useState({ open: false, deployment: null });

  const fetchDeployments = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem('k8s-token');
      const response = await fetch('/api/deployments', { headers: { 'Authorization': `Bearer ${token}` } });
      if (response.ok) {
        const data = await response.json();
        setDeployments(data.items || []);
      } else {
        toast.error('Failed to fetch deployments');
      }
    } catch (error) {
      toast.error('Error fetching deployments: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDeployments();
  }, []);

  const handleDelete = async (namespace, name) => {
    try {
      const token = localStorage.getItem('k8s-token');
      const response = await fetch(`/api/deployments/${namespace}/${name}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
      });
      if (response.ok) {
        toast.success(`Deployment ${name} deleted successfully`);
        fetchDeployments();
      } else {
        toast.error('Failed to delete deployment');
      }
    } catch (error) {
      toast.error('Error deleting deployment: ' + error.message);
    }
  };

  const columns = [
    { field: 'name', headerName: 'Name', flex: 1.5 },
    { field: 'namespace', headerName: 'Namespace', flex: 1 },
    {
      field: 'replicas',
      headerName: 'Replicas',
      flex: 0.7,
      renderCell: (params) => `${params.row.availableReplicas}/${params.row.replicas}`,
    },
    { field: 'strategy', headerName: 'Strategy', flex: 1 },
    { field: 'health', headerName: 'Health', flex: 1, renderCell: getHealthStatus },
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
          <Tooltip title="Edit Deployment">
            <IconButton onClick={() => setEditModal({ open: true, deployment: params.row })}>
              <Edit />
            </IconButton>
          </Tooltip>
          <Tooltip title="Delete Deployment">
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
      <Typography variant="h4" gutterBottom>Deployments</Typography>
      <Button onClick={() => setCreateModalOpen(true)} variant="contained" sx={{ mb: 2 }}>Create Deployment</Button>
      <Button onClick={fetchDeployments} variant="outlined" sx={{ mb: 2, ml: 2 }}>Refresh</Button>
      {loading ? (
        <CircularProgress sx={{ display: 'block', margin: 'auto' }} />
      ) : (
        <DataGrid
          rows={deployments.map(dep => {
            let age = 'No date';
            if (dep.createdAt) {
              try {
                const date = new Date(dep.createdAt);
                if (!isNaN(date.getTime())) {
                  age = `${formatDistanceToNow(date)} ago`;
                }
              } catch (error) {
                age = 'Invalid date';
              }
            }
            return { ...dep, id: `${dep.namespace}-${dep.name}`, age };
          })}
          columns={columns}
          pageSize={10}
          rowsPerPageOptions={[10, 20, 50]}
          autoHeight
        />
      )}
      <DeploymentDetailsModal open={detailsModal.open} onClose={() => setDetailsModal({ open: false, namespace: '', name: '' })} namespace={detailsModal.namespace} name={detailsModal.name} />
      <CreateDeploymentModal open={createModalOpen} onClose={() => setCreateModalOpen(false)} onCreated={fetchDeployments} />
      <EditDeploymentModal open={editModal.open} onClose={() => setEditModal({ open: false, deployment: null })} deployment={editModal.deployment} onUpdated={fetchDeployments} />
    </Box>
  );
};

export default Deployments; 