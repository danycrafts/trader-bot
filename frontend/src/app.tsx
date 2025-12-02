import { useState } from 'preact/hooks';
import { Router, Route } from 'preact-router';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { lightTheme, darkTheme } from './theme';
import Navbar from './components/Navbar';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import Market from './pages/Market';
import Settings from './pages/Settings';

export function App() {
  // TODO: Use global state or context for theme mode
  const [mode, setMode] = useState<'light' | 'dark'>('light');

  const handleLoginSuccess = (user: any) => {
    console.log("Logged in:", user);
    // Persist user state if needed
  };

  const handleLogout = () => {
    // Clear user state
    // Navigate to login (handled by route if protected, but here just simple nav)
    window.location.href = '/';
  };

  return (
    <ThemeProvider theme={mode === 'light' ? lightTheme : darkTheme}>
      <CssBaseline />
      <div id="app" style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
        {/* Only show Navbar if not on login/register pages.
            However, preact-router doesn't easily give us current route outside.
            We can conditionally render Navbar inside a wrapper or based on state.
            For simplicity, let's just put Navbar everywhere except known paths, or wrap Dashboard/Market in a Layout.
        */}

        <Router>
            <Route path="/" component={() => <Login onLoginSuccess={handleLoginSuccess} />} />
            <Route path="/register" component={Register} />
            <Route path="/dashboard" component={() => <><Navbar onLogout={handleLogout} /><Dashboard /></>} />
            <Route path="/market" component={() => <><Navbar onLogout={handleLogout} /><Market /></>} />
            <Route path="/settings" component={() => <><Navbar onLogout={handleLogout} /><Settings /></>} />
        </Router>
      </div>
    </ThemeProvider>
  );
}
