import React, { useState, useEffect } from 'react';
import { DataGrid } from '@mui/x-data-grid';
import { Button, CircularProgress, Box, Typography, Chip, IconButton, Tooltip, Dialog, DialogTitle, DialogContent, DialogActions, TextField } from '@mui/material';
import { Delete, Visibility } from '@mui/icons-material';
import { toast } from 'sonner';
import { formatDistanceToNow } from 'date-fns';

const getServiceTypeChip = (type) => {
  const colorMap = {
    ClusterIP: 'primary',
    NodePort: 'secondary',
    LoadBalancer: 'success',
    ExternalName: 'warning',
  };
  return <Chip label={type} color={colorMap[type] || 'default'} size="small" />;
};

const ServiceDetailsModal = ({ open, onClose, namespace, name }) => {
  const [service, setService] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (!open || !namespace || !name) return;
    setLoading(true);
    setError(null);
    setService(null);
    const token = localStorage.getItem('k8s-token');
    fetch(`/api/services/${namespace}/${name}`, { headers: { 'Authorization': `Bearer ${token}` } })
      .then(res => {
        if (!res.ok) throw new Error('Failed to fetch service details');
        return res.json();
      })
      .then(data => setService(data))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false));
  }, [open, namespace, name]);

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>Service Details: {name}</DialogTitle>
      <DialogContent dividers>
        {loading && <CircularProgress />}
        {error && <Typography color="error">{error}</Typography>}
        {service && (
          <Box>
            <Typography variant="subtitle1">Namespace: {service.namespace}</Typography>
            <Typography variant="subtitle1">Type: {service.type}</Typography>
            <Typography variant="subtitle1">Cluster IP: {service.clusterIP}</Typography>
            <Typography variant="subtitle1">Created At: {service.createdAt}</Typography>
            <Typography variant="subtitle1">Ports:</Typography>
            <pre style={{ background: '#f5f5f5', padding: 8, borderRadius: 4 }}>{JSON.stringify(service.ports, null, 2)}</pre>
          </Box>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Close</Button>
      </DialogActions>
    </Dialog>
  );
};

const CreateServiceModal = ({ open, onClose, onCreated }) => {
  const [form, setForm] = useState({
    name: '',
    namespace: '',
    port: 80,
    targetPort: 80,
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
      const response = await fetch('/api/services', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          ...form,
          port: Number(form.port),
          targetPort: Number(form.targetPort),
        }),
      });
      if (!response.ok) {
        const errText = await response.text();
        throw new Error(errText || 'Failed to create service');
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
      <DialogTitle>Create Service</DialogTitle>
      <DialogContent dividers>
        <form onSubmit={handleSubmit} id="create-service-form">
          <Box display="flex" flexDirection="column" gap={2}>
            <TextField label="Name" name="name" value={form.name} onChange={handleChange} required fullWidth />
            <TextField label="Namespace" name="namespace" value={form.namespace} onChange={handleChange} required fullWidth />
            <TextField label="Port" name="port" type="number" value={form.port} onChange={handleChange} required fullWidth />
            <TextField label="Target Port" name="targetPort" type="number" value={form.targetPort} onChange={handleChange} required fullWidth />
            {error && <Typography color="error">{error}</Typography>}
          </Box>
        </form>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} disabled={loading}>Cancel</Button>
        <Button type="submit" form="create-service-form" variant="contained" disabled={loading}>Create</Button>
      </DialogActions>
    </Dialog>
  );
};

const Services = () => {
  const [services, setServices] = useState([]);
  const [loading, setLoading] = useState(true);
  const [detailsModal, setDetailsModal] = useState({ open: false, namespace: '', name: '' });
  const [createModalOpen, setCreateModalOpen] = useState(false);

  const fetchServices = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem('k8s-token');
      const response = await fetch('/api/services', { headers: { 'Authorization': `Bearer ${token}` } });
      if (response.ok) {
        const data = await response.json();
        console.log('Raw services data:', data.items);
        setServices(data.items || []);
      } else {
        toast.error('Failed to fetch services');
      }
    } catch (error) {
      toast.error('Error fetching services: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchServices();
  }, []);

  const handleDelete = async (namespace, name) => {
    try {
      const token = localStorage.getItem('k8s-token');
      const response = await fetch(`/api/services/${namespace}/${name}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
      });
      if (response.ok) {
        toast.success(`Service ${name} deleted successfully`);
        fetchServices();
      } else {
        toast.error('Failed to delete service');
      }
    } catch (error) {
      toast.error('Error deleting service: ' + error.message);
    }
  };

  const columns = [
    { field: 'name', headerName: 'Name', flex: 1.5 },
    { field: 'namespace', headerName: 'Namespace', flex: 1 },
    { field: 'type', headerName: 'Type', flex: 1, renderCell: (params) => getServiceTypeChip(params.value) },
    { field: 'clusterIP', headerName: 'Cluster IP', flex: 1 },
    {
      field: 'ports',
      headerName: 'Ports',
      flex: 1,
      renderCell: (params) => {
        const value = params.value;
        if (Array.isArray(value) && value.length > 0) {
          return value.map(p => `${p.port}:${p.protocol}`).join(', ');
        }
        return 'N/A';
      },
    },
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
          <Tooltip title="Delete Service">
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
      <Typography variant="h4" gutterBottom>Services</Typography>
      <Button onClick={() => setCreateModalOpen(true)} variant="contained" sx={{ mb: 2 }}>Create Service</Button>
      <Button onClick={fetchServices} variant="outlined" sx={{ mb: 2, ml: 2 }}>Refresh</Button>
      {loading ? (
        <CircularProgress sx={{ display: 'block', margin: 'auto' }} />
      ) : (
        <DataGrid
          rows={services.map(svc => {
            let age = 'No date';
            if (svc.createdAt) {
              try {
                const date = new Date(svc.createdAt);
                if (!isNaN(date.getTime())) {
                  age = `${formatDistanceToNow(date)} ago`;
                }
              } catch (error) {
                age = 'Invalid date';
              }
            }
            return { ...svc, id: `${svc.namespace}-${svc.name}`, age, ports: svc.ports };
          })}
          columns={columns}
          pageSize={10}
          rowsPerPageOptions={[10, 20, 50]}
          autoHeight
        />
      )}
      <ServiceDetailsModal open={detailsModal.open} onClose={() => setDetailsModal({ open: false, namespace: '', name: '' })} namespace={detailsModal.namespace} name={detailsModal.name} />
      <CreateServiceModal open={createModalOpen} onClose={() => setCreateModalOpen(false)} onCreated={fetchServices} />
    </Box>
  );
};

export default Services; 