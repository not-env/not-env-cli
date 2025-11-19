package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigSaveLoad(t *testing.T) {
	// Use a temporary directory for testing
	tmpDir := t.TempDir()
	originalPath := configPath
	configPath = filepath.Join(tmpDir, "config")
	defer func() {
		configPath = originalPath
	}()

	cfg := &Config{
		URL:    "https://test.example.com",
		APIKey: "test-api-key",
	}

	// Save config
	err := cfg.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loaded.URL != cfg.URL {
		t.Errorf("URL mismatch: got %q, want %q", loaded.URL, cfg.URL)
	}

	if loaded.APIKey != cfg.APIKey {
		t.Errorf("APIKey mismatch: got %q, want %q", loaded.APIKey, cfg.APIKey)
	}
}

func TestConfigLoadNotFound(t *testing.T) {
	// Use a temporary directory for testing
	tmpDir := t.TempDir()
	originalPath := configPath
	configPath = filepath.Join(tmpDir, "nonexistent", "config")
	defer func() {
		configPath = originalPath
	}()

	_, err := Load()
	if err == nil {
		t.Error("Expected error when config file doesn't exist")
	}
}

func TestConfigClear(t *testing.T) {
	// Use a temporary directory for testing
	tmpDir := t.TempDir()
	originalPath := configPath
	configPath = filepath.Join(tmpDir, "config")
	defer func() {
		configPath = originalPath
	}()

	cfg := &Config{
		URL:    "https://test.example.com",
		APIKey: "test-api-key",
	}

	// Save config
	err := cfg.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("Config file should exist: %v", err)
	}

	// Clear config
	err = Clear()
	if err != nil {
		t.Fatalf("Failed to clear config: %v", err)
	}

	// Verify file doesn't exist
	if _, err := os.Stat(configPath); err == nil {
		t.Error("Config file should not exist after clear")
	}

	// Clear again (should not error)
	err = Clear()
	if err != nil {
		t.Errorf("Clearing non-existent config should not error: %v", err)
	}
}

