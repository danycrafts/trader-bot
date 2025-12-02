package domain

// MarketDataProvider defines the interface for streaming market data.
type MarketDataProvider interface {
	Subscribe(symbols []string) error
	GetTickChannel() <-chan Tick
	Connect() error
	Close() error
}
