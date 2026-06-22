package npm

import (
	"fmt"
	"os/exec"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

type NPMStatus struct {
	Installed bool   `json:"installed"`
	Version   string `json:"version"`
	Path      string `json:"path"`
}

func (s *Service) CheckNPM() NPMStatus {
	cmd := exec.Command("npm", "--version")
	output, err := cmd.Output()
	if err != nil {
		return NPMStatus{Installed: false}
	}
	return NPMStatus{
		Installed: true,
		Version:   string(output),
	}
}

func (s *Service) InstallPackage(name string) error {
	cmd := exec.Command("npm", "install", "-g", name)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("npm install failed: %s: %w", string(output), err)
	}
	return nil
}
