package service

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ricochet2200/go-disk-usage/du"
)

type SysState struct {
	SysEnabled   bool     `json:"sysenabled"`
	Disks        []string `json:"disks"`
	LowDiskSpace []string
}

// NewSysChecker creates a new instance of SysState struct.
func NewSysChecker() (*SysState, error) {
	return &SysState{}, nil
}

// LoadConfig loads configuration from a file.
func (s *SysState) LoadConfig(file string) {
	f, err := os.ReadFile(file)
	if err != nil {
		panic("Config can't be loaded")
	}

	err = json.Unmarshal(f, &s)
	if err != nil {
		panic("Failed to unmarshal config.json")
	}
}

// CheckDisks checks the disk space of the specified disks.
func (s *SysState) CheckDisks() {
	for _, dsk := range s.Disks {
		d := dsk
		disk := du.NewDiskUsage(d)
		free := disk.Available()

		gigabytes := float64(free) / (1024 * 1024 * 1024)
		if gigabytes < 20.0 {
			fmt.Printf("%s, %.2f, GB\n", d, gigabytes)
			s.LowDiskSpace = append(s.LowDiskSpace, fmt.Sprintf("Low Disk Space Drive:%s, free:%f", d, gigabytes))
		}
	}
}

// ResetLowDiskSpace resets the LowDiskSpace slice.
func (s *SysState) ResetLowDiskSpace() {
	s.LowDiskSpace = []string{} // Assign an empty slice
}
