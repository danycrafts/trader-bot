import { useState, useRef, useEffect } from 'preact/hooks';
import { Box, Typography, Paper, TextField, Button, Autocomplete } from '@mui/material';
import { SearchStocks, PlaceOrder, StartMarketStream } from '../../wailsjs/go/main/App';
import { createChart } from 'lightweight-charts';

export default function Market() {
  const [symbol, setSymbol] = useState('');
  const [options, setOptions] = useState<string[]>([]);
  const [qty, setQty] = useState(1);
  const chartContainerRef = useRef<HTMLDivElement>(null);
  const chartRef = useRef<any>(null);

  // Search stocks effect
  useEffect(() => {
    if (symbol.length > 0) {
      SearchStocks(symbol).then(setOptions).catch(console.error);
    }
  }, [symbol]);

  // Chart setup
  useEffect(() => {
    if (chartContainerRef.current && !chartRef.current) {
      const chart = createChart(chartContainerRef.current, {
        width: chartContainerRef.current.clientWidth,
        height: 300,
        layout: {
           background: { color: '#ffffff' },
           textColor: '#333',
        },
        grid: {
          vertLines: { color: '#f0f3fa' },
          horzLines: { color: '#f0f3fa' },
        },
      });
      const candlestickSeries = (chart as any).addCandlestickSeries();
      candlestickSeries.setData([
         { time: '2018-12-22', open: 75.16, high: 82.84, low: 36.16, close: 45.72 },
         { time: '2018-12-23', open: 45.12, high: 53.90, low: 45.12, close: 48.09 },
         { time: '2018-12-24', open: 60.71, high: 60.71, low: 53.39, close: 59.29 },
         // TODO: Connect to real data via Events.On
      ]);
      chartRef.current = chart;
    }
  }, []);

  const handleBuy = async () => {
    try {
        await PlaceOrder(symbol, Number(qty), "buy");
        alert(`Bought ${qty} ${symbol}`);
    } catch(e) {
        alert(e);
    }
  };

  const handleSell = async () => {
      try {
        await PlaceOrder(symbol, Number(qty), "sell");
        alert(`Sold ${qty} ${symbol}`);
    } catch(e) {
        alert(e);
    }
  };

  const handleSymbolChange = (_: any, newValue: string | null) => {
      if (newValue) {
          setSymbol(newValue);
          StartMarketStream(newValue).then(console.log).catch(console.error);
      }
  }

  // TODO: Unit Test: Verify searching stocks calls backend
  // TODO: Unit Test: Verify buy/sell buttons call PlaceOrder
  return (
    <Box sx={{ flexGrow: 1, p: 3 }}>
      <Typography variant="h4" gutterBottom>Market</Typography>
      <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 3 }}>
        <Box sx={{ width: { xs: '100%', md: '65%' } }}>
            <Paper sx={{ p: 2 }}>
                <Autocomplete
                    options={options}
                    // @ts-ignore
                    renderInput={(params) => <TextField {...params} label="Search Stock" />}
                    onInputChange={(event, newInputValue) => {
                        setSymbol(newInputValue);
                    }}
                    onChange={handleSymbolChange}
                />
                <Box ref={chartContainerRef} sx={{ mt: 2, border: '1px solid #ccc', height: 300 }} />
            </Paper>
        </Box>
        <Box sx={{ width: { xs: '100%', md: '30%' } }}>
            <Paper sx={{ p: 2 }}>
                <Typography variant="h6">Trade {symbol}</Typography>
                <TextField
                    label="Quantity"
                    type="number"
                    fullWidth
                    margin="normal"
                    value={qty}
                    onChange={(e: any) => setQty(e.target.value)}
                />
                <Box sx={{ mt: 2, display: 'flex', gap: 2 }}>
                    <Button variant="contained" color="success" fullWidth onClick={handleBuy}>Buy</Button>
                    <Button variant="contained" color="error" fullWidth onClick={handleSell}>Sell</Button>
                </Box>
            </Paper>
        </Box>
      </Box>
    </Box>
  );
}
