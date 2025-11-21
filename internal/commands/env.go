package commands

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"not-env-cli/internal/client"
	"not-env-cli/internal/config"
)

// validateEnvironmentName validates an environment name
func validateEnvironmentName(name string) error {
	if name == "" {
		return fmt.Errorf("environment name cannot be empty")
	}
	// Allow alphanumeric, dashes, underscores
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", name)
	if !matched {
		return fmt.Errorf("environment name must be alphanumeric with dashes/underscores only")
	}
	return nil
}

// EnvCreate creates a new environment
func EnvCreate(name, description string) error {
	if err := validateEnvironmentName(name); err != nil {
		return err
	}
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cl := client.NewClient(cfg.URL, cfg.APIKey)

	reqBody := map[string]interface{}{
		"name": name,
	}
	if description != "" {
		reqBody["description"] = description
	}

	resp, err := cl.Post("/environments", reqBody)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return client.ParseResponse(resp, nil)
	}

	var result struct {
		ID             int64  `json:"id"`
		OrganizationID int64  `json:"organization_id"`
		Name           string `json:"name"`
		Description    string `json:"description"`
		CreatedAt      string `json:"created_at"`
		Keys           struct {
			EnvAdmin    string `json:"env_admin"`
			EnvReadOnly string `json:"env_read_only"`
		} `json:"keys"`
	}

	if err := client.ParseResponse(resp, &result); err != nil {
		return err
	}

	fmt.Printf("Environment created successfully!\n")
	fmt.Printf("ID: %d\n", result.ID)
	fmt.Printf("Name: %s\n", result.Name)
	fmt.Printf("ENV_ADMIN key: %s\n", result.Keys.EnvAdmin)
	fmt.Printf("ENV_READ_ONLY key: %s\n", result.Keys.EnvReadOnly)
	fmt.Println("\nSave these keys securely!")

	return nil
}

// EnvList lists all environments
func EnvList() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cl := client.NewClient(cfg.URL, cfg.APIKey)

	resp, err := cl.Get("/environments")
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return client.ParseResponse(resp, nil)
	}

	var result struct {
		Environments []struct {
			ID             int64  `json:"id"`
			OrganizationID int64  `json:"organization_id"`
			Name           string `json:"name"`
			Description    string `json:"description"`
			CreatedAt      string `json:"created_at"`
			UpdatedAt      string `json:"updated_at"`
		} `json:"environments"`
	}

	if err := client.ParseResponse(resp, &result); err != nil {
		return err
	}

	if len(result.Environments) == 0 {
		fmt.Println("No environments found.")
		return nil
	}

	fmt.Println("Environments:")
	for _, env := range result.Environments {
		fmt.Printf("  ID: %d, Name: %s", env.ID, env.Name)
		if env.Description != "" {
			fmt.Printf(", Description: %s", env.Description)
		}
		fmt.Println()
	}

	return nil
}

// EnvDelete deletes an environment
func EnvDelete(envID int64) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cl := client.NewClient(cfg.URL, cfg.APIKey)

	resp, err := cl.Delete(fmt.Sprintf("/environments/%d", envID))
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return client.ParseResponse(resp, nil)
	}

	fmt.Printf("Environment %d deleted successfully!\n", envID)
	return nil
}

// EnvImport imports variables from a .env file.
// Flow:
//  1. Parses .env file into key-value pairs
//  2. Creates environment via API (or uses existing if overwrite=true)
//  3. Gets both ENV_ADMIN and ENV_READ_ONLY keys from API response
//  4. Switches to ENV_ADMIN key for setting variables
//  5. Sets all variables from .env file
//  6. Outputs both keys for user (ENV_ADMIN for CLI, ENV_READ_ONLY for SDKs)
func EnvImport(name, description, filePath string, overwrite bool) error {
	if err := validateEnvironmentName(name); err != nil {
		return err
	}
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Read .env file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Parse .env file
	// Supports: KEY=VALUE, KEY="value", KEY='value', comments (#), empty lines
	envVars := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first = to handle values that contain =
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove surrounding quotes (both single and double) if present
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		envVars[key] = value
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Create environment (or use existing if overwrite)
	cl := client.NewClient(cfg.URL, cfg.APIKey)

	var envAdminKey string

	if overwrite {
		// Try to find existing environment
		resp, err := cl.Get("/environments")
		if err == nil && resp.StatusCode == 200 {
			var envs struct {
				Environments []struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
				} `json:"environments"`
			}
			if err := client.ParseResponse(resp, &envs); err == nil {
				for _, env := range envs.Environments {
					if env.Name == name {
						// Need to get keys - but we can't retrieve them via API
						// Provide clear instructions for user
						fmt.Printf("Environment '%s' already exists. To import variables:\n", name)
						fmt.Printf("  1. Run: not-env login\n")
						fmt.Printf("  2. Enter the ENV_ADMIN key for '%s'\n", name)
						fmt.Printf("  3. Run: not-env env import --name %s --file %s --overwrite\n", name, filePath)
						return fmt.Errorf("environment exists - login with ENV_ADMIN key first")
					}
				}
			}
		}
	}

	// Create new environment
	reqBody := map[string]interface{}{
		"name": name,
	}
	if description != "" {
		reqBody["description"] = description
	}

	resp, err := cl.Post("/environments", reqBody)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return client.ParseResponse(resp, nil)
	}

	var createResult struct {
		ID   int64 `json:"id"`
		Keys struct {
			EnvAdmin    string `json:"env_admin"`
			EnvReadOnly string `json:"env_read_only"`
		} `json:"keys"`
	}

	if err := client.ParseResponse(resp, &createResult); err != nil {
		return err
	}

	envAdminKey = createResult.Keys.EnvAdmin
	envReadOnlyKey := createResult.Keys.EnvReadOnly

	// Switch to ENV_ADMIN key for setting variables
	cl = client.NewClient(cfg.URL, envAdminKey)

	// Set all variables
	for key, value := range envVars {
		setResp, err := cl.Put(fmt.Sprintf("/variables/%s", key), map[string]interface{}{
			"value": value,
		})
		if err != nil {
			fmt.Printf("Warning: Failed to set %s: %v\n", key, err)
			continue
		}
		setResp.Body.Close()
		if setResp.StatusCode != 204 {
			fmt.Printf("Warning: Failed to set %s (status %d)\n", key, setResp.StatusCode)
		}
	}

	fmt.Printf("Environment '%s' created and populated with %d variables!\n", name, len(envVars))
	fmt.Printf("ENV_ADMIN key: %s\n", envAdminKey)
	fmt.Printf("ENV_READ_ONLY key: %s\n", envReadOnlyKey)
	fmt.Println("\nSave these keys securely!")
	fmt.Println("Use ENV_ADMIN for managing variables (CLI)")
	fmt.Println("Use ENV_READ_ONLY for applications (SDKs)")

	return nil
}

// EnvShow shows current environment metadata
func EnvShow() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cl := client.NewClient(cfg.URL, cfg.APIKey)

	resp, err := cl.Get("/environment")
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return client.ParseResponse(resp, nil)
	}

	var env struct {
		ID             int64  `json:"id"`
		OrganizationID int64  `json:"organization_id"`
		Name           string `json:"name"`
		Description    string `json:"description"`
		CreatedAt      string `json:"created_at"`
		UpdatedAt      string `json:"updated_at"`
	}

	if err := client.ParseResponse(resp, &env); err != nil {
		return err
	}

	fmt.Printf("Environment: %s (ID: %d)\n", env.Name, env.ID)
	if env.Description != "" {
		fmt.Printf("Description: %s\n", env.Description)
	}
	fmt.Printf("Created: %s\n", env.CreatedAt)
	fmt.Printf("Updated: %s\n", env.UpdatedAt)

	return nil
}

// EnvUpdate updates environment metadata
func EnvUpdate(name, description *string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cl := client.NewClient(cfg.URL, cfg.APIKey)

	reqBody := make(map[string]interface{})
	if name != nil {
		reqBody["name"] = *name
	}
	if description != nil {
		reqBody["description"] = *description
	}

	resp, err := cl.Patch("/environment", reqBody)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return client.ParseResponse(resp, nil)
	}

	fmt.Println("Environment updated successfully!")
	return nil
}

// EnvKeys shows environment keys
func EnvKeys() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cl := client.NewClient(cfg.URL, cfg.APIKey)

	resp, err := cl.Get("/environment/keys")
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return client.ParseResponse(resp, nil)
	}

	var keys struct {
		EnvAdmin    string `json:"env_admin"`
		EnvReadOnly string `json:"env_read_only"`
	}

	if err := client.ParseResponse(resp, &keys); err != nil {
		return err
	}

	fmt.Println("Environment Keys:")
	fmt.Printf("ENV_ADMIN: %s\n", keys.EnvAdmin)
	fmt.Printf("ENV_READ_ONLY: %s\n", keys.EnvReadOnly)

	return nil
}

// EnvSet prints export commands for all variables
func EnvSet() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cl := client.NewClient(cfg.URL, cfg.APIKey)

	resp, err := cl.Get("/variables")
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return client.ParseResponse(resp, nil)
	}

	var result struct {
		Variables []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"variables"`
	}

	if err := client.ParseResponse(resp, &result); err != nil {
		return err
	}

	for _, v := range result.Variables {
		// Escape value for shell
		escapedValue := strings.ReplaceAll(v.Value, `"`, `\"`)
		escapedValue = strings.ReplaceAll(escapedValue, `$`, `\$`)
		escapedValue = strings.ReplaceAll(escapedValue, "`", "\\`")
		fmt.Printf("export %s=\"%s\"\n", v.Key, escapedValue)
	}

	return nil
}

// EnvClear prints unset commands for all variables
func EnvClear() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cl := client.NewClient(cfg.URL, cfg.APIKey)

	resp, err := cl.Get("/variables")
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return client.ParseResponse(resp, nil)
	}

	var result struct {
		Variables []struct {
			Key string `json:"key"`
		} `json:"variables"`
	}

	if err := client.ParseResponse(resp, &result); err != nil {
		return err
	}

	for _, v := range result.Variables {
		fmt.Printf("unset %s\n", v.Key)
	}

	return nil
}

