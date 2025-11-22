// Package main provides the not-env CLI tool.
//
// The CLI is a command-line interface for managing environments and variables
// stored in the not-env backend. It uses Cobra for command structure and
// communicates with the backend via HTTPS.
//
// Command structure:
//   - Authentication: login, logout, use
//   - Environment management: env create/list/delete/import/show/update/keys/set/clear
//   - Variable management: var list/get/set/delete
//
// Configuration is stored in ~/.not-env/config (created via login command).
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"not-env-cli/internal/commands"
)

var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:     "not-env",
	Short:   "not-env CLI - Manage environment variables",
	Long:    "not-env is a CLI tool for managing environment variables stored in not-env-backend",
	Version: version,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to not-env backend",
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := cmd.Flags().GetString("url")
		apiKey, _ := cmd.Flags().GetString("api-key")
		return commands.Login(url, apiKey)
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from not-env backend",
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.Logout()
	},
}

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Switch to a different API key (keeps current backend URL)",
	Long:  "Switch to a different API key while keeping the same backend URL. Useful for switching between environments or API key types (APP_ADMIN, ENV_ADMIN, ENV_READ_ONLY).",
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.Use()
	},
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environments",
}

var envCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new environment (APP_ADMIN)",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		if name == "" {
			return fmt.Errorf("--name is required")
		}

		return commands.EnvCreate(name, description)
	},
}

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all environments (APP_ADMIN, ENV_ADMIN, ENV_READ_ONLY)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.EnvList()
	},
}

var envDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an environment (APP_ADMIN)",
	RunE: func(cmd *cobra.Command, args []string) error {
		envID, _ := cmd.Flags().GetInt64("id")
		if envID == 0 {
			return fmt.Errorf("--id is required")
		}
		return commands.EnvDelete(envID)
	},
}

var envImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import variables from a .env file (ENV_ADMIN)",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		filePath, _ := cmd.Flags().GetString("file")
		overwrite, _ := cmd.Flags().GetBool("overwrite")

		// Name validation is handled in EnvImport based on key type
		if filePath == "" {
			return fmt.Errorf("--file is required")
		}

		return commands.EnvImport(name, description, filePath, overwrite)
	},
}

var envShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current environment metadata (APP_ADMIN, ENV_ADMIN, ENV_READ_ONLY)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.EnvShow()
	},
}

var envUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update environment metadata (APP_ADMIN, ENV_ADMIN)",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		var namePtr, descPtr *string
		if name != "" {
			namePtr = &name
		}
		if description != "" {
			descPtr = &description
		}

		if namePtr == nil && descPtr == nil {
			return fmt.Errorf("at least one of --name or --description is required")
		}

		return commands.EnvUpdate(namePtr, descPtr)
	},
}

var envKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Show environment API keys (ENV_ADMIN)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.EnvKeys()
	},
}

var envSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Print export commands for all variables (ENV_ADMIN, ENV_READ_ONLY)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.EnvSet()
	},
}

var envClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Print unset commands for all variables (ENV_ADMIN, ENV_READ_ONLY)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.EnvClear()
	},
}

var varCmd = &cobra.Command{
	Use:   "var",
	Short: "Manage variables",
}

var varListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all variables (ENV_ADMIN, ENV_READ_ONLY)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.VarList()
	},
}

var varGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a variable value (ENV_ADMIN, ENV_READ_ONLY)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.VarGet(args[0])
	},
}

var varSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a variable value (ENV_ADMIN)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.VarSet(args[0], args[1])
	},
}

var varDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a variable (ENV_ADMIN)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.VarDelete(args[0])
	},
}

func init() {
	// Login/logout/use
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(useCmd)

	// Login command flags
	loginCmd.Flags().String("url", "", "Backend URL (optional, will prompt if not provided)")
	loginCmd.Flags().String("api-key", "", "API key (optional, will prompt if not provided)")

	// Environment commands
	rootCmd.AddCommand(envCmd)
	envCmd.AddCommand(envCreateCmd)
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envDeleteCmd)
	envCmd.AddCommand(envImportCmd)
	envCmd.AddCommand(envShowCmd)
	envCmd.AddCommand(envUpdateCmd)
	envCmd.AddCommand(envKeysCmd)
	envCmd.AddCommand(envSetCmd)
	envCmd.AddCommand(envClearCmd)

	envCreateCmd.Flags().String("name", "", "Environment name")
	envCreateCmd.Flags().String("description", "", "Environment description")
	envDeleteCmd.Flags().Int64("id", 0, "Environment ID")
	envImportCmd.Flags().String("name", "", "Environment name")
	envImportCmd.Flags().String("description", "", "Environment description")
	envImportCmd.Flags().String("file", "", "Path to .env file")
	envImportCmd.Flags().Bool("overwrite", false, "Overwrite existing environment")
	envUpdateCmd.Flags().String("name", "", "New environment name")
	envUpdateCmd.Flags().String("description", "", "New environment description")

	// Variable commands
	rootCmd.AddCommand(varCmd)
	varCmd.AddCommand(varListCmd)
	varCmd.AddCommand(varGetCmd)
	varCmd.AddCommand(varSetCmd)
	varCmd.AddCommand(varDeleteCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
