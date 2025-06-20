import React, { useState, useEffect } from 'react';
import { DataGrid } from '@mui/x-data-grid';
import { Button, CircularProgress, Box, Typography, Chip, Alert } from '@mui/material';
import { toast } from 'sonner';
import { formatDistanceToNow } from 'date-fns';

const getEventTypeChip = (type) => {
  const colorMap = {
    Normal: 'success',
    Warning: 'warning',
    Error: 'error',
  };
  return <Chip label={type} color={colorMap[type] || 'default'} size="small" />;
};

const Events = () => {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);

  const fetchEvents = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem('k8s-token');
      const response = await fetch('/api/events', { 
        headers: { 'Authorization': `Bearer ${token}` } 
      });
      
      console.log('Response status:', response.status);
      console.log('Response ok:', response.ok);
      
      if (response.ok) {
        const data = await response.json();
        console.log('Raw API data:', data);
        console.log('Events array:', data.items);
        
        // Transform the data to match DataGrid expectations
        const transformedEvents = (data.items || []).map((event, index) => {
          // Debug logging for timestamp values
          console.log(`Event ${index}:`, {
            name: event.name,
            lastTimestamp: event.lastTimestamp,
            firstTimestamp: event.firstTimestamp,
            lastTimestampType: typeof event.lastTimestamp,
            firstTimestampType: typeof event.firstTimestamp,
            lastTimestampLength: event.lastTimestamp ? event.lastTimestamp.length : 'null/undefined',
            firstTimestampLength: event.firstTimestamp ? event.firstTimestamp.length : 'null/undefined'
          });
          
          // Simplified timestamp handling
          let lastTimestamp = null;
          
          // Check lastTimestamp first (prefer it over firstTimestamp)
          if (event.lastTimestamp && event.lastTimestamp.length > 0) {
            lastTimestamp = event.lastTimestamp;
          } 
          // Fall back to firstTimestamp
          else if (event.firstTimestamp && event.firstTimestamp.length > 0) {
            lastTimestamp = event.firstTimestamp;
          }
          
          console.log(`Final timestamp for ${event.name}: "${lastTimestamp}"`);
          
          return {
            id: `${event.name || 'unknown'}-${index}`,
            name: event.name || 'N/A',
            namespace: event.namespace || 'N/A',
            reason: event.reason || 'N/A',
            message: event.message || 'N/A',
            type: event.type || 'Normal',
            involvedObject: event.involvedObject || 'N/A',
            lastTimestamp: lastTimestamp,
            count: event.count || 1
          };
        });
        
        console.log('Transformed events:', transformedEvents);
        setEvents(transformedEvents);
      } else {
        const errorText = await response.text();
        console.error('API Error:', errorText);
        toast.error('Failed to fetch events');
      }
    } catch (error) {
      console.error('Fetch error:', error);
      toast.error('Error fetching events: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEvents();
  }, []);

  const columns = [
    { field: 'name', headerName: 'Name', flex: 1.5 },
    { field: 'namespace', headerName: 'Namespace', flex: 1 },
    { field: 'reason', headerName: 'Reason', flex: 1 },
    { field: 'message', headerName: 'Message', flex: 2.5 },
    { 
      field: 'type', 
      headerName: 'Type', 
      flex: 0.8, 
      renderCell: (params) => getEventTypeChip(params.value) 
    },
    { field: 'involvedObject', headerName: 'Involved Object', flex: 1.5 },
    {
      field: 'lastTimestamp',
      headerName: 'Last Seen',
      flex: 1.2,
      renderCell: (params) => {
        const value = params.value;
        console.log('renderCell received:', value, typeof value);
        
        // Handle null, undefined, or empty values
        if (!value) {
          return 'N/A';
        }
        
        try {
          const date = new Date(value);
          if (isNaN(date.getTime())) {
            return 'Invalid date';
          }
          
          return `${formatDistanceToNow(date)} ago`;
        } catch (error) {
          console.error('Error formatting date:', error);
          return 'N/A';
        }
      },
    },
    { field: 'count', headerName: 'Count', type: 'number', flex: 0.5 },
  ];

  if (loading) {
    return (
      <Box sx={{ height: 'calc(100vh - 128px)', width: '100%' }}>
        <Typography variant="h4" gutterBottom>Events</Typography>
        <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
          <CircularProgress />
          <Typography sx={{ ml: 2 }}>Loading events...</Typography>
        </Box>
      </Box>
    );
  }

  return (
    <Box sx={{ height: 'calc(100vh - 128px)', width: '100%' }}>
      <Typography variant="h4" gutterBottom>Events</Typography>
      <Button onClick={fetchEvents} variant="contained" sx={{ mb: 2 }}>
        Refresh
      </Button>
      
      {events.length === 0 ? (
        <Alert severity="info" sx={{ mb: 2 }}>
          <Typography variant="body2">
            No events found in the cluster. This is normal if:
            <br />• No recent activity has occurred
            <br />• Events have expired (default TTL is 1 hour)  
            <br />• This is a new or quiet cluster
            <br /><br />
            Try creating a test resource: <code>kubectl run test-pod --image=nginx</code>
          </Typography>
        </Alert>
      ) : (
        <DataGrid
          rows={events}
          columns={columns}
          pageSize={25}
          rowsPerPageOptions={[10, 25, 50]}
          autoHeight
          disableSelectionOnClick
          sx={{
            '& .MuiDataGrid-row:hover': {
              backgroundColor: 'rgba(0, 0, 0, 0.04)',
            },
          }}
        />
      )}
    </Box>
  );
};

export default Events;