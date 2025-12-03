package adapter

import (
	"database/sql"
	"fmt"
	"time"
	"trading-bot/internal/domain"

	_ "modernc.org/sqlite"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	repo := &SQLiteRepository{db: db}
	if err := repo.initDB(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *SQLiteRepository) initDB() error {
	// Enable WAL mode
	if _, err := r.db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return fmt.Errorf("failed to set WAL mode: %w", err)
	}

	// Create trades table
	createTrades := `
	CREATE TABLE IF NOT EXISTS trades (
		id TEXT PRIMARY KEY,
		symbol TEXT NOT NULL,
		strategy TEXT NOT NULL,
		side TEXT NOT NULL,
		entry_price REAL NOT NULL,
		entry_time INTEGER NOT NULL,
		exit_price REAL,
		exit_time INTEGER,
		pnl REAL,
		status TEXT NOT NULL
	);`
	if _, err := r.db.Exec(createTrades); err != nil {
		return fmt.Errorf("failed to create trades table: %w", err)
	}

	// Create signals table
	createSignals := `
	CREATE TABLE IF NOT EXISTS signals (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp INTEGER NOT NULL,
		symbol TEXT NOT NULL,
		signal_type TEXT NOT NULL,
		indicators TEXT,
		action_taken TEXT
	);`
	if _, err := r.db.Exec(createSignals); err != nil {
		return fmt.Errorf("failed to create signals table: %w", err)
	}

	// Create users table
	createUsers := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at INTEGER NOT NULL
	);`
	if _, err := r.db.Exec(createUsers); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create settings table
	createSettings := `
	CREATE TABLE IF NOT EXISTS settings (
		user_id INTEGER PRIMARY KEY,
		alpaca_api_key TEXT,
		alpaca_secret_key TEXT,
		theme TEXT,
		notifications_email INTEGER,
		notifications_push INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`
	if _, err := r.db.Exec(createSettings); err != nil {
		return fmt.Errorf("failed to create settings table: %w", err)
	}

	return nil
}

func (r *SQLiteRepository) SaveSignal(signal *domain.Signal) error {
	query := `
	INSERT INTO signals (timestamp, symbol, signal_type, indicators, action_taken)
	VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, signal.Timestamp, signal.Symbol, signal.SignalType, signal.Indicators, signal.ActionTaken)
	if err != nil {
		return fmt.Errorf("failed to save signal: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) SaveTrade(trade *domain.Trade) error {
	query := `
	INSERT INTO trades (id, symbol, strategy, side, entry_price, entry_time, exit_price, exit_time, pnl, status)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		exit_price = excluded.exit_price,
		exit_time = excluded.exit_time,
		pnl = excluded.pnl,
		status = excluded.status
	`
	_, err := r.db.Exec(query, trade.ID, trade.Symbol, trade.Strategy, trade.Side, trade.EntryPrice, trade.EntryTime, trade.ExitPrice, trade.ExitTime, trade.PnL, trade.Status)
	if err != nil {
		return fmt.Errorf("failed to save trade: %w", err)
	}
	return nil
}

// User & Settings Methods

func (r *SQLiteRepository) CreateUser(email, password string) error {
	query := `INSERT INTO users (email, password, created_at) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, email, password, time.Now().Unix())
	return err
}

func (r *SQLiteRepository) GetUserByEmail(email string) (*domain.User, error) {
	query := `SELECT id, email, password, created_at FROM users WHERE email = ?`
	row := r.db.QueryRow(query, email)
	var u domain.User
	var createdAt int64
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &createdAt); err != nil {
		return nil, err
	}
	u.CreatedAt = time.Unix(createdAt, 0)
	return &u, nil
}

func (r *SQLiteRepository) SaveSettings(s *domain.UserSettings) error {
	query := `
	INSERT INTO settings (user_id, alpaca_api_key, alpaca_secret_key, theme, notifications_email, notifications_push)
	VALUES (?, ?, ?, ?, ?, ?)
	ON CONFLICT(user_id) DO UPDATE SET
		alpaca_api_key = excluded.alpaca_api_key,
		alpaca_secret_key = excluded.alpaca_secret_key,
		theme = excluded.theme,
		notifications_email = excluded.notifications_email,
		notifications_push = excluded.notifications_push
	`
	notifEmail := 0
	if s.NotificationsEmail {
		notifEmail = 1
	}
	notifPush := 0
	if s.NotificationsPush {
		notifPush = 1
	}

	_, err := r.db.Exec(query, s.UserID, s.AlpacaAPIKey, s.AlpacaSecretKey, s.Theme, notifEmail, notifPush)
	return err
}

func (r *SQLiteRepository) GetSettings(userID int) (*domain.UserSettings, error) {
	query := `SELECT alpaca_api_key, alpaca_secret_key, theme, notifications_email, notifications_push FROM settings WHERE user_id = ?`
	row := r.db.QueryRow(query, userID)
	var s domain.UserSettings
	s.UserID = userID
	var notifEmail, notifPush int
	if err := row.Scan(&s.AlpacaAPIKey, &s.AlpacaSecretKey, &s.Theme, &notifEmail, &notifPush); err != nil {
		if err == sql.ErrNoRows {
			// Return empty settings if not found
			return &s, nil
		}
		return nil, err
	}
	s.NotificationsEmail = notifEmail == 1
	s.NotificationsPush = notifPush == 1
	return &s, nil
}

// LogTick is a helper for Phase 2 verification to show we are streaming data.
// In production we might not log every tick to DB, but for now we might want to log it to a separate table or just verify via logs.
// For the purpose of the requirement "verify streaming ticks are being logged", let's add a simple log table.
func (r *SQLiteRepository) InitTickLog() error {
	createTicks := `
	CREATE TABLE IF NOT EXISTS tick_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		symbol TEXT NOT NULL,
		price REAL NOT NULL,
		volume INTEGER,
		timestamp INTEGER NOT NULL
	);`
	if _, err := r.db.Exec(createTicks); err != nil {
		return fmt.Errorf("failed to create tick_log table: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) LogTick(tick *domain.Tick) error {
	query := `INSERT INTO tick_log (symbol, price, volume, timestamp) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, tick.Symbol, tick.Price, tick.Volume, tick.Timestamp.UnixMilli())
	return err
}

func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}
