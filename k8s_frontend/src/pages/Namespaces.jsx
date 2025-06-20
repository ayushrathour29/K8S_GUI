import React, { useState, useEffect } from 'react';
import { DataGrid } from '@mui/x-data-grid';
import { Button, CircularProgress, Box, Typography, Chip, IconButton, Tooltip, Dialog, DialogTitle, DialogContent, DialogActions, TextField } from '@mui/material';
import { Delete, Visibility } from '@mui/icons-material';
import { toast } from 'sonner';
import { formatDistanceToNow } from 'date-fns';

const getStatusChip = (status) => {
  const colorMap = {
    Active: 'success',
    Terminating: 'error',
  };
  return <Chip label={status} color={colorMap[status] || 'default'} size="small" />;
};

const NamespaceDetailsModal = ({ open, onClose, name }) => {
  const [namespace, setNamespace] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (!open || !name) return;
    setLoading(true);
    setError(null);
    setNamespace(null);
    const token = localStorage.getItem('k8s-token');
    fetch(`/api/namespaces/${name}`, { headers: { 'Authorization': `Bearer ${token}` } })
      .then(res => {
        if (!res.ok) throw new Error('Failed to fetch namespace details');
        return res.json();
      })
      .then(data => setNamespace(data))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false));
  }, [open, name]);

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Namespace Details: {name}</DialogTitle>
      <DialogContent dividers>
        {loading && <CircularProgress />}
        {error && <Typography color="error">{error}</Typography>}
        {namespace && (
          <Box>
            <Typography variant="subtitle1">Status: {namespace.status}</Typography>
            <Typography variant="subtitle1">Created At: {namespace.createdAt}</Typography>
            <Typography variant="subtitle1">Labels:</Typography>
            <pre style={{ background: '#f5f5f5', padding: 8, borderRadius: 4 }}>{JSON.stringify(namespace.labels, null, 2)}</pre>
          </Box>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Close</Button>
      </DialogActions>
    </Dialog>
  );
};

const CreateNamespaceModal = ({ open, onClose, onCreated }) => {
  const [form, setForm] = useState({
    name: '',
    labels: '{}',
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
    let labelsObj = undefined;
    if (form.labels.trim()) {
      try {
        labelsObj = JSON.parse(form.labels);
      } catch (err) {
        setError('Labels must be valid JSON');
        setLoading(false);
        return;
      }
    }
    try {
      const token = localStorage.getItem('k8s-token');
      const response = await fetch('/api/namespaces', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          name: form.name,
          labels: labelsObj,
        }),
      });
      if (!response.ok) {
        const errText = await response.text();
        throw new Error(errText || 'Failed to create namespace');
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
      <DialogTitle>Create Namespace</DialogTitle>
      <DialogContent dividers>
        <form onSubmit={handleSubmit} id="create-namespace-form">
          <Box display="flex" flexDirection="column" gap={2}>
            <TextField label="Name" name="name" value={form.name} onChange={handleChange} required fullWidth />
            <TextField label="Labels (JSON)" name="labels" value={form.labels} onChange={handleChange} fullWidth multiline minRows={2} />
            {error && <Typography color="error">{error}</Typography>}
          </Box>
        </form>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} disabled={loading}>Cancel</Button>
        <Button type="submit" form="create-namespace-form" variant="contained" disabled={loading}>Create</Button>
      </DialogActions>
    </Dialog>
  );
};

const Namespaces = () => {
  const [namespaces, setNamespaces] = useState([]);
  const [loading, setLoading] = useState(true);
  const [detailsModal, setDetailsModal] = useState({ open: false, name: '' });
  const [createModalOpen, setCreateModalOpen] = useState(false);

  const fetchNamespaces = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem('k8s-token');
      const response = await fetch('/api/namespaces', { headers: { 'Authorization': `Bearer ${token}` } });
      if (response.ok) {
        const data = await response.json();
        setNamespaces(data.items || []);
      } else {
        toast.error('Failed to fetch namespaces');
      }
    } catch (error) {
      toast.error('Error fetching namespaces: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchNamespaces();
  }, []);

  const handleDelete = async (name) => {
    if (['default', 'kube-system', 'kube-public', 'kube-node-lease'].includes(name)) {
      toast.error(`Cannot delete system namespace: ${name}`);
      return;
    }
    try {
      const token = localStorage.getItem('k8s-token');
      const response = await fetch(`/api/namespaces/${name}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
      });
      if (response.ok) {
        toast.success(`Namespace ${name} deleted successfully`);
        fetchNamespaces();
      } else {
        toast.error('Failed to delete namespace');
      }
    } catch (error) {
      toast.error('Error deleting namespace: ' + error.message);
    }
  };

  const columns = [
    { field: 'name', headerName: 'Name', flex: 2 },
    { field: 'status', headerName: 'Status', flex: 1, renderCell: (params) => getStatusChip(params.value) },
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
            <IconButton onClick={() => setDetailsModal({ open: true, name: params.row.name })}>
              <Visibility />
            </IconButton>
          </Tooltip>
          <Tooltip title="Delete Namespace">
            <span>
              <IconButton
                color="error"
                onClick={() => handleDelete(params.row.name)}
                disabled={['default', 'kube-system', 'kube-public', 'kube-node-lease'].includes(params.row.name)}
              >
                <Delete />
              </IconButton>
            </span>
          </Tooltip>
        </>
      ),
    },
  ];

  return (
    <Box sx={{ height: 'calc(100vh - 128px)', width: '100%' }}>
      <Typography variant="h4" gutterBottom>Namespaces</Typography>
      <Button onClick={() => setCreateModalOpen(true)} variant="contained" sx={{ mb: 2 }}>Create Namespace</Button>
      <Button onClick={fetchNamespaces} variant="outlined" sx={{ mb: 2, ml: 2 }}>Refresh</Button>
      {loading ? (
        <CircularProgress sx={{ display: 'block', margin: 'auto' }} />
      ) : (
        <DataGrid
          rows={namespaces.map(ns => {
            let age = 'No date';
            if (ns.createdAt) {
              try {
                const date = new Date(ns.createdAt);
                if (!isNaN(date.getTime())) {
                  age = `${formatDistanceToNow(date)} ago`;
                }
              } catch (error) {
                age = 'Invalid date';
              }
            }
            return { ...ns, id: ns.name, age };
          })}
          columns={columns}
          pageSize={10}
          rowsPerPageOptions={[10, 20, 50]}
          autoHeight
        />
      )}
      <NamespaceDetailsModal open={detailsModal.open} onClose={() => setDetailsModal({ open: false, name: '' })} name={detailsModal.name} />
      <CreateNamespaceModal open={createModalOpen} onClose={() => setCreateModalOpen(false)} onCreated={fetchNamespaces} />
    </Box>
  );
};

export default Namespaces; 