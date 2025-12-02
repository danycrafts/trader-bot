package adapter

import (
	"fmt"
	"trading-bot/internal/domain"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

type AlpacaBroker struct {
	client *alpaca.Client
}

func NewAlpacaBroker() *AlpacaBroker {
	// API keys are expected to be in environment variables:
	// APCA_API_KEY_ID and APCA_API_SECRET_KEY
	// The Base URL defaults to paper trading if not specified, 
	// or can be set via APCA_API_BASE_URL.
	
	// For testing purposes, we can try to load them, but the SDK does it automatically if we use NewClient(alpaca.ClientOpts{})
    // However, we need to ensure we can use it.
    
    // Explicitly reading env vars or letting the library handle it. 
    // The library handles APCA_API_KEY_ID, APCA_API_SECRET_KEY, APCA_API_BASE_URL
	client := alpaca.NewClient(alpaca.ClientOpts{})
	return &AlpacaBroker{
		client: client,
	}
}

func (b *AlpacaBroker) GetAccountBalance() (float64, error) {
	account, err := b.client.GetAccount()
	if err != nil {
		return 0, fmt.Errorf("failed to get account: %w", err)
	}
	// Equity is the total value of the account (cash + positions)
	// Cash is the available cash.
	// The prompt asks for "cash equity". Usually means "Cash" or "Equity". 
	// "Check Balance ... display the cash equity". 
	// I'll return Equity which is the most useful metric, or maybe both?
	// Let's return Equity as float64.
	return account.Equity.InexactFloat64(), nil
}

// Ensure AlpacaBroker implements domain.Broker
var _ domain.Broker = (*AlpacaBroker)(nil)
