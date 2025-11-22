package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"not-env-cli/internal/client"
	"not-env-cli/internal/config"
)

// Use switches to a different API key while keeping the same backend URL
func Use() error {
	// Load existing config to get URL
	existingConfig, err := config.Load()
	if err != nil {
		return fmt.Errorf("not logged in. Run 'not-env login' first to set backend URL")
	}

	if existingConfig.URL == "" {
		return fmt.Errorf("no backend URL configured. Run 'not-env login' first")
	}

	reader := bufio.NewReader(os.Stdin)
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
	cl := client.NewClient(existingConfig.URL, apiKey)
	resp, err := cl.Get("/health")
	if err != nil {
		return fmt.Errorf("failed to connect to backend: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid API key or backend unreachable")
	}

	// Get API key type from /me endpoint
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

	// Update config with new API key (keep existing URL)
	existingConfig.APIKey = apiKey
	existingConfig.KeyType = meInfo.KeyType
	existingConfig.EnvIDFromKey = meInfo.EnvironmentID
	// Clear env_id since we're switching to a potentially different environment
	existingConfig.EnvID = nil

	if err := existingConfig.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Switched to new API key (backend: %s, key type: %s)\n", existingConfig.URL, existingConfig.KeyType)
	return nil
}

