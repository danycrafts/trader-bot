import { Router, Route } from 'preact-router';
import CssBaseline from '@mui/material/CssBaseline';
import { ThemeProvider } from './context/ThemeContext';
import Navbar from './components/Navbar';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import Market from './pages/Market';
import Settings from './pages/Settings';

export function App() {
  const handleLoginSuccess = (user: any) => {
    console.log("Logged in:", user);
    // Persist user state if needed
  };

  const handleLogout = () => {
    // Clear user state
    window.location.href = '/';
  };

  return (
    <ThemeProvider>
      <CssBaseline />
      <div id="app" style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
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
