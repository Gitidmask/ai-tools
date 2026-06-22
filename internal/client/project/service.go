package project

import (
	"database/sql"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	CreatedAt int64  `json:"created_at"`
	LastUsed  int64  `json:"last_used"`
	Pinned    bool   `json:"pinned"`
}

func (s *Service) ListProjects() ([]Project, error) {
	rows, err := s.db.Query(`SELECT id, name, path, created_at, last_used, pinned 
		FROM projects ORDER BY pinned DESC, last_used DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Path, &p.CreatedAt, &p.LastUsed, &p.Pinned); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}
