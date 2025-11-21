package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"not-env-cli/internal/client"
	"not-env-cli/internal/config"
)

// Login handles the login command
func Login() error {
	reader := bufio.NewReader(os.Stdin)

	// Try to load existing config to get last used URL
	var defaultURL string
	existingConfig, err := config.Load()
	if err == nil && existingConfig.URL != "" {
		defaultURL = existingConfig.URL
	}

	// Prompt for URL with default
	if defaultURL != "" {
		fmt.Printf("Backend URL [%s]: ", defaultURL)
	} else {
		fmt.Print("Backend URL: ")
	}
	urlInput, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read URL: %w", err)
	}
	url := strings.TrimSpace(urlInput)

	// Use default if empty
	if url == "" {
		if defaultURL == "" {
			return fmt.Errorf("URL is required")
		}
		url = defaultURL
	}

	// Ensure URL has protocol
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	fmt.Print("API Key: ")
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read API key: %w", err)
	}
	apiKey = strings.TrimSpace(apiKey)

	if apiKey == "" {
		return fmt.Errorf("API key is required")
	}

	// Validate credentials by making a test request
	cl := client.NewClient(url, apiKey)
	resp, err := cl.Get("/health")
	if err != nil {
		return fmt.Errorf("failed to connect to backend: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid credentials or backend unreachable")
	}

	// Save config
	cfg := &config.Config{
		URL:    url,
		APIKey: apiKey,
	}

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("Logged in successfully!")
	return nil
}

