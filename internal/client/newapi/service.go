package newapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// Account represents an AI provider account
type Account struct {
	ID                    int64   `json:"id"`
	Name                  string  `json:"name"`
	Email                 string  `json:"email"`
	APIKey                string  `json:"api_key,omitempty"`
	Provider              string  `json:"provider"`
	Status                string  `json:"status"`
	PlanType              string  `json:"plan_type"`
	QuotaUsed             float64 `json:"quota_used"`
	QuotaLimit            float64 `json:"quota_limit"`
	LastUsedAt            int64   `json:"last_used_at"`
	RequestCount          int     `json:"request_count"`
	ErrorCount            int     `json:"error_count"`
	CreatedAt             int64   `json:"created_at"`
	UpdatedAt             int64   `json:"updated_at"`
	Issuer                string  `json:"issuer"`
	CooldownUntil         int64   `json:"cooldown_until"`
	LastError             string  `json:"last_error"`
	PrimaryUsedPercent    float64 `json:"primary_used_percent"`
	PrimaryWindowMinutes  int     `json:"primary_window_minutes"`
	PrimaryResetsAt       int64   `json:"primary_resets_at"`
	SecondaryUsedPercent  float64 `json:"secondary_used_percent"`
	SecondaryWindowMinutes int    `json:"secondary_window_minutes"`
	SecondaryResetsAt     int64   `json:"secondary_resets_at"`
}

// ListAccounts returns all accounts
func (s *Service) ListAccounts() ([]Account, error) {
	rows, err := s.db.Query(`SELECT id, name, email, api_key, provider, status, 
		last_synced_at, quota_used, quota_limit, last_used_at, request_count, error_count, 
		created_at, updated_at, plan_type, chatgpt_account_id, workspace_id, issuer, 
		cooldown_until, last_error, primary_used_percent, primary_window_minutes, 
		primary_resets_at, secondary_used_percent, secondary_window_minutes, 
		secondary_resets_at FROM accounts ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var a Account
		err := rows.Scan(&a.ID, &a.Name, &a.Email, &a.APIKey, &a.Provider,
			&a.Status, &a.LastUsedAt, &a.QuotaUsed, &a.QuotaLimit, &a.LastUsedAt,
			&a.RequestCount, &a.ErrorCount, &a.CreatedAt, &a.UpdatedAt, &a.PlanType,
			&a.Issuer, &a.CooldownUntil, &a.LastError,
			&a.PrimaryUsedPercent, &a.PrimaryWindowMinutes, &a.PrimaryResetsAt,
			&a.SecondaryUsedPercent, &a.SecondaryWindowMinutes, &a.SecondaryResetsAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, nil
}

// GetUserID returns the current user ID
func (s *Service) GetUserID() (string, error) {
	var userID string
	err := s.db.QueryRow("SELECT value FROM app_settings WHERE key = 'user_id'").Scan(&userID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return userID, err
}

// SaveUserID persists the user ID
func (s *Service) SaveUserID(userID string) error {
	_, err := s.db.Exec(
		"INSERT INTO app_settings(key, value) VALUES('user_id', ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value",
		userID)
	return err
}

// HTTP client helper for newapi requests
func (s *Service) doRequest(method, url, apiKey string, body interface{}) (*http.Response, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}
