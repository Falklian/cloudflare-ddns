/*
Copyright © 2022 Bill Walker <bill@billw.dev>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configureCmd = &cobra.Command{
	Use:     "configure",
	Aliases: []string{"config", "conf"},
	Short:   "Create configuration file (required)",
	Long: `Creates a configuration file for cloudflare-ddns. This is required before running any other commands.

Scoped Cloudflare API tokens can be used, as well as the Global API Key and account email address.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(color.GreenString("\n---=== Configuring Cloudflare DDNS client ===---"))
		if ok := createConfig(); !ok {
			fmt.Println(color.RedString("\nFailed to create configuration file."))
			os.Exit(1)
		}
	},
	DisableFlagsInUseLine: true,
}

func createConfig() bool {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print(color.GreenString(`
Please enter either your Cloudflare API key and email address, or API token.
You can find API keys/tokens at https://dash.cloudflare.com/profile/api-tokens.

CTRL+C to quit at any time

`))

	if _, err := os.Stat(viper.ConfigFileUsed()); !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Configuration file already exists at %s. Overwrite? [y/N] ", viper.ConfigFileUsed())
		scanner.Scan()
		if strings.ToLower(scanner.Text()) != "y" {
			return false
		}
		fmt.Print("\n")
	}

	if viper.ConfigFileUsed() == "" {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.SetConfigFile(fmt.Sprintf("%s/.cloudflare-ddns/config.yml", home))
	}

	if _, err := os.Stat(filepath.Dir(viper.ConfigFileUsed())); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(filepath.Dir(viper.ConfigFileUsed()), 0755); err != nil {
			fmt.Println(color.RedString("Failed to create configuration directory: %s", err))
			return false
		}
	}

	fmt.Print("Cloudflare API token (leave blank to use API key and email address): ")
	scanner.Scan()
	cfAPIToken := scanner.Text()

	if cfAPIToken != "" {
		viper.Set("api-token", cfAPIToken)
	} else {
		fmt.Print("Cloudflare API key: ")
		scanner.Scan()
		viper.Set("api-key", scanner.Text())

		fmt.Print("Cloudflare email address: ")
		scanner.Scan()
		viper.Set("email", scanner.Text())
	}

	fmt.Print("List of domains to update (comma-separated): ")
	scanner.Scan()
	zones := strings.Split(scanner.Text(), ",")
	for i, zone := range zones {
		zones[i] = strings.TrimSpace(zone)
	}
	viper.Set("zones", zones)

	viper.SetConfigType("yaml")

	if err := viper.WriteConfigAs(viper.ConfigFileUsed()); err != nil {
		fmt.Println(color.RedString("Failed to write configuration file: %s", err))
		return false
	}

	return true
}
