# Cloudflare DDNS

A tool for automatically updating dynamic IP addresses in Cloudflare

## Purpose

Having an internet connection with a dynamic address can cause hiccups for self-hosters. This tool can be set up using a cronjob to periodically update Cloudflare with any possible IP adress changes

## Usage

Binaries exist for Windows, Linux, and macOS (arm64 and amd64 builds) and can be found in the Releases section. Download the appropriate build for your operating system and save it to a directory that makes sense for your setup (e.g. `/usr/local/bin`)

Before using the utility, you will want to create a config file with the command `cloudflare-ddns configure`. Both global Cloudflare API key/email addresss, or a scoped API token with `Zone:DNS:Edit` permissions will work for configuration

To update DNS automatically, use the `update` command, optionally with the `--cron` flag to disable colored output. This would best be set up with a cronjob

```shell
# Open crontab
$ crontab -e

# Add a schedule that makes sense for you (example: Run every 5 mintues):
*/5 * * * * /path/to/cloudflare-ddns update --cron >> /path/to/logfile.log 2>&1
```

There is also a `publicIp` command that will echo your public IP address to the terminal. This command does not require configuration to run

## Known issues

- Doesn't always pick up the config file, which seems to be a possible problem with [Viper](https://github.com/spf13/viper)
