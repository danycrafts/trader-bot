import { createContext } from 'preact';
import { useContext, useState, useEffect } from 'preact/hooks';
import { createTheme, ThemeProvider as MuiThemeProvider } from '@mui/material/styles';
import { lightTheme, darkTheme } from '../theme';
import { GetSettings } from '../../wailsjs/go/main/App';

type ThemeMode = 'light' | 'dark';

interface ThemeContextType {
  mode: ThemeMode;
  toggleTheme: () => void;
  setThemeMode: (mode: ThemeMode) => void;
}

const ThemeContext = createContext<ThemeContextType>({
  mode: 'light',
  toggleTheme: () => {},
  setThemeMode: () => {},
});

export const useTheme = () => useContext(ThemeContext);

export const ThemeProvider = ({ children }: { children: any }) => {
  const [mode, setMode] = useState<ThemeMode>('light');

  useEffect(() => {
      // Load theme from settings on startup
      // Assuming user ID 1 for MVP
      GetSettings(1).then((s: any) => {
          if (s && s.theme) {
              setMode(s.theme as ThemeMode);
          }
      }).catch(console.error);
  }, []);

  const toggleTheme = () => {
    setMode((prevMode) => (prevMode === 'light' ? 'dark' : 'light'));
  };

  const setThemeMode = (newMode: ThemeMode) => {
      setMode(newMode);
  }

  const theme = mode === 'light' ? lightTheme : darkTheme;

  return (
    <ThemeContext.Provider value={{ mode, toggleTheme, setThemeMode }}>
      <MuiThemeProvider theme={theme}>
        {children}
      </MuiThemeProvider>
    </ThemeContext.Provider>
  );
};
