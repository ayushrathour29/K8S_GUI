import React, { useState, useEffect } from 'react';
import { 
  AppBar, 
  Toolbar, 
  IconButton, 
  Typography, 
  InputBase, 
  Badge, 
  Avatar,
  Menu,
  MenuItem,
  ListItemText,
  ListItemIcon,
  Divider,
  Box,
  Chip,
  CircularProgress,
  Typography as MuiTypography,
  List,
  ListItem,
  ListItemButton,
  Paper,
  Popper,
  ClickAwayListener,
  Tooltip
} from '@mui/material';
import { 
  Menu as MenuIcon, 
  Search, 
  Notifications, 
  Logout,
  Warning,
  Info,
  Error,
  CheckCircle,
  Schedule,
  Storage,
  Cloud,
  Settings,
  ViewList,
  Security
} from '@mui/icons-material';
import { styled, alpha } from '@mui/material/styles';
import { toast } from 'sonner';
import { useNavigate } from 'react-router-dom';
import useSession from '../hooks/useSession';

const drawerWidth = 240;

const SearchBar = styled('div')(({ theme }) => ({
  position: 'relative',
  borderRadius: theme.shape.borderRadius,
  backgroundColor: alpha(theme.palette.common.white, 0.15),
  '&:hover': {
    backgroundColor: alpha(theme.palette.common.white, 0.25),
  },
  marginRight: theme.spacing(2),
  marginLeft: 0,
  width: '100%',
  [theme.breakpoints.up('sm')]: {
    marginLeft: theme.spacing(3),
    width: 'auto',
  },
}));

const SearchIconWrapper = styled('div')(({ theme }) => ({
  padding: theme.spacing(0, 2),
  height: '100%',
  position: 'absolute',
  pointerEvents: 'none',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
}));

const StyledInputBase = styled(InputBase)(({ theme }) => ({
  color: 'inherit',
  '& .MuiInputBase-input': {
    padding: theme.spacing(1, 1, 1, 0),
    paddingLeft: `calc(1em + ${theme.spacing(4)})`,
    transition: theme.transitions.create('width'),
    width: '100%',
    [theme.breakpoints.up('md')]: {
      width: '20ch',
    },
  },
}));

const NotificationItem = ({ notification }) => {
  const getIcon = (type) => {
    switch (type) {
      case 'Warning':
        return <Warning color="warning" />;
      case 'Error':
        return <Error color="error" />;
      case 'Normal':
        return <Info color="info" />;
      default:
        return <CheckCircle color="success" />;
    }
  };

  const getTimeAgo = (timestamp) => {
    const now = new Date();
    const eventTime = new Date(timestamp);
    const diffInMinutes = Math.floor((now - eventTime) / (1000 * 60));
    
    if (diffInMinutes < 1) return 'Just now';
    if (diffInMinutes < 60) return `${diffInMinutes}m ago`;
    if (diffInMinutes < 1440) return `${Math.floor(diffInMinutes / 60)}h ago`;
    return `${Math.floor(diffInMinutes / 1440)}d ago`;
  };

  return (
    <MenuItem sx={{ minWidth: 300, maxWidth: 400 }}>
      <ListItemIcon>
        {getIcon(notification.type)}
      </ListItemIcon>
      <ListItemText
        primary={notification.reason}
        secondary={
          <Box>
            <Typography variant="body2" color="text.secondary">
              {notification.message}
            </Typography>
            <Box sx={{ display: 'flex', alignItems: 'center', mt: 0.5 }}>
              <Schedule sx={{ fontSize: 12, mr: 0.5 }} />
              <Typography variant="caption" color="text.secondary">
                {getTimeAgo(notification.lastTimestamp)}
              </Typography>
              <Chip 
                label={notification.namespace} 
                size="small" 
                sx={{ ml: 1, height: 16 }}
              />
            </Box>
          </Box>
        }
      />
    </MenuItem>
  );
};

const Header = ({ onLogout, onToggleSidebar, sidebarOpen }) => {
  const [notifications, setNotifications] = useState([]);
  const [loading, setLoading] = useState(false);
  const [anchorEl, setAnchorEl] = useState(null);
  const [notificationCount, setNotificationCount] = useState(0);
  
  // Search state
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [searchLoading, setSearchLoading] = useState(false);
  const [searchAnchorEl, setSearchAnchorEl] = useState(null);
  
  const navigate = useNavigate();
  const { authenticatedFetch, isValidating, isAuthenticated } = useSession();

  // Search functionality
  const performSearch = async (query) => {
    if (!query.trim()) {
      setSearchResults([]);
      return;
    }

    setSearchLoading(true);
    try {
      const results = [];

      // Search pods
      try {
        const podsResponse = await authenticatedFetch('/api/pods');
        if (podsResponse.ok) {
          const podsData = await podsResponse.json();
          const matchingPods = podsData.items?.filter(pod => 
            pod.name.toLowerCase().includes(query.toLowerCase()) ||
            pod.namespace.toLowerCase().includes(query.toLowerCase())
          ) || [];
          
          matchingPods.forEach(pod => {
            results.push({
              id: `pod-${pod.namespace}-${pod.name}`,
              name: pod.name,
              namespace: pod.namespace,
              type: 'Pod',
              status: pod.status,
              icon: <Storage />,
              path: '/pods'
            });
          });
        }
      } catch (error) {
        console.error('Error searching pods:', error);
      }

      // Search deployments
      try {
        const deploymentsResponse = await authenticatedFetch('/api/deployments');
        if (deploymentsResponse.ok) {
          const deploymentsData = await deploymentsResponse.json();
          const matchingDeployments = deploymentsData.items?.filter(deployment => 
            deployment.name.toLowerCase().includes(query.toLowerCase()) ||
            deployment.namespace.toLowerCase().includes(query.toLowerCase())
          ) || [];
          
          matchingDeployments.forEach(deployment => {
            results.push({
              id: `deployment-${deployment.namespace}-${deployment.name}`,
              name: deployment.name,
              namespace: deployment.namespace,
              type: 'Deployment',
              status: deployment.status,
              icon: <Cloud />,
              path: '/deployments'
            });
          });
        }
      } catch (error) {
        console.error('Error searching deployments:', error);
      }

      // Search services
      try {
        const servicesResponse = await authenticatedFetch('/api/services');
        if (servicesResponse.ok) {
          const servicesData = await servicesResponse.json();
          const matchingServices = servicesData.items?.filter(service => 
            service.name.toLowerCase().includes(query.toLowerCase()) ||
            service.namespace.toLowerCase().includes(query.toLowerCase())
          ) || [];
          
          matchingServices.forEach(service => {
            results.push({
              id: `service-${service.namespace}-${service.name}`,
              name: service.name,
              namespace: service.namespace,
              type: 'Service',
              status: service.status,
              icon: <Settings />,
              path: '/services'
            });
          });
        }
      } catch (error) {
        console.error('Error searching services:', error);
      }

      // Search nodes
      try {
        const nodesResponse = await authenticatedFetch('/api/nodes');
        if (nodesResponse.ok) {
          const nodesData = await nodesResponse.json();
          const matchingNodes = nodesData.items?.filter(node => 
            node.name.toLowerCase().includes(query.toLowerCase())
          ) || [];
          
          matchingNodes.forEach(node => {
            results.push({
              id: `node-${node.name}`,
              name: node.name,
              namespace: 'N/A',
              type: 'Node',
              status: node.status,
              icon: <ViewList />,
              path: '/nodes'
            });
          });
        }
      } catch (error) {
        console.error('Error searching nodes:', error);
      }

      setSearchResults(results.slice(0, 10)); // Limit to 10 results
    } catch (error) {
      console.error('Search error:', error);
      if (error.message !== 'Session expired') {
        toast.error('Search failed');
      }
    } finally {
      setSearchLoading(false);
    }
  };

  // Debounced search
  useEffect(() => {
    const timeoutId = setTimeout(() => {
      performSearch(searchQuery);
    }, 300);

    return () => clearTimeout(timeoutId);
  }, [searchQuery]);

  const handleSearchClick = (event) => {
    setSearchAnchorEl(event.currentTarget);
  };

  const handleSearchClose = () => {
    setSearchAnchorEl(null);
    setSearchQuery('');
    setSearchResults([]);
  };

  const handleSearchResultClick = (result) => {
    navigate(result.path);
    handleSearchClose();
    toast.success(`Navigated to ${result.type}: ${result.name}`);
  };

  const searchOpen = Boolean(searchAnchorEl);

  const fetchNotifications = async () => {
    setLoading(true);
    try {
      const response = await authenticatedFetch('/api/events');
      
      if (response.ok) {
        const data = await response.json();
        const recentEvents = data.items
          .filter(event => {
            const eventTime = new Date(event.lastTimestamp);
            const oneHourAgo = new Date(Date.now() - 60 * 60 * 1000);
            return eventTime > oneHourAgo;
          })
          .slice(0, 10); // Show only the 10 most recent events
        
        setNotifications(recentEvents);
        setNotificationCount(recentEvents.length);
      } else {
        console.error('Failed to fetch notifications');
        setNotificationCount(0);
      }
    } catch (error) {
      console.error('Error fetching notifications:', error);
      if (error.message !== 'Session expired') {
        setNotificationCount(0);
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchNotifications();
    // Refresh notifications every 30 seconds
    const interval = setInterval(fetchNotifications, 30000);
    return () => clearInterval(interval);
  }, []);

  const handleNotificationClick = (event) => {
    setAnchorEl(event.currentTarget);
    fetchNotifications(); // Refresh when opening
  };

  const handleNotificationClose = () => {
    setAnchorEl(null);
  };

  const handleNotificationItemClick = (notification) => {
    // You can add specific actions here based on notification type
    toast.info(`Notification: ${notification.reason} - ${notification.message}`);
    handleNotificationClose();
  };

  const open = Boolean(anchorEl);

  return (
    <AppBar
      position="fixed"
      sx={{
        width: { sm: `calc(100% - ${sidebarOpen ? drawerWidth : 0}px)` },
        ml: { sm: `${sidebarOpen ? drawerWidth : 0}px` },
        transition: (theme) =>
          theme.transitions.create(['margin', 'width'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
          }),
      }}
    >
      <Toolbar>
        <IconButton
          color="inherit"
          aria-label="open drawer"
          edge="start"
          onClick={onToggleSidebar}
          sx={{ mr: 2 }}
        >
          <MenuIcon />
        </IconButton>
        <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1, display: { xs: 'none', sm: 'block' } }}>
          Kubernetes Dashboard
        </Typography>
        
        {/* Session Status Indicator */}
        <Tooltip title={isValidating ? "Validating session..." : isAuthenticated ? "Session active" : "Session expired"}>
          <Box sx={{ display: 'flex', alignItems: 'center', mr: 2 }}>
            {isValidating ? (
              <CircularProgress size={20} color="inherit" />
            ) : (
              <Security 
                color={isAuthenticated ? "inherit" : "warning"} 
                sx={{ fontSize: 20 }}
              />
            )}
          </Box>
        </Tooltip>
        
        {/* Search Bar */}
        <SearchBar>
          <SearchIconWrapper>
            <Search />
          </SearchIconWrapper>
          <StyledInputBase 
            placeholder="Search resources…" 
            inputProps={{ 'aria-label': 'search' }}
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            onClick={handleSearchClick}
          />
        </SearchBar>
        
        {/* Search Results Dropdown */}
        <Popper
          open={searchOpen}
          anchorEl={searchAnchorEl}
          placement="bottom-start"
          style={{ zIndex: 1300, width: searchAnchorEl?.offsetWidth }}
        >
          <ClickAwayListener onClickAway={handleSearchClose}>
            <Paper 
              elevation={8}
              sx={{ 
                maxHeight: 400, 
                overflow: 'auto',
                mt: 1,
                minWidth: 300
              }}
            >
              {searchLoading ? (
                <Box sx={{ display: 'flex', justifyContent: 'center', p: 2 }}>
                  <CircularProgress size={24} />
                </Box>
              ) : searchResults.length === 0 && searchQuery ? (
                <Box sx={{ p: 2 }}>
                  <Typography variant="body2" color="text.secondary">
                    No results found for "{searchQuery}"
                  </Typography>
                </Box>
              ) : searchResults.length > 0 ? (
                <List>
                  {searchResults.map((result) => (
                    <ListItem key={result.id} disablePadding>
                      <ListItemButton onClick={() => handleSearchResultClick(result)}>
                        <ListItemIcon>
                          {result.icon}
                        </ListItemIcon>
                        <ListItemText
                          primary={result.name}
                          secondary={
                            <Box>
                              <Typography variant="body2" color="text.secondary">
                                {result.type} • {result.namespace}
                              </Typography>
                              <Chip 
                                label={result.status} 
                                size="small" 
                                color={result.status === 'Running' ? 'success' : 'default'}
                                sx={{ mt: 0.5 }}
                              />
                            </Box>
                          }
                        />
                      </ListItemButton>
                    </ListItem>
                  ))}
                </List>
              ) : null}
            </Paper>
          </ClickAwayListener>
        </Popper>
        
        {/* Notification Icon with Dropdown */}
        <IconButton 
          size="large" 
          aria-label="notifications" 
          color="inherit"
          onClick={handleNotificationClick}
          aria-controls={open ? 'notifications-menu' : undefined}
          aria-haspopup="true"
          aria-expanded={open ? 'true' : undefined}
        >
          <Badge badgeContent={notificationCount} color="error">
            <Notifications />
          </Badge>
        </IconButton>
        
        {/* Notification Menu */}
        <Menu
          id="notifications-menu"
          anchorEl={anchorEl}
          open={open}
          onClose={handleNotificationClose}
          MenuListProps={{
            'aria-labelledby': 'notifications-button',
          }}
          PaperProps={{
            sx: { maxHeight: 400, width: 400 }
          }}
        >
          <Box sx={{ p: 2, borderBottom: 1, borderColor: 'divider' }}>
            <Typography variant="h6">Notifications</Typography>
            <Typography variant="body2" color="text.secondary">
              Recent cluster events
            </Typography>
          </Box>
          
          {loading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', p: 2 }}>
              <CircularProgress size={24} />
            </Box>
          ) : notifications.length === 0 ? (
            <MenuItem>
              <ListItemText 
                primary="No recent notifications"
                secondary="All systems are running smoothly"
              />
            </MenuItem>
          ) : (
            <>
              {notifications.map((notification, index) => (
                <React.Fragment key={`${notification.name}-${index}`}>
                  <MenuItem onClick={() => handleNotificationItemClick(notification)}>
                    <NotificationItem notification={notification} />
                  </MenuItem>
                  {index < notifications.length - 1 && <Divider />}
                </React.Fragment>
              ))}
            </>
          )}
        </Menu>

        <IconButton size="large" edge="end" aria-label="account of current user" color="inherit" onClick={onLogout}>
          <Logout />
        </IconButton>
        <Avatar sx={{ ml: 2, bgcolor: 'secondary.main' }}>A</Avatar>
      </Toolbar>
    </AppBar>
  );
};

export default Header; 