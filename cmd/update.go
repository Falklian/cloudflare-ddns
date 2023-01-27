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
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Falklian/cloudflare-ddns/utils"
	"github.com/cloudflare/cloudflare-go"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

type Config struct {
	APIKey   string   `mapstructure:"cf_api_key"`
	Email    string   `mapstructure:"cf_email_address"`
	APIToken string   `mapstructure:"cf_api_token"`
	Zones    []string `mapstructure:"zones"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: `Update "A" records for all configured zones`,
	Long: `Updates the "A" records for all domains/zones found in the config file.

NOTE: *ALL* "A" records will be updated. If your zone has multiple "A" records, you
may want to exclude it from updating`,
	Run: func(cmd *cobra.Command, args []string) {
		if ok := checkConfigExists(); !ok {
			fmt.Println(color.RedString("Please run `cloudflare-ddns configure` to configure the application"))
			os.Exit(1)
		}

		var config Config
		viper.Unmarshal(&config)

		var api *cloudflare.API
		var err error

		if config.APIToken != "" {
			api, err = cloudflare.NewWithAPIToken(config.APIToken)
		} else {
			api, err = cloudflare.New(config.APIKey, config.Email)
		}

		if err != nil {
			fmt.Println(color.RedString("Error creating Cloudflare API client: %s", err))
			os.Exit(1)
		}

		context := context.Background()

		zones, err := api.ListZones(context)
		if err != nil {
			fmt.Println(color.RedString("Error listing zones: %s", err))
			os.Exit(1)
		}

		for _, zone := range zones {
			if slices.Contains(config.Zones, zone.Name) {
				fmt.Println(color.GreenString("Updating DNS records for %s", zone.Name))
				records, _, err := api.ListDNSRecords(context, cloudflare.ResourceIdentifier(zone.ID), cloudflare.ListDNSRecordsParams{})
				if err != nil {
					fmt.Println(color.RedString("Error listing DNS records: %s", err))
					os.Exit(1)
				}

				currentIp := utils.GetIp()
				fmt.Println(color.GreenString("Current IP: %s", currentIp))

				for _, record := range records {
					if record.Type == "A" {
						fmt.Println(color.GreenString("Updating DNS record %s", record.Name))

						if record.Content == currentIp {
							fmt.Println(color.YellowString("DNS record %s is already up to date", record.Name))
							continue
						}
						updatedRecord := cloudflare.UpdateDNSRecordParams{
							ID:      record.ID,
							Content: currentIp,
						}

						err := api.UpdateDNSRecord(context, cloudflare.ResourceIdentifier(zone.ID), updatedRecord)
						if err != nil {
							fmt.Println(color.RedString("Error updating DNS record: %s", err))
							os.Exit(1)
						}
						fmt.Println(color.GreenString("DNS record %s updated successfully", record.Name))
					}
				}
			}
		}
	},
}

func checkConfigExists() bool {
	if _, err := os.Stat(viper.ConfigFileUsed()); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
