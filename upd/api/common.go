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

package api

import (
	"fmt"
	"regexp"
	"time"
)

const ApiVersion string = "/v1"

// used to return to the api client
type Result struct {
	Success bool
	Message string
	Code    int
}

// Binding from JSON, data coming from user, not tainted
type Meta struct {
	Expire string `json:"expire" form:"expire"`
}

// vaious helbers
func Log(format string, values ...any) {
	fmt.Printf("[DEBUG] "+format+"\n", values...)
}

func Ts() string {
	t := time.Now()
	return t.Format("2006-01-02-15-04-")
}

func NormalizeFilename(file string) string {
	r := regexp.MustCompile(`[^\w\d\-_\\.]`)

	return Ts() + r.ReplaceAllString(file, "")
}
