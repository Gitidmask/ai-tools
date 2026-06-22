package sidecar

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Service struct {
	cmd       *exec.Cmd
	db        *sql.DB
	running   bool
	port      int
	healthURL string
}

type SidecarStatus struct {
	Running bool   `json:"running"`
	PID     int    `json:"pid"`
	Port    int    `json:"port"`
	Version string `json:"version"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
	Version   string `json:"version"`
	PID       int    `json:"pid"`
}

func NewService() *Service {
	return &Service{}
}

// SetDB sets the database instance for the sidecar service
func (s *Service) SetDB(db *sql.DB) {
	s.db = db
}

// Start launches the sidecar process
func (s *Service) Start(port int) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	sidecarDir := filepath.Join(filepath.Dir(exe), "sidecars")
	sidecarPath := filepath.Join(sidecarDir, "claude-sidecar-x86_64-pc-windows-msvc.exe")

	// Check multiple possible locations
	if _, err := os.Stat(sidecarPath); os.IsNotExist(err) {
		// Try relative to working directory
		sidecarPath = "sidecars/claude-sidecar-x86_64-pc-windows-msvc.exe"
		if _, err := os.Stat(sidecarPath); os.IsNotExist(err) {
			return fmt.Errorf("sidecar binary not found (tried: %s, %s)",
				filepath.Join(sidecarDir, "claude-sidecar-x86_64-pc-windows-msvc.exe"),
				sidecarPath)
		}
	}

	// Set port: use provided or find free port
	if port <= 0 {
		port = 9800
	}

	s.cmd = exec.Command(sidecarPath)
	s.cmd.Env = append(os.Environ(),
		fmt.Sprintf("SIDECAR_PORT=%d", port),
		"LISTEN_ADDR=127.0.0.1",
		"NODE_ENV=production",
	)

	// Capture output for logging
	logDir := filepath.Join(filepath.Dir(exe), "logs")
	os.MkdirAll(logDir, 0755)
	logPath := filepath.Join(logDir, "sidecar.log")

	logFile, err := os.Create(logPath)
	if err == nil {
		s.cmd.Stdout = io.MultiWriter(os.Stdout, logFile)
		s.cmd.Stderr = io.MultiWriter(os.Stderr, logFile)
	} else {
		s.cmd.Stdout = os.Stdout
		s.cmd.Stderr = os.Stderr
	}

	if err := s.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start sidecar: %w", err)
	}

	s.running = true
	s.port = port
	s.healthURL = fmt.Sprintf("http://127.0.0.1:%d/api/health", port)

	// Wait for health check
	for i := 0; i < 30; i++ {
		if s.checkHealth() {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}

	return fmt.Errorf("sidecar started but health check failed after 6s")
}

// Stop terminates the sidecar process
func (s *Service) Stop() error {
	if s.cmd != nil && s.cmd.Process != nil {
		if err := s.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill sidecar: %w", err)
		}
		// Wait for process to exit
		s.cmd.Wait()
	}
	s.running = false
	return nil
}

// Status returns the current sidecar status
func (s *Service) Status() SidecarStatus {
	status := SidecarStatus{
		Running: s.running,
		Port:    s.port,
		Version: "1.0.0",
	}
	if s.cmd != nil && s.cmd.Process != nil {
		status.PID = s.cmd.Process.Pid
	}
	return status
}

// GetConfig returns the sidecar configuration via API
func (s *Service) GetConfig() (map[string]interface{}, error) {
	if !s.running || s.healthURL == "" {
		return nil, fmt.Errorf("sidecar not running")
	}

	configURL := fmt.Sprintf("http://127.0.0.1:%d/api/sidecar-config", s.port)
	resp, err := http.Get(configURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get sidecar config: %w", err)
	}
	defer resp.Body.Close()

	var config map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}
	return config, nil
}

// checkHealth performs a health check on the sidecar
func (s *Service) checkHealth() bool {
	if s.healthURL == "" {
		return false
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(s.healthURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return false
	}

	return health.Status == "ok"
}

// HealthCheck performs a health check (exported for external use)
func (s *Service) HealthCheck() bool {
	return s.checkHealth()
}
