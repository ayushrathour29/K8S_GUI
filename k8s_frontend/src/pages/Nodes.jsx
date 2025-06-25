import React, { useState, useEffect } from 'react';
import { DataGrid } from '@mui/x-data-grid';
import { Button, CircularProgress, Box, Typography, Chip, Tooltip, IconButton, Dialog, DialogTitle, DialogContent, DialogActions } from '@mui/material';
import { Visibility } from '@mui/icons-material';
import { toast } from 'sonner';
import { formatDistanceToNow } from 'date-fns';

const getStatusChip = (status) => {
  return (
    <Chip
      label={status}
      color="success"      
      size="small"
      variant="filled"     
    />
  );
};

const NodeDetailsModal = ({ open, onClose, nodeName }) => {
  const [node, setNode] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (!open || !nodeName) return;
    setLoading(true);
    setError(null);
    setNode(null);
    const token = localStorage.getItem('k8s-token');
    fetch(`/api/nodes/${nodeName}`, { headers: { 'Authorization': `Bearer ${token}` } })
      .then(res => {
        if (!res.ok) throw new Error('Failed to fetch node details');
        return res.json();
      })
      .then(data => setNode(data))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false));
  }, [open, nodeName]);

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>Node Details: {nodeName}</DialogTitle>
      <DialogContent dividers>
        {loading && <CircularProgress />}
        {error && <Typography color="error">{error}</Typography>}
        {node && (
          <Box>
            <Typography variant="subtitle1">Status: {node.status}</Typography>
            <Typography variant="subtitle1">Kubelet Version: {node.version}</Typography>
            <Typography variant="subtitle1">OS Image: {node.osImage}</Typography>
            <Typography variant="subtitle1">Created At: {node.createdAt}</Typography>
            <Typography variant="subtitle1">Labels:</Typography>
            <pre style={{ background: '#f5f5f5', padding: 8, borderRadius: 4 }}>{JSON.stringify(node.labels, null, 2)}</pre>
            <Typography variant="subtitle1">Capacity:</Typography>
            <pre style={{ background: '#f5f5f5', padding: 8, borderRadius: 4 }}>{JSON.stringify(node.capacity, null, 2)}</pre>
            <Typography variant="subtitle1">Allocatable:</Typography>
            <pre style={{ background: '#f5f5f5', padding: 8, borderRadius: 4 }}>{JSON.stringify(node.allocatable, null, 2)}</pre>
          </Box>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Close</Button>
      </DialogActions>
    </Dialog>
  );
};

const Nodes = () => {
  const [nodes, setNodes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [selectedNode, setSelectedNode] = useState(null);

  const fetchNodes = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem('k8s-token');
      const response = await fetch('/api/nodes', { headers: { 'Authorization': `Bearer ${token}` } });
      if (response.ok) {
        const data = await response.json();
        setNodes(data.items || []);
      } else {
        toast.error('Failed to fetch nodes');
      }
    } catch (error) {
      toast.error('Error fetching nodes: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchNodes();
  }, []);

  const handleViewDetails = (nodeName) => {
    setSelectedNode(nodeName);
    setModalOpen(true);
  };

  const columns = [
    { field: 'name', headerName: 'Name', flex: 1.5 },
    { field: 'status', headerName: 'Status', flex: 1, renderCell: (params) => getStatusChip(params.value) },
    { field: 'version', headerName: 'Kubelet Version', flex: 1 },
    { field: 'osImage', headerName: 'OS Image', flex: 1.5 },
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
        <Tooltip title="View Node Details">
          <IconButton onClick={() => handleViewDetails(params.row.name)}>
            <Visibility />
          </IconButton>
        </Tooltip>
      ),
    },
  ];

  return (
    <Box sx={{ height: 'calc(100vh - 128px)', width: '100%' }}>
      <Typography variant="h4" gutterBottom>Nodes</Typography>
      <Button onClick={fetchNodes} variant="contained" sx={{ mb: 2 }}>Refresh</Button>
      {loading ? (
        <CircularProgress sx={{ display: 'block', margin: 'auto' }} />
      ) : (
        <DataGrid
          rows={nodes.map(node => {
            let age = 'No date';
            if (node.createdAt) {
              try {
                const date = new Date(node.createdAt);
                if (!isNaN(date.getTime())) {
                  age = `${formatDistanceToNow(date)} ago`;
                }
              } catch (error) {
                age = 'Invalid date';
              }
            }
            return { ...node, id: node.name, age };
          })}
          columns={columns}
          pageSize={10}
          rowsPerPageOptions={[10, 20, 50]}
          autoHeight
        />
      )}
      <NodeDetailsModal open={modalOpen} onClose={() => setModalOpen(false)} nodeName={selectedNode} />
    </Box>
  );
};

export default Nodes; 