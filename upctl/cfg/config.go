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
	//"strings"
)

const Version string = "v0.0.1"

var VERSION string // maintained by -x

type Config struct {
	// globals
	Endpoint string
	Debug    bool
	Retries  int
	Silent   bool

	// used for authentication
	Apikey string

	// upload
	Expire string

	// used for filtering (list command)
	Apicontext string

	// required to intercept requests using httpmock in tests
	Mock bool

	// used to filter lists
	Query string

	// required for forms
	Description string
	Notify      string
}

func Getversion() string {
	// main program version

	// generated  version string, used  by -v contains  cfg.Version on
	//  main  branch,   and  cfg.Version-$branch-$lastcommit-$date  on
	// development branch

	return fmt.Sprintf("This is upctl version %s", VERSION)
}
