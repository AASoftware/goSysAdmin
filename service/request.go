package service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Web struct {
	URL       string `json:"url"`
	CertCheck bool   `json:"certcheck"`
	IsUp      bool   `json:"website"`
	WebErrors []string
}

// NewWebChecker creates a new instance of Web struct.
func NewWebChecker() (*Web, error) {
	return &Web{}, nil
}

// LoadConfig loads configuration from a file.
func (s *Web) LoadConfig(file string) {
	f, err := os.ReadFile(file)
	if err != nil {
		panic("Config can't be loaded")
	}

	err = json.Unmarshal(f, &s)
	if err != nil {
		panic("Failed to unmarshal config.json")
	}
	fmt.Println("Loaded URL:", s.URL)
}

// CheckWebsite checks if the website is reachable.
func (w *Web) CheckWebsite() {
	resp, err := http.Get(w.URL)
	if err != nil {
		fmt.Println("Failed to check if Website is up")
		return
	}
	if resp.StatusCode != http.StatusOK {
		w.WebErrors = append(w.WebErrors, fmt.Sprintf("%s: not reachable (Status: %s)", w.URL, resp.Status))
	}
}

// CheckCert checks the SSL certificate of the website.
func (w *Web) CheckCert() {
	u, err := url.Parse(w.URL)
	if err != nil {
		fmt.Println("Failed to parse URL:", err)
		return
	}

	conn, err := tls.Dial("tcp", u.Host+":443", &tls.Config{
		InsecureSkipVerify: true, // Skip certificate verification
	})
	if err != nil {
		fmt.Println("Failed to establish TLS connection:", err)
		return
	}
	defer conn.Close()

	// Get the server's certificate
	state := conn.ConnectionState()
	cert := state.PeerCertificates[0]

	// Get the certificate's expiry date
	expiryDate := cert.NotAfter

	// Get the current date and time
	currentTime := time.Now()

	// Check if the certificate expires in less than 3 months
	if expiryDate.Before(currentTime.AddDate(0, 3, 0)) {
		w.WebErrors = append(w.WebErrors, fmt.Sprintf("%s: certificate expires in less than 3 months (expiry date: %s)", w.URL, expiryDate))
	}
}

// ResetWebErrors resets the WebErrors slice.
func (w *Web) ResetWebErrors() {
	w.WebErrors = []string{} // Assign an empty slice
}
