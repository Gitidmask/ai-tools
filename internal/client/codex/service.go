package codex

import (
	"database/sql"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// Profile represents a Codex profile
type Profile struct {
	Name              string `json:"name"`
	DisplayName       string `json:"display_name"`
	APIKey            string `json:"api_key,omitempty"`
	BaseURL           string `json:"base_url"`
	DefaultModel      string `json:"default_model"`
	Model             string `json:"model"`
	ModelProvider     string `json:"model_provider"`
	IsActive          bool   `json:"is_active"`
	CreatedAt         int64  `json:"created_at"`
	UpdatedAt         int64  `json:"updated_at"`
}

// ExtraProvider represents a custom API provider
type ExtraProvider struct {
	Name               string `json:"name"`
	BaseURL            string `json:"base_url"`
	EnvKey             string `json:"env_key"`
	WireAPI            string `json:"wire_api"`
	RequiresOpenAIAuth bool   `json:"requires_openai_auth"`
	QueryParams        string `json:"query_params"`
	HTTPHeaders        string `json:"http_headers"`
}

// ProviderConfig represents the active provider configuration
type ProviderConfig struct {
	Model                  string `json:"model"`
	ModelProvider          string `json:"model_provider"`
	ModelContextWindow     int    `json:"model_context_window"`
	ModelReasoningEffort   string `json:"model_reasoning_effort"`
	ModelReasoningSummary  string `json:"model_reasoning_summary"`
	ModelVerbosity         string `json:"model_verbosity"`
	ApprovalPolicy         string `json:"approval_policy"`
	AllowLoginShell        bool   `json:"allow_login_shell"`
	SandboxMode            string `json:"sandbox_mode"`
	SandboxWritableRoots   string `json:"sandbox_writable_roots"`
	SandboxNetworkAccess   bool   `json:"sandbox_network_access"`
	SandboxExcludeTmpdir   bool   `json:"sandbox_exclude_tmpdir"`
	SandboxExcludeSlashTmp bool   `json:"sandbox_exclude_slash_tmp"`
	NetworkEnabled         bool   `json:"network_enabled"`
	NetworkMode            string `json:"network_mode"`
	NetworkDomains         string `json:"network_domains"`
}

// ListProfiles returns all codex profiles
func (s *Service) ListProfiles() ([]Profile, error) {
	rows, err := s.db.Query(`SELECT name, display_name, api_key, base_url, default_model, 
		is_active, created_at, updated_at FROM codex_profiles WHERE api_key != '' ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []Profile
	for rows.Next() {
		var p Profile
		if err := rows.Scan(&p.Name, &p.DisplayName, &p.APIKey, &p.BaseURL,
			&p.DefaultModel, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

// GetActiveProfile returns the active profile
func (s *Service) GetActiveProfile() (*Profile, error) {
	var p Profile
	err := s.db.QueryRow(`SELECT name, display_name, api_key, base_url, default_model, 
		is_active, created_at, updated_at FROM codex_profiles WHERE is_active=1 AND api_key != ''`).
		Scan(&p.Name, &p.DisplayName, &p.APIKey, &p.BaseURL,
			&p.DefaultModel, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
