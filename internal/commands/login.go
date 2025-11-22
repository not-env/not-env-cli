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
// If url or apiKey are provided (non-empty), they will be used instead of prompting
func Login(url, apiKey string) error {
	reader := bufio.NewReader(os.Stdin)

	// Try to load existing config to get last used URL
	var defaultURL string
	existingConfig, err := config.Load()
	if err == nil && existingConfig.URL != "" {
		defaultURL = existingConfig.URL
	}

	// Get URL: use flag if provided, otherwise prompt
	if url == "" {
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
		url = strings.TrimSpace(urlInput)

		// Use default if empty
		if url == "" {
			if defaultURL == "" {
				return fmt.Errorf("URL is required")
			}
			url = defaultURL
		}
	}

	// Ensure URL has protocol
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	// Get API key: use flag if provided, otherwise prompt
	if apiKey == "" {
		fmt.Print("API Key: ")
		apiKeyInput, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read API key: %w", err)
		}
		apiKey = strings.TrimSpace(apiKeyInput)
	}

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

	// Get API key type from /me endpoint
	cfg := &config.Config{
		URL:    url,
		APIKey: apiKey,
	}

	meResp, err := cl.Get("/me")
	if err != nil {
		return fmt.Errorf("failed to get API key info: %w", err)
	}
	defer meResp.Body.Close()

	if meResp.StatusCode != 200 {
		return fmt.Errorf("failed to get API key info: status %d", meResp.StatusCode)
	}

	var meInfo struct {
		KeyType       string `json:"key_type"`
		EnvironmentID *int64 `json:"environment_id,omitempty"`
	}
	if err := client.ParseResponse(meResp, &meInfo); err != nil {
		return fmt.Errorf("failed to parse API key info: %w", err)
	}

	cfg.KeyType = meInfo.KeyType
	cfg.EnvIDFromKey = meInfo.EnvironmentID

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Logged in successfully! (Key type: %s)\n", cfg.KeyType)
	return nil
}

