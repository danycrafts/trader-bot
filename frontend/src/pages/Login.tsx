import { useState } from 'preact/hooks';
import { Box, Button, TextField, Typography, Paper, Container } from '@mui/material';
import { Login as LoginBackend } from '../../wailsjs/go/main/App';
import { route } from 'preact-router';

interface LoginProps {
  onLoginSuccess: (user: any) => void;
}

export default function Login({ onLoginSuccess }: LoginProps) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleLogin = async () => {
    try {
      const user = await LoginBackend(email, password);
      onLoginSuccess(user);
      route('/dashboard');
    } catch (e: any) {
      setError(e.toString());
    }
  };

  return (
    <Container component="main" maxWidth="xs">
      <Paper sx={{ marginTop: 8, display: 'flex', flexDirection: 'column', alignItems: 'center', padding: 3 }}>
        <Typography component="h1" variant="h5">
          Sign in
        </Typography>
        <Box component="form" noValidate sx={{ mt: 1 }}>
          <TextField
            margin="normal"
            required
            fullWidth
            id="email"
            label="Email Address"
            name="email"
            autoComplete="email"
            autoFocus
            value={email}
            onChange={(e: any) => setEmail(e.target.value)}
          />
          <TextField
            margin="normal"
            required
            fullWidth
            name="password"
            label="Password"
            type="password"
            id="password"
            autoComplete="current-password"
            value={password}
            onChange={(e: any) => setPassword(e.target.value)}
          />
          {error && (
            <Typography color="error" variant="body2">
              {error}
            </Typography>
          )}
          <Button
            type="button"
            fullWidth
            variant="contained"
            sx={{ mt: 3, mb: 2 }}
            onClick={handleLogin}
          >
            Sign In
          </Button>
          <Button
            fullWidth
            variant="text"
            onClick={() => route('/register')}
          >
            Don't have an account? Register
          </Button>
        </Box>
      </Paper>
    </Container>
  );
}
