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
package cfg

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const Version string = "v0.0.1"

var VERSION string // maintained by -x

type Apicontext struct {
	Context string `koanf:"context"` // aka name or tenant
	Key     string `koanf:"key"`
}

// holds the whole configs, filled by commandline flags, env and config file
type Config struct {
	ApiPrefix  string `koanf:"apiprefix"` // path prefix
	Debug      bool   `koanf:"debug"`
	Listen     string `koanf:"listen"`     // [host]:port
	StorageDir string `koanf:"storagedir"` // db and uploads go there
	Url        string `koanf:"url"`        // public visible url, might be different from Listen
	DbFile     string `koanf:"dbfile"`
	Super      string `koanf:"super"` // the apicontext which has all permissions

	// fiber settings, see:
	// https://docs.gofiber.io/api/fiber/#config
	Prefork   bool   `koanf:"prefork"`   // default: nope
	AppName   string `koanf:"appname"`   // "upd"
	BodyLimit int    `koanf:"bodylimit"` // much
	V4only    bool   `koanf:"ipv4"`
	V6only    bool   `koanf:"ipv6"`
	Network   string

	// only settable via config
	Apicontext []Apicontext `koanf:"apicontext"`

	// Internals only
	RegNormalizedFilename *regexp.Regexp
	RegDuration           *regexp.Regexp
	RegKey                *regexp.Regexp
	CleanInterval         time.Duration
	DefaultExpire         int
}

func Getversion() string {
	// main program version

	// generated  version string, used  by -v contains  cfg.Version on
	//  main  branch,   and  cfg.Version-$branch-$lastcommit-$date  on
	// development branch

	return fmt.Sprintf("This is upd version %s", VERSION)
}

func (c *Config) GetVersion() string {
	return VERSION
}

// post processing of options, if any
func (c *Config) ApplyDefaults() {
	if len(c.Url) == 0 {
		if strings.HasPrefix(c.Listen, ":") {
			c.Url = "http://localhost" + c.Listen
		} else {
			c.Url = "http://" + c.Listen
		}
	}

	switch {
	case c.V4only && c.V6only:
		c.Network = "tcp" // dual stack
	case c.V4only:
		c.Network = "tcp4"
	case c.V6only:
		c.Network = "tcp6"
	default:
		if c.Prefork {
			c.Network = "tcp4"
		} else {
			c.Network = "tcp" // dual stack
		}
	}

	c.RegNormalizedFilename = regexp.MustCompile(`[^\w\d\-_\.]`)
	c.RegDuration = regexp.MustCompile(`[^dhms0-9]`)
	c.RegKey = regexp.MustCompile(`[^a-zA-Z0-9\-]`)

	c.CleanInterval = 10 * time.Second
	c.DefaultExpire = 30 * 86400 // 1 month
}
