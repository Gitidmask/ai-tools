package devbridge

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

// Server provides the DevBridge HTTP API that the frontend and sidecar communicate with.
// This runs inside the main Go process (ai_tools.exe), not in the sidecar subprocess.
type Server struct {
	port       int
	server     *http.Server
	listener   net.Listener
	startedAt  time.Time
}

// New creates a new DevBridge server
func New(port int) *Server {
	return &Server{
		port:      port,
		startedAt: time.Now(),
	}
}

// Start begins listening on the configured port
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/health", s.withCORS(s.handleHealth))
	mux.HandleFunc("/api/sidecar-config", s.withCORS(s.handleSidecarConfig))

	// Try the configured port, fall back to a random port if busy
	addr := fmt.Sprintf("127.0.0.1:%d", s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		// Port busy, try random
		listener, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return fmt.Errorf("devbridge: failed to listen: %w", err)
		}
	}
	s.listener = listener
	s.port = listener.Addr().(*net.TCPAddr).Port

	s.server = &http.Server{
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go s.server.Serve(listener)
	return nil
}

// Stop shuts down the server
func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

// Port returns the actual listening port
func (s *Server) Port() int {
	return s.port
}

// HealthHandler is a standalone HTTP handler for health checks
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
	})
}

// ConfigHandler is a standalone HTTP handler for config info
func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"version": "1.0.0",
	})
}

// handleHealth returns the health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
		"pid":       fmt.Sprintf("%d", 0), // simplified
	})
}

// handleSidecarConfig returns configuration info
func (s *Server) handleSidecarConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"port":    s.port,
		"version": "1.0.0",
		"uptime":  time.Since(s.startedAt).String(),
	})
}

// withCORS wraps a handler with CORS headers
func (s *Server) withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
