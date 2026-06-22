package clihistory

import (
	"database/sql"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

type CLISession struct {
	Source      string `json:"source"`
	SessionID   string `json:"session_id"`
	CWD         string `json:"cwd"`
	FilePath    string `json:"file_path"`
	StartedAt   int64  `json:"started_at"`
	LastActive  int64  `json:"last_active"`
	FirstMsg    string `json:"first_msg"`
	IndexedAt   int64  `json:"indexed_at"`
	DisplayName string `json:"display_name"`
}

func (s *Service) ListSessions() ([]CLISession, error) {
	rows, err := s.db.Query(`SELECT source, session_id, cwd, file_path, started_at, 
		last_active, first_msg, indexed_at, display_name 
		FROM cli_sessions ORDER BY last_active DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []CLISession
	for rows.Next() {
		var sess CLISession
		if err := rows.Scan(&sess.Source, &sess.SessionID, &sess.CWD, &sess.FilePath,
			&sess.StartedAt, &sess.LastActive, &sess.FirstMsg, &sess.IndexedAt, &sess.DisplayName); err != nil {
			return nil, err
		}
		sessions = append(sessions, sess)
	}
	return sessions, nil
}
