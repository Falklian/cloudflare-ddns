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
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cron    bool
	version string = "(development)"
)

var rootCmd = &cobra.Command{
	Use:     "cloudflare-ddns",
	Version: version,
}

func Execute() error {
	err := rootCmd.Execute()
	return err
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cloudflare-ddns/config.yml)")
	rootCmd.PersistentFlags().BoolVar(&cron, "cron", false, "run in cron mode (color output disabled; default is false)")
	rootCmd.AddCommand(configureCmd, publicIpCmd, updateCmd)
	rootCmd.SetVersionTemplate("cloudflare-ddns {{.Version}}\n")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.SetConfigName("config.yml")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(fmt.Sprintf("%s/.cloudflare-ddns", home))
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv()

	if cron {
		color.NoColor = true
	}

	if err := viper.ReadInConfig(); err != nil {
		// The below should work according to the docs, but it doesn't
		// if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		if ok := errors.Is(err, os.ErrNotExist); ok {
			fmt.Println(color.YellowString("Config file `%s` not found", filepath.Base(viper.ConfigFileUsed())))
		}
	}
}
