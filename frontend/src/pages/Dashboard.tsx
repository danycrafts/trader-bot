import { useEffect, useState } from 'preact/hooks';
import { Box, Typography, Paper, Table, TableBody, TableCell, TableHead, TableRow } from '@mui/material';
import { GetPortfolio, GetAccountBalance, GetRecentTrades } from '../../wailsjs/go/main/App';

export default function Dashboard() {
  const [balance, setBalance] = useState<number | null>(null);
  const [portfolio, setPortfolio] = useState<any[]>([]);
  const [history, setHistory] = useState<any[]>([]);

  useEffect(() => {
    // Mock loading or real loading
    GetAccountBalance().then(setBalance).catch(console.error);
    GetPortfolio().then(setPortfolio).catch(console.error);
    GetRecentTrades().then(setHistory).catch(console.error);
  }, []);

  return (
    <Box sx={{ flexGrow: 1, p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Dashboard
      </Typography>
      <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 3 }}>
        <Box sx={{ width: { xs: '100%', md: '30%' } }}>
          <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
            <Typography variant="h6" color="primary" gutterBottom>
              Total Balance
            </Typography>
            <Typography component="p" variant="h4">
              ${balance?.toFixed(2) ?? '---'}
            </Typography>
          </Paper>
        </Box>
        <Box sx={{ width: '100%' }}>
          <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
            <Typography variant="h6" color="primary" gutterBottom>
              Current Positions
            </Typography>
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>Symbol</TableCell>
                  <TableCell>Qty</TableCell>
                  <TableCell>Entry Price</TableCell>
                  <TableCell>Current Price</TableCell>
                  <TableCell>PnL</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {portfolio.map((row) => (
                  <TableRow key={row.Symbol}>
                    <TableCell>{row.Symbol}</TableCell>
                    <TableCell>{row.Qty}</TableCell>
                    <TableCell>${row.EntryPrice.toFixed(2)}</TableCell>
                    <TableCell>${row.CurrentPrice.toFixed(2)}</TableCell>
                    <TableCell sx={{ color: row.PnL >= 0 ? 'green' : 'red' }}>
                      ${row.PnL.toFixed(2)}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </Paper>
        </Box>
        <Box sx={{ width: '100%' }}>
          <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
            <Typography variant="h6" color="primary" gutterBottom>
              Recent Trades
            </Typography>
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>ID</TableCell>
                  <TableCell>Symbol</TableCell>
                  <TableCell>Side</TableCell>
                  <TableCell>Strategy</TableCell>
                  <TableCell>Entry Price</TableCell>
                  <TableCell>PnL</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {history.map((row) => (
                  <TableRow key={row.id}>
                    <TableCell>{row.id}</TableCell>
                    <TableCell>{row.Symbol}</TableCell>
                    <TableCell>{row.Side}</TableCell>
                    <TableCell>{row.Strategy}</TableCell>
                    <TableCell>${row.EntryPrice.toFixed(2)}</TableCell>
                    <TableCell sx={{ color: row.PnL >= 0 ? 'green' : 'red' }}>
                      ${row.PnL.toFixed(2)}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </Paper>
        </Box>
      </Box>
    </Box>
  );
}
