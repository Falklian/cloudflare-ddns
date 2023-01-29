# Cloudflare DDNS

A tool for automatically updating dynamic IP addresses in Cloudflare

## Purpose

Having an internet connection with a dynamic IP address can cause hiccups for self-hosters. This tool can be set up using a cronjob to periodically update Cloudflare with any possible IP address changes

## Usage

### Commands

- `completion`: Generates shell completion scripts
  - Sub-commands
    - `bash`: Generates bash completions
    - `fish`: Generates fish completions
    - `powershell`: Generates powershell completions
    - `zsh`: Generates zsh completions
    - Flags
      - `--no-descriptions`: Disables descriptions in completions
- `configure`: Run configuration wizard (required)
  - Aliases
    - `config`
    - `conf`
- `help`: Help about any command
- `publicIp`: Gets your public IP address
  - Aliases
    - `ip`
- `update`: Updates DNS "A" records with your current IP address
  - Flags
    - `--api-key`: Cloudflare Global API key. Must be used with `--email`
    - `--email`: Cloudflare account email address. Must be used with `--api-key`
    - `--token`: Cloudflare API token with `Zone:DNS:Edit` permissions. Can be used in place of `--api-key` and `--email`
    - `--zone`: comma-seperated list of Cloudflare zone IDs
- Global flags
  - `--config`: Path to config file (default: `~/.cloudflare-ddns/config.yaml`)
  - `--cron`: Disables colored output (default: `false`)

Binaries exist for Windows, Linux, and macOS and can be found in the Releases section. Download the appropriate build for your operating system and save it to a directory that makes sense for your setup (e.g. `/usr/local/bin`). You can also run `go install github.com/Falklian/cloudflare-ddns@latest` to install the latest version of the tool to your `$GOPATH/bin` directory

Before using the utility, you will want to create a config file with the command `cloudflare-ddns configure`. Both global Cloudflare API key/email addresss, or a scoped API token with `Zone:DNS:Edit` permissions will work for configuration

As an alternative to creating a config file, you can pass the `--api-key`, `--email`, and `--zone` flags to the `update` command. If you are using a scoped API token, you can pass the `--token` flag instead of `--api-key` and `--email`

To update DNS automatically, use the `update` command. This would best be set up with a cronjob:

```shell
# Open crontab
$ crontab -e

# Add a schedule that makes sense for you (example: Run every 5 mintues):
*/5 * * * * /path/to/cloudflare-ddns update --cron >> /path/to/logfile.log 2>&1
```

There is also a `publicIp` command that will echo your public IP address to the terminal. This command does not require configuration to run

## Building and Contributing

### Prerequisites

- Go 1.19+

### Building

```shell
# Clone the repository
$ git clone https://github.com/Falklian/cloudflare-ddns.git

# Change into the directory
$ cd cloudflare-ddns

# Build the binary
$ go build -o build/cloudflare-ddns -ldflags '-s -w'
```

## Feature Wishlist

- [x] ~~Add flags for use in place of generating a config file~~
- [ ] Add logging support
- [ ] Add support to allow specifiying "A" records to update
- [ ] Add ipv6 support
- [ ] Create Docker image
