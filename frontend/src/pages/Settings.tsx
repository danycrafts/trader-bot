import { useEffect, useState } from 'preact/hooks';
import { Box, Typography, Paper, TextField, Button, Switch, FormControlLabel } from '@mui/material';
import { GetSettings, SaveSettings } from '../../wailsjs/go/main/App';

export default function Settings() {
  const [apiKey, setApiKey] = useState('');
  const [secretKey, setSecretKey] = useState('');
  const [theme, setTheme] = useState('light');
  const [notificationsEmail, setNotificationsEmail] = useState(false);
  const [notificationsPush, setNotificationsPush] = useState(false);

  useEffect(() => {
    // Mock user ID 1 for now
    GetSettings(1).then((s: any) => {
        if(s) {
            setApiKey(s.alpaca_api_key || '');
            setSecretKey(s.alpaca_secret_key || '');
            setTheme(s.theme || 'light');
            setNotificationsEmail(s.notifications_email);
            setNotificationsPush(s.notifications_push);
        }
    }).catch(console.error);
  }, []);

  const handleSave = async () => {
    try {
        await SaveSettings({
            user_id: 1, // Mock
            alpaca_api_key: apiKey,
            alpaca_secret_key: secretKey,
            theme: theme,
            notifications_email: notificationsEmail,
            notifications_push: notificationsPush
        });
        alert('Settings saved!');
    } catch(e) {
        alert(e);
    }
  };

  return (
    <Box sx={{ flexGrow: 1, p: 3 }}>
      <Typography variant="h4" gutterBottom>Settings</Typography>
      <Paper sx={{ p: 3, maxWidth: 600 }}>
        <Typography variant="h6" gutterBottom>API Keys</Typography>
        <TextField
            label="Alpaca API Key"
            fullWidth
            margin="normal"
            value={apiKey}
            onChange={(e:any) => setApiKey(e.target.value)}
        />
        <TextField
            label="Alpaca Secret Key"
            fullWidth
            margin="normal"
            type="password"
            value={secretKey}
            onChange={(e:any) => setSecretKey(e.target.value)}
        />

        <Typography variant="h6" gutterBottom sx={{ mt: 3 }}>Appearance</Typography>
        <FormControlLabel
            control={<Switch checked={theme === 'dark'} onChange={(e:any) => setTheme(e.target.checked ? 'dark' : 'light')} />}
            label="Dark Mode"
        />

        <Typography variant="h6" gutterBottom sx={{ mt: 3 }}>Notifications</Typography>
        <FormControlLabel
            control={<Switch checked={notificationsEmail} onChange={(e:any) => setNotificationsEmail(e.target.checked)} />}
            label="Email Notifications"
        />
        <FormControlLabel
            control={<Switch checked={notificationsPush} onChange={(e:any) => setNotificationsPush(e.target.checked)} />}
            label="Push Notifications"
        />

        <Button variant="contained" sx={{ mt: 3 }} onClick={handleSave}>Save Settings</Button>
      </Paper>
    </Box>
  );
}
