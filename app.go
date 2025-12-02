package main

import (
	"context"
	"fmt"
	"log"
	"trading-bot/internal/adapter"
	"trading-bot/internal/domain"
	"trading-bot/internal/usecase"
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
	// We need a context for the streamer, but we usually want it bound to the app lifecycle.
	// We'll pass a background context for now, or manage it in startup.
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
	// Re-initialize streamer with the app context if needed, or ensuring it respects cancellation.
	// For now, the streamer uses its own context derived from Background in NewApp, which is fine
	// as long as we close it.
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
