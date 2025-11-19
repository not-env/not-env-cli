# not-env-cli Requirements

This document specifies the functional and non-functional requirements for not-env-cli.

## Functional Requirements

### FR1: Configuration Management

**FR1.1:** The CLI must store configuration in `~/.not-env/config` (TOML format).

**FR1.2:** Configuration must include:
- `url`: Backend URL (string)
- `api_key`: API key (string)
- `env_id`: Current environment ID (optional integer)

**FR1.3:** The CLI must provide `login` command that:
- Prompts for backend URL and API key
- Validates credentials by making a test request to `/health`
- Saves configuration to `~/.not-env/config`
- Creates config directory if it doesn't exist
- Sets file permissions to 0600

**FR1.4:** The CLI must provide `logout` command that:
- Removes the configuration file
- Handles missing file gracefully

**FR1.5:** All commands (except `login` and `logout`) must load configuration and fail with a clear error if not logged in.

### FR2: Authentication

**FR2.1:** The CLI must include the API key in the `Authorization: Bearer <API_KEY>` header for all API requests.

**FR2.2:** The CLI must handle authentication errors (401, 403) and display clear error messages.

**FR2.3:** The CLI must validate API key type requirements and show appropriate errors for insufficient permissions.

### FR3: Environment Commands

**FR3.1:** `not-env env create --name NAME [--description DESC]`
- Requires APP_ADMIN key
- Creates a new environment
- Prints environment ID and generated ENV_ADMIN/ENV_READ_ONLY keys
- Fails if environment name already exists

**FR3.2:** `not-env env list`
- Requires APP_ADMIN key
- Lists all environments in the organization
- Displays ID, name, and description for each environment

**FR3.3:** `not-env env delete --id ENV_ID`
- Requires APP_ADMIN key
- Deletes environment and all its variables/keys
- Confirms deletion

**FR3.4:** `not-env env import --name NAME [--description DESC] --file PATH [--overwrite]`
- Reads .env file from PATH
- Parses KEY=VALUE pairs (supports quoted values)
- Creates environment if it doesn't exist (or uses existing if --overwrite)
- Populates all variables using ENV_ADMIN key
- Prints number of variables imported and ENV_ADMIN key

**FR3.5:** `not-env env show`
- Works with any ENV_* key type
- Displays current environment metadata (ID, name, description, timestamps)

**FR3.6:** `not-env env update [--name NAME] [--description DESC]`
- Requires ENV_ADMIN key
- Updates environment name and/or description
- Requires at least one flag

**FR3.7:** `not-env env keys`
- Requires ENV_ADMIN key
- Displays ENV_ADMIN and ENV_READ_ONLY keys for current environment

**FR3.8:** `not-env env set`
- Works with any ENV_* key type
- Fetches all variables from current environment
- Prints `export KEY="value"` lines for each variable
- Escapes special characters in values for shell safety
- Designed for use with `eval "$(not-env env set)"`

**FR3.9:** `not-env env clear`
- Works with any ENV_* key type
- Fetches all variable keys from current environment
- Prints `unset KEY` lines for each variable
- Designed for use with `eval "$(not-env env clear)"`

### FR4: Variable Commands

**FR4.1:** `not-env var list`
- Works with any ENV_* key type
- Lists all variables in current environment
- Displays KEY=VALUE pairs

**FR4.2:** `not-env var get KEY`
- Works with any ENV_* key type
- Gets a single variable by key
- Prints KEY=VALUE
- Returns error if variable not found

**FR4.3:** `not-env var set KEY VALUE`
- Requires ENV_ADMIN key
- Creates or updates a variable
- Confirms success

**FR4.4:** `not-env var delete KEY`
- Requires ENV_ADMIN key
- Deletes a variable
- Confirms success

### FR5: API Client

**FR5.1:** The CLI must provide an HTTP client wrapper that:
- Uses HTTPS for all requests
- Adds Authorization header automatically
- Handles JSON request/response bodies
- Provides methods: Get, Post, Put, Patch, Delete

**FR5.2:** The client must parse error responses and return meaningful error messages.

**FR5.3:** The client must handle network errors and display clear messages.

### FR6: Command Structure

**FR6.1:** The CLI must use cobra for command parsing.

**FR6.2:** Commands must be organized as:
- `not-env login`
- `not-env logout`
- `not-env env <subcommand>`
- `not-env var <subcommand>`

**FR6.3:** All commands must validate required flags and arguments.

**FR6.4:** Commands must display usage help with `--help` flag.

### FR7: Error Messages

**FR7.1:** The CLI must display clear, actionable error messages for:
- Missing configuration (not logged in)
- Invalid credentials
- Insufficient permissions
- Network errors
- Invalid arguments/flags
- Backend errors

**FR7.2:** Error messages must include:
- What went wrong
- Why it might have happened
- How to fix it (when applicable)

### FR8: Shell Integration

**FR8.1:** `not-env env set` output must be compatible with shell `eval`:
- Uses `export KEY="value"` format
- Escapes special characters (quotes, dollars, backticks)
- One variable per line

**FR8.2:** `not-env env clear` output must be compatible with shell `eval`:
- Uses `unset KEY` format
- One variable per line

**FR8.3:** The CLI must not print any non-export/unset output when used with `eval`.

## Non-Functional Requirements

### NFR1: Usability

**NFR1.1:** Commands must be intuitive and follow common CLI patterns.

**NFR1.2:** Help text must be clear and include examples.

**NFR1.3:** Error messages must be user-friendly and actionable.

### NFR2: Security

**NFR2.1:** Configuration file must have permissions 0600 (read/write for owner only).

**NFR2.2:** API keys must never be logged or printed except when explicitly requested (env keys command).

**NFR2.3:** The CLI must use HTTPS for all backend communication.

### NFR3: Performance

**NFR3.1:** Commands must complete in under 2 seconds for typical operations.

**NFR3.2:** The CLI must handle network timeouts gracefully.

### NFR4: Compatibility

**NFR4.1:** The CLI must work on:
- Linux
- macOS
- Windows (with WSL or Git Bash)

**NFR4.2:** Shell integration must work with:
- bash
- zsh
- fish (with minor modifications)

### NFR5: Dependencies

**NFR5.1:** The CLI must use:
- Go 1.21+
- `github.com/spf13/cobra` for command parsing
- `github.com/pelletier/go-toml/v2` for config file parsing
- Standard library `net/http` for HTTP client

## Implementation Constraints

### IC1: Config File Format

**IC1.1:** Config file must be TOML format.

**IC1.2:** Config file location: `~/.not-env/config`.

**IC1.3:** Config directory must be created with permissions 0700.

### IC2: API Communication

**IC2.1:** All requests must use HTTPS (or HTTP for localhost).

**IC2.2:** All requests must include `Authorization: Bearer <API_KEY>` header.

**IC2.3:** Request/response bodies must be JSON.

### IC3: Command Output

**IC3.1:** Success messages must be printed to stdout.

**IC3.2:** Error messages must be printed to stderr.

**IC3.3:** `env set` and `env clear` output must be to stdout only (no errors mixed in).

### IC4: .env File Parsing

**IC4.1:** Must support:
- `KEY=VALUE` format
- Quoted values: `KEY="value"` or `KEY='value'`
- Comments: lines starting with `#`
- Empty lines (ignored)

**IC4.2:** Must handle:
- Special characters in values
- Multi-line values (not supported in v1)
- Variable expansion (not supported in v1)

## Error Handling Specifications

### EH1: Not Logged In

**Message:** `not logged in. Run 'not-env login' first`

**When:** Configuration file doesn't exist or is invalid.

### EH2: Invalid Credentials

**Message:** `invalid credentials or backend unreachable`

**When:** Health check fails or returns non-200 status.

### EH3: Insufficient Permissions

**Message:** `HTTP 403: insufficient permissions: requires [KEY_TYPE]`

**When:** API returns 403 Forbidden.

### EH4: Environment Not Found

**Message:** `environment not found`

**When:** Environment ID doesn't exist or doesn't belong to organization.

### EH5: Variable Not Found

**Message:** `HTTP 404: variable not found`

**When:** Variable key doesn't exist in environment.

### EH6: Network Error

**Message:** `failed to connect to backend: <error details>`

**When:** Network request fails (timeout, DNS error, etc.).

## Expected Behaviors

### EB1: Login Flow

1. User runs `not-env login`
2. CLI prompts for URL (defaults to https:// if no protocol)
3. CLI prompts for API key
4. CLI validates by calling `/health` endpoint
5. CLI saves config to `~/.not-env/config`
6. CLI prints "Logged in successfully!"

### EB2: Environment Creation

1. User runs `not-env env create --name dev`
2. CLI loads config
3. CLI sends POST `/environments` with APP_ADMIN key
4. Backend creates environment and generates keys
5. CLI prints environment ID and both keys
6. CLI warns user to save keys

### EB3: Shell Integration

1. User runs `eval "$(not-env env set)"`
2. CLI loads config
3. CLI sends GET `/variables` with ENV_* key
4. CLI prints `export KEY="value"` for each variable
5. Shell evaluates exports and sets variables
6. User can now use `$KEY` in shell

### EB4: Variable Management

1. User runs `not-env var set DB_HOST localhost`
2. CLI loads config
3. CLI sends PUT `/variables/DB_HOST` with ENV_ADMIN key
4. Backend encrypts and stores value
5. CLI prints success message

