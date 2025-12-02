package adapter

import (
	"context"
	"log"
	"trading-bot/internal/domain"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
)

type AlpacaStreamer struct {
	client    stream.StocksClient
	tickChan  chan domain.Tick
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewAlpacaStreamer(ctx context.Context) *AlpacaStreamer {
	ctx, cancel := context.WithCancel(ctx)
	// Buffered channel to prevent blocking the stream callback
	tickChan := make(chan domain.Tick, 1000)
	return &AlpacaStreamer{
		tickChan: tickChan,
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (s *AlpacaStreamer) tradeHandler(t stream.Trade) {
	tick := domain.Tick{
		Symbol:    t.Symbol,
		Price:     t.Price,
		Timestamp: t.Timestamp,
		Volume:    uint64(t.Size),
	}

	// Non-blocking send or log drop
	select {
	case s.tickChan <- tick:
	default:
		log.Println("Warning: Tick channel full, dropping tick")
	}
}

func (s *AlpacaStreamer) Connect() error {
	// Create a new client. It automatically uses environment variables for auth.
	// By default, stream.NewStocksClient connects to IEX (free real-time data) or SIP (paid) based on subscription.
	c := stream.NewStocksClient(
		"iex", // Use IEX for free data
		stream.WithLogger(stream.DefaultLogger()),
	)

	s.client = *c

	// Connect in a background goroutine because Connect() blocks
	go func() {
		if err := c.Connect(s.ctx); err != nil {
			log.Printf("Stream connection error: %v", err)
		}
	}()

	return nil
}

func (s *AlpacaStreamer) Subscribe(symbols []string) error {
	// Note: Subscribing usually happens before Connect, or dynamically if supported.
	// Alpaca allows dynamic subscription.
	// We use the same handler logic defined in tradeHandler.
	return s.client.SubscribeToTrades(s.tradeHandler, symbols...)
}

func (s *AlpacaStreamer) GetTickChannel() <-chan domain.Tick {
	return s.tickChan
}

func (s *AlpacaStreamer) Close() error {
	s.cancel()
	return nil
}

var _ domain.MarketDataProvider = (*AlpacaStreamer)(nil)
