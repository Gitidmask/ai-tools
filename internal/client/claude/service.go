package claude

import (
	"database/sql"
	"fmt"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// Preset represents a Claude preset configuration
type Preset struct {
	ID                        string `json:"id"`
	Name                      string `json:"name"`
	Color                     string `json:"color"`
	Icon                      string `json:"icon"`
	Provider                  string `json:"provider"`
	BaseURL                   string `json:"base_url"`
	DefaultModel              string `json:"default_model"`
	PresetDefaultOpusModel    string `json:"preset_default_opus_model"`
	PresetDefaultSonnetModel  string `json:"preset_default_sonnet_model"`
	PresetDefaultHaikuModel   string `json:"preset_default_haiku_model"`
	Badge                     string `json:"badge"`
	SortOrder                 int    `json:"sort_order"`
}

// ListPresets returns all Claude presets
func (s *Service) ListPresets() ([]Preset, error) {
	rows, err := s.db.Query(`SELECT id, name, color, icon, provider, base_url, 
		default_model, preset_default_opus_model, preset_default_sonnet_model, 
		preset_default_haiku_model, badge, sort_order 
		FROM claude_presets ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var presets []Preset
	for rows.Next() {
		var p Preset
		if err := rows.Scan(&p.ID, &p.Name, &p.Color, &p.Icon, &p.Provider,
			&p.BaseURL, &p.DefaultModel, &p.PresetDefaultOpusModel,
			&p.PresetDefaultSonnetModel, &p.PresetDefaultHaikuModel,
			&p.Badge, &p.SortOrder); err != nil {
			return nil, err
		}
		presets = append(presets, p)
	}
	return presets, nil
}

// ApplyPreset applies a preset configuration (writes env vars)
func (s *Service) ApplyPreset(presetID string) error {
	// This would set environment variables for Claude configuration
	return fmt.Errorf("not implemented")
}
