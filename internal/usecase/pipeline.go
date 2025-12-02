package usecase

import (
	"log"
	"trading-bot/internal/adapter"
	"trading-bot/internal/domain"
)

type DataPipeline struct {
	provider domain.MarketDataProvider
	repo     *adapter.SQLiteRepository
}

func NewDataPipeline(provider domain.MarketDataProvider, repo *adapter.SQLiteRepository) *DataPipeline {
	return &DataPipeline{
		provider: provider,
		repo:     repo,
	}
}

func (p *DataPipeline) Start(symbols []string) error {
	// 1. Connect to stream
	if err := p.provider.Connect(); err != nil {
		return err
	}

	// 2. Subscribe to symbols
	if len(symbols) > 0 {
		if err := p.provider.Subscribe(symbols); err != nil {
			return err
		}
	}

	// 3. Start Consumer (Worker)
	go p.consume()

	return nil
}

func (p *DataPipeline) Close() error {
	return p.provider.Close()
}

func (p *DataPipeline) consume() {
	ticks := p.provider.GetTickChannel()
	for tick := range ticks {
		// Log to console
		log.Printf("Received tick: %s $%.2f", tick.Symbol, tick.Price)

		// Persist to DB
		if err := p.repo.LogTick(&tick); err != nil {
			log.Printf("Failed to log tick: %v", err)
		}
		
		// Here we would also push to "Transformation" stage
	}
}
