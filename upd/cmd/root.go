/*
Copyright © 2023 Thomas von Dein

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tlinden/up/upd/api"
	"github.com/tlinden/up/upd/cfg"
	"os"
	"strings"
)

var (
	cfgFile string
)

func completion(cmd *cobra.Command, mode string) error {
	switch mode {
	case "bash":
		return cmd.Root().GenBashCompletion(os.Stdout)
	case "zsh":
		return cmd.Root().GenZshCompletion(os.Stdout)
	case "fish":
		return cmd.Root().GenFishCompletion(os.Stdout, true)
	case "powershell":
		return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
	default:
		return errors.New("Invalid shell parameter! Valid ones: bash|zsh|fish|powershell")
	}
}

func Execute() {
	var (
		conf           cfg.Config
		ShowVersion    bool
		ShowCompletion string
	)

	var rootCmd = &cobra.Command{
		Use:   "upd [options]",
		Short: "upload daemon",
		Long:  `Run an upload daemon reachable via REST API`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if ShowVersion {
				fmt.Println(cfg.Getversion())
				return nil
			}

			if len(ShowCompletion) > 0 {
				return completion(cmd, ShowCompletion)
			}

			conf.ApplyDefaults()

			// actual execution starts here
			return api.Runserver(&conf, args)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig(cmd)
		},
	}

	// options
	rootCmd.PersistentFlags().BoolVarP(&ShowVersion, "version", "v", false, "Print program version")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "custom config file")
	rootCmd.PersistentFlags().BoolVarP(&conf.Debug, "debug", "d", false, "Enable debugging")
	rootCmd.PersistentFlags().StringVarP(&conf.Listen, "listen", "l", ":8080", "listen to custom ip:port (use [ip]:port for ipv6)")
	rootCmd.PersistentFlags().StringVarP(&conf.StorageDir, "storagedir", "s", "/tmp", "storage directory for uploaded files")
	rootCmd.PersistentFlags().StringVarP(&conf.ApiPrefix, "apiprefix", "a", "/api", "API endpoint path")
	rootCmd.PersistentFlags().StringVarP(&conf.Url, "url", "u", "", "HTTP endpoint w/o path")
	rootCmd.PersistentFlags().StringVarP(&conf.DbFile, "dbfile", "D", "/tmp/uploads.db", "Bold database file to use")

	// server settings
	rootCmd.PersistentFlags().BoolVarP(&conf.V4only, "ipv4", "4", false, "Only listen on ipv4")
	rootCmd.PersistentFlags().BoolVarP(&conf.V6only, "ipv6", "6", false, "Only listen on ipv6")
	rootCmd.MarkFlagsMutuallyExclusive("ipv4", "ipv6")

	rootCmd.PersistentFlags().BoolVarP(&conf.Prefork, "prefork", "p", false, "Prefork server threads")
	rootCmd.PersistentFlags().StringVarP(&conf.AppName, "appname", "n", "upd "+conf.GetVersion(), "App name to say hi as")
	rootCmd.PersistentFlags().IntVarP(&conf.BodyLimit, "bodylimit", "b", 10250000000, "Max allowed upload size in bytes")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// initialize viper, read config and ENV, bind flags
func initConfig(cmd *cobra.Command) error {
	v := viper.New()
	viper.SetConfigType("hcl")

	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName("upd")

		// default location[s]
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME/.config/upd")
		v.AddConfigPath("/etc")
		v.AddConfigPath("/usr/local/etc")

	}

	// ignore read errors, report all others
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
		fmt.Println(err)
	}

	fmt.Println("Using config file:", v.ConfigFileUsed())

	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	v.SetEnvPrefix("upd")

	// map flags to viper
	bindFlags(cmd, v)

	return nil
}

// bind flags to viper settings (env+cfgfile)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// map flag name to config variable
		configName := f.Name

		// use config variable if flag is not set and config is set
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
