# not-env-cli

not-env-cli is the command-line interface for not-env, a self-hosted environment variable management system. It provides commands to manage environments and variables, and integrates with your shell for easy variable loading.

## Overview

The CLI allows you to:
- Manage environments (create, list, delete)
- Import variables from .env files
- Manage variables (set, get, list, delete)
- Load variables into your shell session
- View environment metadata and API keys

## Installation

### From Source

```bash
cd not-env-cli
go build -o not-env
sudo mv not-env /usr/local/bin/
```

### Verify Installation

```bash
not-env --version
```

## Quick Start

### 1. Login

First, you need to login to your not-env backend:

```bash
not-env login
```

You'll be prompted for:
- **Backend URL**: The URL of your not-env backend (e.g., `https://not-env.example.com` or `http://localhost:1212`)
- **API Key**: Your API key (APP_ADMIN, ENV_ADMIN, or ENV_READ_ONLY)

**Example:**
```
Backend URL: https://not-env.example.com
API Key: dGVzdF9hcHBfYWRtaW5fa2V5X2hlcmU...
```

**If this works correctly, you should see:**
```
Logged in successfully!
```

Your credentials are saved to `~/.not-env/config`.

### 2. Create an Environment (APP_ADMIN)

Using an APP_ADMIN key:

```bash
not-env env create --name development --description "Development environment"
```

**Expected output:**
```
Environment created successfully!
ID: 1
Name: development
ENV_ADMIN key: dGVzdF9lbnZfYWRtaW5fa2V5X2hlcmU...
ENV_READ_ONLY key: dGVzdF9lbnZfcmVhZG9ubHlfa2V5X2hlcmU...

Save these keys securely!
```

**If this works correctly, you should see:**
- Environment ID and name
- Two API keys: ENV_ADMIN and ENV_READ_ONLY
- Save these keys - you'll need them to manage variables

### 3. Import from .env File

Create an environment and populate it from an existing .env file:

```bash
not-env login
# Enter ENV_ADMIN key when prompted, or use APP_ADMIN to create environment first

not-env env import --name development --file .env
```

**Example .env file:**
```
DB_HOST=localhost
DB_PORT=5432
DB_PASSWORD=secret123
API_KEY=abc123
```

**Expected output:**
```
Environment 'development' created and populated with 4 variables!
ENV_ADMIN key: dGVzdF9lbnZfYWRtaW5fa2V5X2hlcmU...
Save this key securely!
```

**If this works correctly, you should see:**
- Confirmation that the environment was created
- Number of variables imported
- The ENV_ADMIN key for managing variables

### 4. Login with ENV_ADMIN or ENV_READ_ONLY

To manage variables, login with an ENV_ADMIN or ENV_READ_ONLY key:

```bash
not-env login
# Enter backend URL and ENV_ADMIN key
```

**Verify login:**
```bash
not-env env show
```

**Expected output:**
```
Environment: development (ID: 1)
Description: Development environment
Created: 2024-01-15T10:30:00Z
Updated: 2024-01-15T10:30:00Z
```

**If this works correctly, you should see:**
- The environment name and ID
- Description and timestamps
- This confirms you're logged in with a valid ENV_* key

### 5. List Variables

```bash
not-env var list
```

**Expected output:**
```
Variables:
  DB_HOST=localhost
  DB_PORT=5432
  DB_PASSWORD=secret123
  API_KEY=abc123
```

**If this works correctly, you should see:**
- All variables in the current environment
- Each variable shown as KEY=VALUE

### 6. Load Variables into Shell

Load all variables into your current shell session:

```bash
eval "$(not-env env set)"
```

**Verify variables are loaded:**
```bash
echo $DB_HOST
```

**Expected output:**
```
localhost
```

**If this works correctly, you should see:**
- The variable value printed
- You can now use `$DB_HOST`, `$DB_PORT`, etc. in your shell

**Clear variables:**
```bash
eval "$(not-env env clear)"
```

**Verify variables are cleared:**
```bash
echo $DB_HOST
```

**Expected output:**
```
(empty line)
```

**If this works correctly, you should see:**
- No output (variable is unset)

## Commands Reference

### Authentication

#### `not-env login`
Login to not-env backend. Prompts for URL and API key.

#### `not-env logout`
Logout and clear saved credentials.

### Environment Management

#### `not-env env create --name NAME [--description DESC]`
Create a new environment. Requires APP_ADMIN key.

**Example:**
```bash
not-env env create --name production --description "Production environment"
```

#### `not-env env list`
List all environments. Requires APP_ADMIN key.

**Example output:**
```
Environments:
  ID: 1, Name: development, Description: Development environment
  ID: 2, Name: production, Description: Production environment
```

#### `not-env env delete --id ENV_ID`
Delete an environment. Requires APP_ADMIN key.

**Example:**
```bash
not-env env delete --id 1
```

#### `not-env env import --name NAME [--description DESC] --file PATH [--overwrite]`
Import variables from a .env file. Creates environment if it doesn't exist.

**Example:**
```bash
not-env env import --name development --file .env --description "Dev env"
```

#### `not-env env show`
Show current environment metadata. Works with any ENV_* key.

**Example output:**
```
Environment: development (ID: 1)
Description: Development environment
Created: 2024-01-15T10:30:00Z
Updated: 2024-01-15T10:30:00Z
```

#### `not-env env update [--name NAME] [--description DESC]`
Update environment name and/or description. Requires ENV_ADMIN key.

**Example:**
```bash
not-env env update --name staging --description "Staging environment"
```

#### `not-env env keys`
Show ENV_ADMIN and ENV_READ_ONLY keys for current environment. Requires ENV_ADMIN key.

**Example output:**
```
Environment Keys:
ENV_ADMIN: dGVzdF9lbnZfYWRtaW5fa2V5X2hlcmU...
ENV_READ_ONLY: dGVzdF9lbnZfcmVhZG9ubHlfa2V5X2hlcmU...
```

#### `not-env env set`
Print `export KEY="value"` commands for all variables. Use with `eval` to load into shell.

**Example:**
```bash
eval "$(not-env env set)"
```

**Example output (when run directly):**
```
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_PASSWORD="secret123"
```

#### `not-env env clear`
Print `unset KEY` commands for all variables. Use with `eval` to clear from shell.

**Example:**
```bash
eval "$(not-env env clear)"
```

### Variable Management

#### `not-env var list`
List all variables in current environment.

**Example output:**
```
Variables:
  DB_HOST=localhost
  DB_PORT=5432
```

#### `not-env var get KEY`
Get a single variable value.

**Example:**
```bash
not-env var get DB_HOST
```

**Expected output:**
```
DB_HOST=localhost
```

#### `not-env var set KEY VALUE`
Set a variable value. Requires ENV_ADMIN key.

**Example:**
```bash
not-env var set DB_HOST localhost
```

**Expected output:**
```
Variable DB_HOST set successfully!
```

#### `not-env var delete KEY`
Delete a variable. Requires ENV_ADMIN key.

**Example:**
```bash
not-env var delete DB_HOST
```

**Expected output:**
```
Variable DB_HOST deleted successfully!
```

## Configuration

The CLI stores configuration in `~/.not-env/config` (TOML format):

```toml
url = "https://not-env.example.com"
api_key = "your-api-key-here"
```

## Error Handling

### Wrong Key Type

If you use the wrong key type for a command, you'll see:

```
Error: HTTP 403: insufficient permissions: requires [APP_ADMIN]
```

**Solution:** Use the correct key type (APP_ADMIN for environment management, ENV_ADMIN for variable management).

### Not Logged In

If you're not logged in:

```
Error: not logged in. Run 'not-env login' first
```

**Solution:** Run `not-env login` to authenticate.

### Backend Unreachable

If the backend is unreachable:

```
Error: failed to connect to backend: dial tcp: lookup not-env.example.com: no such host
```

**Solution:** Check the backend URL and ensure the backend is running.

## Integration with Backend and SDK

- **Backend**: The CLI communicates with [not-env-backend](../not-env-backend/README.md) via HTTPS
- **SDK**: Variables loaded via `not-env env set` can be used alongside the [JavaScript SDK](../not-env-sdk-js/README.md) in applications

## Examples

### Complete Workflow

```bash
# 1. Login with APP_ADMIN
not-env login
# Enter: https://not-env.example.com
# Enter: <APP_ADMIN_KEY>

# 2. Create environment
not-env env create --name myapp --description "My Application"

# 3. Login with ENV_ADMIN (from step 2 output)
not-env login
# Enter: https://not-env.example.com
# Enter: <ENV_ADMIN_KEY>

# 4. Set variables
not-env var set DB_HOST localhost
not-env var set DB_PORT 5432

# 5. List variables
not-env var list

# 6. Load into shell
eval "$(not-env env set)"
echo $DB_HOST  # Should print: localhost
```

## Troubleshooting

### Variables not loading in shell

- Ensure you're using `eval "$(not-env env set)"` with quotes
- Check that you're logged in with a valid ENV_* key
- Verify variables exist: `not-env var list`

### Permission denied errors

- Check your API key type matches the command requirements
- APP_ADMIN: environment management
- ENV_ADMIN: variable management
- ENV_READ_ONLY: read-only access

### Import fails

- Ensure the .env file exists and is readable
- Check file format (KEY=VALUE, one per line)
- Verify you have ENV_ADMIN permissions

