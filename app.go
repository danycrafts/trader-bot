package main

import (
	"context"
	"fmt"
	"trading-bot/internal/adapter"
	"trading-bot/internal/domain"
)

// App struct
type App struct {
	ctx    context.Context
	broker domain.Broker
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		broker: adapter.NewAlpacaBroker(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetAccountBalance returns the current account equity
func (a *App) GetAccountBalance() (float64, error) {
	return a.broker.GetAccountBalance()
}
