# not-env-cli

Command-line interface for not-env, a self-hosted environment variable management system. Manage environments and variables, and integrate with your shell.

## Common Tasks

| Task | Command |
|------|---------|
| **Login** | `not-env login` |
| **Create environment** | `not-env env create --name dev` |
| **Import .env file** | `not-env env import --name dev --file .env` |
| **List environments** | `not-env env list` |
| **Show current environment** | `not-env env show` |
| **Set variable** | `not-env var set DB_HOST localhost` |
| **Get variable** | `not-env var get DB_HOST` |
| **List variables** | `not-env var list` |
| **Load into shell** | `eval "$(not-env env set)"` |
| **Clear from shell** | `eval "$(not-env env clear)"` |

## Overview

The CLI allows you to:
- Manage environments (create, list, delete)
- Import variables from .env files
- Manage variables (set, get, list, delete)
- Load variables into your shell session
- View environment metadata and API keys

## Installation

**Linux/macOS (one command):**
```bash
curl -L https://github.com/not-env/not-env-cli/releases/latest/download/not-env-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m) -o /usr/local/bin/not-env && chmod +x /usr/local/bin/not-env
```

**Windows:** Download from [Releases page](https://github.com/not-env/not-env-cli/releases/latest)

**From Source:**
```bash
git clone https://github.com/not-env/not-env-cli.git && cd not-env-cli && go build -o not-env && sudo mv not-env /usr/local/bin/
```

## Environment Variables

The CLI does not require any environment variables. Configuration is stored in `~/.not-env/config` (created via `not-env login`).

## Quick Start

### 1. Login
```bash
not-env login
# Enter backend URL (defaults to last used)
# Enter API key (APP_ADMIN for creating environments)
```

### 2. Import .env File
```bash
# Create a sample .env file (or use your existing one)
cat > .env <<EOF
DB_HOST=localhost
DB_PORT=5432
API_KEY=your-secret-key
EOF

# Import .env file (creates environment AND imports all variables)
not-env env import --name dev --file .env
# Save both keys! ENV_READ_ONLY for SDKs, ENV_ADMIN for CLI.
```

**That's it!** The import command creates the environment and imports all variables in one step.

## Common Workflows

### Which Key Should I Use?

- **ENV_ADMIN**: Use this key with the CLI to manage variables (set, get, delete, import)
- **ENV_READ_ONLY**: Use this key in your applications (SDKs) - this is what you'll set as `NOT_ENV_API_KEY`
- **APP_ADMIN**: Use this key to create/manage environments (only needed for initial setup)

### Manual Variable Setting (Alternative)
```bash
# After creating environment with 'not-env env create --name dev'
not-env login  # Use ENV_ADMIN key
not-env var set DB_HOST "localhost"
not-env var set DB_PORT "5432"
```

### Load Variables into Shell
```bash
eval "$(not-env env set)"
```

### List All Variables
```bash
not-env var list
```

## Commands Reference

### Authentication

- `not-env login` - Login to backend (prompts for URL and API key)
- `not-env logout` - Clear saved credentials

### Environment Management

- `not-env env create --name NAME [--description DESC]` - Create environment (APP_ADMIN)
- `not-env env list` - List all environments (APP_ADMIN)
- `not-env env delete --id ENV_ID` - Delete environment (APP_ADMIN)
- `not-env env import --name NAME --file PATH [--overwrite]` - Import from .env file
- `not-env env show` - Show current environment metadata
- `not-env env update [--name NAME] [--description DESC]` - Update environment (ENV_ADMIN)
- `not-env env keys` - Show API keys for current environment (ENV_ADMIN)
- `not-env env set` - Print `export` commands (use with `eval`)
- `not-env env clear` - Print `unset` commands (use with `eval`)

### Variable Management

- `not-env var list` - List all variables
- `not-env var get KEY` - Get variable value
- `not-env var set KEY VALUE` - Set variable (ENV_ADMIN)
- `not-env var delete KEY` - Delete variable (ENV_ADMIN)

## Configuration

Configuration stored in `~/.not-env/config` (TOML format):

```toml
url = "https://not-env.example.com"
api_key = "your-api-key-here"
```

## Troubleshooting

**Wrong key type:**
- Use APP_ADMIN for environment management
- Use ENV_ADMIN for variable management
- Use ENV_READ_ONLY for read-only access

**Not logged in:**
- Run `not-env login` to authenticate

**Variables not loading in shell:**
- Use `eval "$(not-env env set)"` with quotes
- Verify you're logged in with ENV_* key
- Check variables exist: `not-env var list`

**Import fails:**
- Ensure .env file exists and is readable
- Check format (KEY=VALUE, one per line)
- Verify ENV_ADMIN permissions

## Integration

- **Backend**: Communicates with [not-env-backend](../not-env-backend/README.md) via HTTPS
- **SDKs**: Variables can be used with [JavaScript SDK](../SDKs/not-env-sdk-js/README.md) or [Python SDK](../SDKs/not-env-sdk-python/README.md)
