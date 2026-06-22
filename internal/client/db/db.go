package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Initialize() (*sql.DB, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("could not determine user config dir: %w", err)
	}

	dbDir := filepath.Join(configDir, "ai-toolbox")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create database directory: %w", err)
	}

	dbPath := filepath.Join(dbDir, "ai_tools.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrency
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS codex_profiles (
			name TEXT PRIMARY KEY,
			model TEXT NOT NULL DEFAULT '',
			model_provider TEXT NOT NULL DEFAULT '',
			model_reasoning_effort TEXT NOT NULL DEFAULT '',
			approval_policy TEXT NOT NULL DEFAULT '',
			sandbox_mode TEXT NOT NULL DEFAULT '',
			created_at INTEGER NOT NULL DEFAULT 0,
			updated_at INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS cli_sessions (
			source       TEXT NOT NULL CHECK(source IN ('claude','codex')),
			session_id   TEXT NOT NULL,
			cwd          TEXT NOT NULL,
			file_path    TEXT NOT NULL,
			started_at   INTEGER NOT NULL,
			last_active  INTEGER NOT NULL,
			first_msg    TEXT,
			indexed_at   INTEGER NOT NULL,
			display_name TEXT NOT NULL DEFAULT '',
			PRIMARY KEY (source, session_id)
		)`,
		`CREATE TABLE IF NOT EXISTS cli_session_tabs (
			source TEXT NOT NULL,
			session_id TEXT NOT NULL,
			tab_id TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			FOREIGN KEY (source, session_id) REFERENCES cli_sessions(source, session_id)
		)`,
		`CREATE TABLE IF NOT EXISTS accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL DEFAULT '',
			api_key TEXT NOT NULL DEFAULT '',
			provider TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL DEFAULT 'active',
			plan_type TEXT NOT NULL DEFAULT '',
			quota_used REAL DEFAULT 0,
			quota_limit REAL DEFAULT 100,
			last_used_at INTEGER DEFAULT 0,
			request_count INTEGER DEFAULT 0,
			error_count INTEGER DEFAULT 0,
			created_at INTEGER DEFAULT 0,
			updated_at INTEGER DEFAULT 0,
			issuer TEXT DEFAULT 'https://auth.openai.com',
			cooldown_until INTEGER DEFAULT 0,
			last_error TEXT DEFAULT '',
			primary_used_percent REAL DEFAULT 0,
			primary_window_minutes INTEGER DEFAULT 0,
			primary_resets_at INTEGER DEFAULT 0,
			secondary_used_percent REAL DEFAULT 0,
			secondary_window_minutes INTEGER DEFAULT 0,
			secondary_resets_at INTEGER DEFAULT 0,
			chatgpt_account_id TEXT DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS codex_tokens (
			account_id INTEGER PRIMARY KEY,
			id_token TEXT NOT NULL DEFAULT '',
			access_token TEXT NOT NULL DEFAULT '',
			refresh_token TEXT NOT NULL DEFAULT '',
			expires_at INTEGER DEFAULT 0,
			last_refresh INTEGER DEFAULT 0,
			FOREIGN KEY (account_id) REFERENCES accounts(id)
		)`,
		`CREATE TABLE IF NOT EXISTS codex_extra_providers (
			name TEXT PRIMARY KEY,
			base_url TEXT NOT NULL DEFAULT '',
			env_key TEXT NOT NULL DEFAULT '',
			wire_api TEXT NOT NULL DEFAULT 'responses',
			requires_openai_auth INTEGER NOT NULL DEFAULT 1,
			query_params TEXT NOT NULL DEFAULT '{}',
			http_headers TEXT NOT NULL DEFAULT '{}',
			created_at INTEGER NOT NULL DEFAULT 0,
			updated_at INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS codex_provider_config (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			model TEXT NOT NULL DEFAULT '',
			model_provider TEXT NOT NULL DEFAULT '',
			model_context_window INTEGER NOT NULL DEFAULT 0,
			model_reasoning_effort TEXT NOT NULL DEFAULT '',
			model_reasoning_summary TEXT NOT NULL DEFAULT '',
			model_verbosity TEXT NOT NULL DEFAULT '',
			approval_policy TEXT NOT NULL DEFAULT '',
			allow_login_shell INTEGER NOT NULL DEFAULT 0,
			sandbox_mode TEXT NOT NULL DEFAULT '',
			sandbox_writable_roots TEXT NOT NULL DEFAULT '[]',
			sandbox_network_access INTEGER NOT NULL DEFAULT 0,
			sandbox_exclude_tmpdir INTEGER NOT NULL DEFAULT 0,
			sandbox_exclude_slash_tmp INTEGER NOT NULL DEFAULT 0,
			network_enabled INTEGER NOT NULL DEFAULT 0,
			network_mode TEXT NOT NULL DEFAULT '',
			network_domains TEXT NOT NULL DEFAULT '{}',
			updated_at INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS claude_presets (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			color TEXT NOT NULL DEFAULT '',
			icon TEXT NOT NULL DEFAULT '',
			provider TEXT NOT NULL,
			base_url TEXT NOT NULL DEFAULT '',
			default_model TEXT NOT NULL DEFAULT '',
			preset_default_opus_model TEXT NOT NULL DEFAULT '',
			preset_default_sonnet_model TEXT NOT NULL DEFAULT '',
			preset_default_haiku_model TEXT NOT NULL DEFAULT '',
			badge TEXT NOT NULL DEFAULT '',
			sort_order INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS gateway_config (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			listen_addr TEXT DEFAULT '127.0.0.1:9800',
			routing_strategy TEXT NOT NULL DEFAULT '',
			max_concurrency INTEGER NOT NULL DEFAULT 0,
			session_affinity INTEGER NOT NULL DEFAULT 0,
			target_base_url TEXT NOT NULL DEFAULT '',
			upstream_proxy_url TEXT NOT NULL DEFAULT '',
			api_key TEXT NOT NULL DEFAULT '',
			updated_at INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS gateway_keys (
			id TEXT PRIMARY KEY,
			key_hash TEXT NOT NULL,
			key_prefix TEXT NOT NULL,
			label TEXT NOT NULL DEFAULT '',
			created_at INTEGER NOT NULL DEFAULT 0,
			last_used_at INTEGER DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			path TEXT NOT NULL,
			created_at INTEGER NOT NULL DEFAULT 0,
			last_used INTEGER NOT NULL DEFAULT 0,
			pinned INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS app_settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE INDEX IF NOT EXISTS idx_cli_sessions_cwd ON cli_sessions(cwd, last_active DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_claude_presets_sort ON claude_presets(sort_order)`,
	}

	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return fmt.Errorf("migration failed: %w\nSQL: %s", err, m)
		}
	}

	return nil
}
