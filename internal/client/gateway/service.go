package gateway

import (
	"database/sql"
	"net/http"
	"strings"
)

type Service struct {
	db     *sql.DB
	server *http.Server
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

type GatewayConfig struct {
	ListenAddr       string `json:"listen_addr"`
	RoutingStrategy  string `json:"routing_strategy"`
	MaxConcurrency   int    `json:"max_concurrency"`
	SessionAffinity  bool   `json:"session_affinity"`
	TargetBaseURL    string `json:"target_base_url"`
	UpstreamProxyURL string `json:"upstream_proxy_url"`
	APIKey           string `json:"api_key"`
}

func (s *Service) GetConfig() (*GatewayConfig, error) {
	var cfg GatewayConfig
	err := s.db.QueryRow(`SELECT listen_addr, routing_strategy, max_concurrency, 
		session_affinity, target_base_url, upstream_proxy_url, api_key 
		FROM gateway_config WHERE id = 1`).
		Scan(&cfg.ListenAddr, &cfg.RoutingStrategy, &cfg.MaxConcurrency,
			&cfg.SessionAffinity, &cfg.TargetBaseURL, &cfg.UpstreamProxyURL, &cfg.APIKey)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// StartProxy starts the AI gateway proxy server
func (s *Service) StartProxy() error {
	cfg, err := s.GetConfig()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/", s.proxyHandler)

	s.server = &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: mux,
	}

	go s.server.ListenAndServe()
	return nil
}

func (s *Service) proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Forward requests to upstream provider
	cfg, err := s.GetConfig()
	if err != nil {
		http.Error(w, "gateway not configured", http.StatusInternalServerError)
		return
	}

	targetURL := strings.TrimRight(cfg.TargetBaseURL, "/") + r.URL.Path
	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "failed to create proxy request", http.StatusInternalServerError)
		return
	}

	// Forward headers
	for k, v := range r.Header {
		proxyReq.Header[k] = v
	}
	if cfg.APIKey != "" {
		proxyReq.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	}

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "proxy request failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	// Copy body
	// ...
}
