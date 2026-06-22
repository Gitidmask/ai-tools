package env

type Service struct{}

func NewService() *Service {
	return &Service{}
}

type EnvironmentInfo struct {
	OS      string `json:"os"`
	Arch    string `json:"arch"`
	Debug   bool   `json:"debug"`
	GoVersion string `json:"go_version"`
	GoOS    string `json:"go_os"`
	GoArch  string `json:"go_arch"`
}

func (s *Service) GetEnvironment() EnvironmentInfo {
	return EnvironmentInfo{
		OS:    "windows",
		Arch:  "amd64",
		Debug: false,
	}
}
