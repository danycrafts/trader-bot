package domain

import "time"

// User represents a registered user
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"` // In a real app, this should be hashed
	CreatedAt time.Time `json:"created_at"`
}

// UserSettings represents user preferences and API keys
type UserSettings struct {
	UserID             int    `json:"user_id"`
	AlpacaAPIKey       string `json:"alpaca_api_key"`
	AlpacaSecretKey    string `json:"alpaca_secret_key"`
	Theme              string `json:"theme"` // "dark" or "light"
	NotificationsEmail bool   `json:"notifications_email"`
	NotificationsPush  bool   `json:"notifications_push"`
}
