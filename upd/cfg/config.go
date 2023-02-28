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
	"strings"
)

const Version string = "v0.0.1"

var VERSION string // maintained by -x

type Config struct {
	ApiPrefix  string
	Debug      bool
	Listen     string
	StorageDir string
	Url        string
	DbFile     string

	// fiber settings, see:
	// https://docs.gofiber.io/api/fiber/#config
	Prefork   bool
	AppName   string
	BodyLimit int
	V4only    bool
	V6only    bool
	Network   string
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
}
