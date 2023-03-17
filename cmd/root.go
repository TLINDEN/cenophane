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

	"github.com/knadh/koanf/parsers/hcl"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"

	flag "github.com/spf13/pflag"

	"github.com/alecthomas/repr"
	"github.com/tlinden/cenophane/api"
	"github.com/tlinden/cenophane/cfg"

	"os"
	"path/filepath"
	"strings"
)

var (
	cfgFile string
)

func Execute() error {
	var (
		conf        cfg.Config
		ShowVersion bool
	)

	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}

	f.BoolVarP(&ShowVersion, "version", "v", false, "Print program version")
	f.StringVarP(&cfgFile, "config", "c", "", "custom config file")
	f.BoolVarP(&conf.Debug, "debug", "d", false, "Enable debugging")
	f.StringVarP(&conf.Listen, "listen", "l", ":8080", "listen to custom ip:port (use [ip]:port for ipv6)")
	f.StringVarP(&conf.StorageDir, "storagedir", "s", "/tmp", "storage directory for uploaded files")
	f.StringVarP(&conf.ApiPrefix, "apiprefix", "a", "/api", "API endpoint path")
	f.StringVarP(&conf.Url, "url", "u", "", "HTTP endpoint w/o path")
	f.StringVarP(&conf.DbFile, "dbfile", "D", "/tmp/uploads.db", "Bold database file to use")
	f.StringVarP(&conf.Super, "super", "", "", "The API Context which has permissions on all contexts")

	// server settings
	f.BoolVarP(&conf.V4only, "ipv4", "4", false, "Only listen on ipv4")
	f.BoolVarP(&conf.V6only, "ipv6", "6", false, "Only listen on ipv6")

	f.BoolVarP(&conf.Prefork, "prefork", "p", false, "Prefork server threads")
	f.StringVarP(&conf.AppName, "appname", "n", "cenod "+conf.GetVersion(), "App name to say hi as")
	f.IntVarP(&conf.BodyLimit, "bodylimit", "b", 10250000000, "Max allowed upload size in bytes")
	f.StringSliceP("apikeys", "", []string{}, "Api key[s] to allow access")

	f.Parse(os.Args[1:])

	// exclude -6 and -4
	if conf.V4only && conf.V6only {
		return errors.New("You cannot mix -4 and -6!")
	}

	// config provider
	var k = koanf.New(".")

	// Load the config files provided in the commandline or the default locations
	configfiles := []string{}
	configfile, _ := f.GetString("config")
	if configfile != "" {
		configfiles = []string{configfile}
	} else {
		configfiles = []string{
			"/etc/cenod.hcl", "/usr/local/etc/cenod.hcl", // unix variants
			filepath.Join(os.Getenv("HOME"), ".config", "cenod", "cenod.hcl"),
			filepath.Join(os.Getenv("HOME"), ".cenod"),
			"cenod.hcl",
		}
	}

	for _, cfgfile := range configfiles {
		if _, err := os.Stat(cfgfile); err == nil {
			if err := k.Load(file.Provider(cfgfile), hcl.Parser(true)); err != nil {
				return errors.New("error loading config file: " + err.Error())
			}
		}
		// else: we ignore the file if it doesn't exists
	}

	// env overrides config file
	k.Load(env.Provider("CENOD_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "CENOD_")), "_", ".", -1)
	}), nil)

	// command line overrides env
	if err := k.Load(posflag.Provider(f, ".", k), nil); err != nil {
		return errors.New("error loading config: " + err.Error())
	}

	// fetch values
	k.Unmarshal("", &conf)

	if conf.Debug {
		repr.Print(conf)
	}

	switch {
	case ShowVersion:
		fmt.Println(cfg.Getversion())
		return nil
	default:
		conf.ApplyDefaults()
		return api.Runserver(&conf, flag.Args())
	}
}
