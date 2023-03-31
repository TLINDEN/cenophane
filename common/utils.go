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
package common

import (
	"errors"
	"regexp"
)

/*
   Untaint user input, that is: remove all non supported chars.

   wanted is a  regexp matching chars we shall  leave. Everything else
   will be removed. Eg:

   untainted := Untaint(input, `[^a-zA-Z0-9\-]`)

   Returns a  new string  and an  error if the  input string  has been
   modified.  It's the  callers  choice  to decide  what  to do  about
   it. You may  ignore the error and use the  untainted string or bail
   out.
*/
func Untaint(input string, wanted *regexp.Regexp) (string, error) {
	untainted := wanted.ReplaceAllString(input, "")

	if len(untainted) != len(input) {
		return untainted, errors.New("Invalid input string!")
	}

	return untainted, nil
}
