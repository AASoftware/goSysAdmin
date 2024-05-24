package service

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/sys/windows/svc/mgr"
)

type Service struct {
	ServiceList []string `json:"services"`
	Enabled     bool     `json:"enabled"`
	NotRunning  []string
}

// NewServiceChecker creates a new instance of Service struct.
func NewServiceChecker() (*Service, error) {
	return &Service{}, nil
}

// LoadConfig loads configuration from a file.
func (s *Service) LoadConfig(file string) {
	f, err := os.ReadFile(file)
	if err != nil {
		panic("Config can't be loaded")
	}

	err = json.Unmarshal(f, &s)
	if err != nil {
		panic("Failed to unmarshal config.json")
	}
}

// CheckStatus checks the status of specified services.
func (s *Service) CheckStatus() error {
	m, err := mgr.Connect()
	if err != nil {
		fmt.Println("Failed to connect with Service Manager")
		return err
	}
	defer m.Disconnect() // Close connection to Service Manager

	for _, srvName := range s.ServiceList {
		srv, err := m.OpenService(srvName)
		if err != nil {
			fmt.Printf("Failed to open service %s: %v\n", srvName, err)
			continue
		}
		defer srv.Close() // Close service

		status, err := srv.Query()
		if err != nil {
			fmt.Printf("Failed to query service %s: %v\n", srvName, err)
			continue
		}
		fmt.Printf("Service %s status: %v\n", srvName, status.State)
		if status.State != 4 {
			s.NotRunning = append(s.NotRunning, srvName)
		}
	}
	return nil
}

// ResetNotRunning resets the NotRunning slice.
func (s *Service) ResetNotRunning() {
	s.NotRunning = []string{} // Assign an empty slice
}
