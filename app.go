package main

import (
	"context"
	"fmt"
	"log"
	"trading-bot/internal/adapter"
	"trading-bot/internal/domain"
	"trading-bot/internal/usecase"
	"time"
)

// App struct
type App struct {
	ctx      context.Context
	broker   domain.Broker
	pipeline *usecase.DataPipeline
	repo     *adapter.SQLiteRepository
}

// NewApp creates a new App application struct
func NewApp() *App {
	// Initialize DB
	repo, err := adapter.NewSQLiteRepository("trading_bot.db")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Streamer
	streamer := adapter.NewAlpacaStreamer(context.Background())

	pipeline := usecase.NewDataPipeline(streamer, repo)

	return &App{
		broker:   adapter.NewAlpacaBroker(),
		pipeline: pipeline,
		repo:     repo,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	if a.pipeline != nil {
		if err := a.pipeline.Close(); err != nil {
			fmt.Printf("Error closing pipeline: %v\n", err)
		}
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetAccountBalance returns the current account equity
func (a *App) GetAccountBalance() (float64, error) {
	return a.broker.GetAccountBalance()
}

// StartMarketStream starts streaming data for the given symbol
func (a *App) StartMarketStream(symbol string) string {
	err := a.pipeline.Start([]string{symbol})
	if err != nil {
		return fmt.Sprintf("Error starting stream: %v", err)
	}
	return fmt.Sprintf("Started streaming %s", symbol)
}

// --- Auth ---

// Login authenticates a user
func (a *App) Login(email, password string) (*domain.User, error) {
	u, err := a.repo.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}
	if u.Password != password { // In real app, use bcrypt
		return nil, fmt.Errorf("invalid email or password")
	}
	return u, nil
}

// Register registers a new user
func (a *App) Register(email, password string) error {
	_, err := a.repo.GetUserByEmail(email)
	if err == nil {
		return fmt.Errorf("email already exists")
	}
	return a.repo.CreateUser(email, password)
}

// --- Settings ---

// SaveSettings saves user settings
func (a *App) SaveSettings(settings domain.UserSettings) error {
	return a.repo.SaveSettings(&settings)
}

// GetSettings retrieves user settings
func (a *App) GetSettings(userID int) (*domain.UserSettings, error) {
	return a.repo.GetSettings(userID)
}

// --- Market & Trading ---

// SearchStocks searches for stocks (mock implementation for now as Alpaca needs API Key)
func (a *App) SearchStocks(query string) ([]string, error) {
	// In a real implementation, we would use a.broker.SearchAssets(query) or similar
	// For now, return a mock list
	stocks := []string{"AAPL", "GOOGL", "MSFT", "TSLA", "AMZN", "FB", "NFLX", "NVDA", "AMD", "INTC"}
	var result []string
	for _, s := range stocks {
		if len(result) < 5 { // Limit to 5
			result = append(result, s)
		}
	}
	return result, nil
}

// PlaceOrder places an order
func (a *App) PlaceOrder(symbol string, qty float64, side string) error {
	// a.broker.PlaceOrder(...)
	// For now, mock it and save to DB
	fmt.Printf("Placing order: %s %f %s\n", side, qty, symbol)
	return nil
}

// GetPortfolio returns current portfolio
func (a *App) GetPortfolio() ([]domain.Position, error) {
	// return a.broker.GetPositions()
	// Mock
	return []domain.Position{
		{Symbol: "AAPL", Qty: 10, EntryPrice: 150.0, CurrentPrice: 155.0, PnL: 50.0},
		{Symbol: "TSLA", Qty: 5, EntryPrice: 700.0, CurrentPrice: 690.0, PnL: -50.0},
	}, nil
}

// GetRecentTrades returns recent trades
func (a *App) GetRecentTrades() ([]domain.Trade, error) {
	// Mock return for now since DB might be empty
	return []domain.Trade{
		{ID: "1", Symbol: "AAPL", Strategy: "Momentum", Side: "buy", EntryPrice: 150.0, EntryTime: time.Now().Unix(), Status: "closed", PnL: 10.0},
		{ID: "2", Symbol: "TSLA", Strategy: "MeanReversion", Side: "sell", EntryPrice: 700.0, EntryTime: time.Now().Add(-1 * time.Hour).Unix(), Status: "closed", PnL: -5.0},
	}, nil
}
