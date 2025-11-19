package commands

import (
	"fmt"

	"not-env-cli/internal/config"
)

// Logout handles the logout command
func Logout() error {
	if err := config.Clear(); err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	fmt.Println("Logged out successfully!")
	return nil
}

