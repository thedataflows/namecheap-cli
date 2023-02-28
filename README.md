# Namecheap DNS command line interface

This tool was primarily built for backup and restore of Namecheap configuration for a particular domain.
Only one domain can be manipulated in a run, but you can create multiple configurations for each domain you want to manage.

## Namecheap setup

1. For testing, is highly recommended to use a sandbox account: <https://www.sandbox.namecheap.com/settings/tools/apiaccess/>
2. For production, enable API access: <https://ap.www.namecheap.com/settings/tools/apiaccess/>
3. Generate an API key
4. Add your IP to the whitelist

> **Note**: Namecheap API does not support update or append of one record. So whatever you pass as input for `set` command will overwrite the entire DNS configuration! To go around this all-or-nothing approach, use `setone` command to upsert/delete a single entry so `namecheap-cli` will download existing configuration, patch it, then upload it back in one go.

## Run It 🏃

`go run main.go --config sample/sandbox.yaml get`

## Usage

- `namecheap-cli`

    ```properties
    Namecheap DNS command line interface

    All flags values can be provided via env vars starting with NAMECHEAP_*
    To pass a subcommand (e.g. 'serve') flag, use NAMECHEAP_GET_FLAGNAME=somevalue

    Usage:
    namecheap-cli [flags]
    namecheap-cli [command]

    Available Commands:
    completion  Generate the autocompletion script for the specified shell
    convert     Convert Namecheap DNS configuration between local storage formats
    get         Download Namecheap DNS configuration
    help        Help about any command
    set         Upload Namecheap DNS configuration
    setone      create/update/delete a single DNS entry
    version     Display version and exit

    Flags:
        --config strings      Config file(s) or directories. When just dirs, file 'main' with extensions 'json, toml, yaml, yml, properties, props, prop, hcl, tfvars, dotenv, env, ini' is looked up. Can be specified multiple times (default [.,C:\Users\cri\AppData\Roaming\main])
    -h, --help                help for namecheap-cli
        --log-format string   Set log format to one of: 'console, json' (default "console")
        --log-level string    Set log level to one of: 'trace, debug, info, warn, error, fatal, panic, disabled' (default "info")

    Use "namecheap-cli [command] --help" for more information about a command.
    ```

- `namecheap-cli covert -h`

    ```properties
    Convert Namecheap DNS configuration between local storage formats

    Usage:
    namecheap-cli convert [flags]

    Aliases:
    convert, c

    Flags:
        --force                  Overwrite the file if exists
    -h, --help                   help for convert
    -i, --input-file string      Input file. If omitted, stdin is used until 2 consecutive newlines are detected
        --input-format string    Input format. Supported: [xml yaml json] (default "xml")
    -o, --output-file string     Output file. If omitted, outputs to stdout
        --output-format string   Output format. Supported: [xml yaml json] (default "yaml")
    -s, --sld string             Namecheap second-level domain, e.g.: 'example'.
    -t, --tld string             Namecheap top-level domain, e.g.: 'com'.

    Global Flags:
        --config strings      Config file(s) or directories. When just dirs, file 'main' with extensions 'json, toml, yaml, yml, properties, props, prop, hcl, tfvars, dotenv, env, ini' is looked up. Can be specified multiple times (default [.,C:\Users\cri\AppData\Roaming\main])
        --log-format string   Set log format to one of: 'console, json' (default "console")
        --log-level string    Set log level to one of: 'trace, debug, info, warn, error, fatal, panic, disabled' (default "info")
    ```

- `namecheap-cli get -h`

    ```properties
    Download Namecheap DNS configuration

    Usage:
    namecheap-cli get [flags]

    Aliases:
    get, g

    Flags:
        --client-ip string       Client IP. This is not really required (default "127.0.0.1")
        --force                  Force overwriting the file if exists
    -h, --help                   help for get
    -k, --key string             [Required] Namecheap API key
    -o, --output-file string     Output file. If omitted, outputs to stdout
        --output-format string   Output format. Supported: [xml yaml json] (default "xml")
        --sandbox                Use Namecheap sandbox API
    -s, --sld string             [Required] Namecheap second-level domain, e.g.: 'example'
        --timeout duration       Request timeout (default 10ns)
    -t, --tld string             [Required] Namecheap top-level domain, e.g.: 'com'
    -u, --username string        [Required] Namecheap user

    Global Flags:
        --config strings      Config file(s) or directories. When just dirs, file 'main' with extensions 'json, toml, yaml, yml, properties, props, prop, hcl, tfvars, dotenv, env, ini' is looked up. Can be specified multiple times (default [.,C:\Users\cri\AppData\Roaming\main])
        --log-format string   Set log format to one of: 'console, json' (default "console")
        --log-level string    Set log level to one of: 'trace, debug, info, warn, error, fatal, panic, disabled' (default "info")
    ```

- `namecheap-cli set -h`

    ```properties
    Upload Namecheap DNS configuration

    Usage:
    namecheap-cli set [flags]

    Aliases:
    set, s

    Flags:
        --client-ip string      Client IP. This is not really required (default "127.0.0.1")
    -h, --help                  help for set
    -i, --input-file string     Input file. If omitted, stdin is used until 2 consecutive newlines are detected
        --input-format string   Input format. Supported: [xml yaml json] (default "xml")
    -k, --key string            [Required] Namecheap API key
        --sandbox               Use Namecheap sandbox API
    -s, --sld string            Namecheap second-level domain, e.g.: 'example'. Can be read from the input file
        --timeout duration      Request timeout (default 10ns)
    -t, --tld string            Namecheap top-level domain, e.g.: 'com'. Can be read from the input file
    -u, --username string       [Required] Namecheap user

    Global Flags:
        --config strings      Config file(s) or directories. When just dirs, file 'main' with extensions 'json, toml, yaml, yml, properties, props, prop, hcl, tfvars, dotenv, env, ini' is looked up. Can be specified multiple times (default [.,C:\Users\cri\AppData\Roaming\main])
        --log-format string   Set log format to one of: 'console, json' (default "console")
        --log-level string    Set log level to one of: 'trace, debug, info, warn, error, fatal, panic, disabled' (default "info")
    ```

- `namecheap-cli setone -h`

    ```properties
    Create/update/delete a single DNS entry

    Usage:
    namecheap-cli setone [flags]

    Aliases:
    setone, o

    Flags:
        --address string        [Required] Record value
        --client-ip string      Client IP. This is not really required (default "127.0.0.1")
        --delete                Delete DNS entry
        --friendlyname string   Friendly name
    -h, --help                  help for setone
        --isactive              Active state (default true)
    -k, --key string            [Required] Namecheap API key
        --mxpref string         MXPref
        --name string           [Required] Record name
        --sandbox               Use Namecheap sandbox API
    -s, --sld string            [Required] Namecheap second-level domain, e.g.: 'example'
        --timeout duration      Request timeout (default 10ns)
    -t, --tld string            [Required] Namecheap top-level domain, e.g.: 'com'
        --ttl string            Time to live in seconds. 1799 is Namecheap's equivalent to 'Automatic' (default "1799")
        --type string           [Required] Record type
    -u, --username string       [Required] Namecheap user

    Global Flags:
        --config strings      Config file(s) or directories. When just dirs, file 'main' with extensions 'json, toml, yaml, yml, properties, props, prop, hcl, tfvars, dotenv, env, ini' is looked up. Can be specified multiple times (default [.,C:\Users\cri\AppData\Roaming\main])
        --log-format string   Set log format to one of: 'console, json' (default "console")
        --log-level string    Set log level to one of: 'trace, debug, info, warn, error, fatal, panic, disabled' (default "info")
    ```

## Configure It ☑️

- See [sample/sandbox.yaml](./sample/sandbox.yaml) for config file
- All parameters can be set via flags or env as well: `NAMECHEAP_<subcommand>_<flag>`, example: `NAMECHEAP_GET_KEY=1122334455`

## Test It 🧪

Test for coverage and race conditions

`make coverage`

## Lint It 👕

`make pre-commit run --all-files --show-diff-on-failure`

## Roadmap

- [ ] ?

## Development

### Build

- Preferably: `goreleaser build --clean --single-target` or
- `make build` or
- `scripts/local-build.sh` (deprecated)
