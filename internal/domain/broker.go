package domain

// Broker defines the interface for interacting with a brokerage.
type Broker interface {
	GetAccountBalance() (float64, error)
}
