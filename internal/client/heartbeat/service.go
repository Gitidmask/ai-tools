package heartbeat

type Service struct{}

func NewService() *Service {
	return &Service{}
}

type HealthStatus struct {
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
	Version   string `json:"version"`
}

func (s *Service) Ping() HealthStatus {
	return HealthStatus{
		Status:    "ok",
		Timestamp: 0,
		Version:   "1.0.0",
	}
}
