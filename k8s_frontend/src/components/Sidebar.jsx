import React from 'react';
import { Box, Drawer, List, ListItem, ListItemButton, ListItemIcon, ListItemText, Toolbar, Typography } from '@mui/material';
import { Link, useLocation } from 'react-router-dom';
import { Dashboard, Dns, Apps, Layers, NetworkWifi, Folder, Event, BarChart, ChevronLeft, Storage, Assessment } from '@mui/icons-material';

const drawerWidth = 240;

const menuItems = [
  { path: '/', label: 'Dashboard', icon: <Dashboard /> },
  { path: '/nodes', label: 'Nodes', icon: <Dns /> },
  { path: '/pods', label: 'Pods', icon: <Apps /> },
  { path: '/deployments', label: 'Deployments', icon: <Layers /> },
  { path: '/services', label: 'Services', icon: <NetworkWifi /> },
  { path: '/namespaces', label: 'Namespaces', icon: <Folder /> },
  { path: '/events', label: 'Events', icon: <Event /> },
  { path: '/metrics', label: 'Metrics', icon: <BarChart /> },
];

const Sidebar = ({ open }) => {
  const location = useLocation();

  const drawerContent = (
    <div>
      <Toolbar sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', px: [1] }}>
        <Typography variant="h6" noWrap>K8s Dashboard</Typography>
      </Toolbar>
      <List>
        {menuItems.map((item) => {
          const isActive = location.pathname === item.path;
          return (
            <ListItem key={item.path} disablePadding component={Link} to={item.path} sx={{ color: 'inherit', textDecoration: 'none' }}>
              <ListItemButton selected={isActive} sx={{
                '&.Mui-selected': {
                  backgroundColor: 'primary.main',
                  color: 'primary.contrastText',
                  ':hover': {
                    backgroundColor: 'primary.dark',
                  },
                },
              }}>
                <ListItemIcon sx={{ color: isActive ? 'primary.contrastText' : 'inherit' }}>
                  {item.icon}
                </ListItemIcon>
                <ListItemText primary={item.label} />
              </ListItemButton>
            </ListItem>
          );
        })}
      </List>
    </div>
  );

  return (
    <Box component="nav" sx={{ width: { sm: open ? drawerWidth : 0 }, flexShrink: { sm: 0 } }} aria-label="mailbox folders">
      <Drawer
        variant="persistent"
        open={open}
        sx={{
          '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
        }}
      >
        {drawerContent}
      </Drawer>
    </Box>
  );
};

export default Sidebar; 