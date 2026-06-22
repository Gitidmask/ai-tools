package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Service struct {
	client *http.Client
}

func NewService() *Service {
	return &Service{
		client: &http.Client{},
	}
}

type DownloadProgress struct {
	Total     int64 `json:"total"`
	Completed int64 `json:"completed"`
	Speed     int64 `json:"speed"`
}

func (s *Service) DownloadFile(url, destDir, filename string) error {
	destPath := filepath.Join(destDir, filename)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	resp, err := s.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
