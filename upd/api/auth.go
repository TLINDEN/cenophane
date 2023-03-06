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
	"crypto/sha256"
	"crypto/subtle"
	"github.com/gofiber/fiber/v2"
	//"github.com/gofiber/keyauth/v2"
	"regexp"
	"strings"
)

var Authurls []*regexp.Regexp
var Apikeys []string

func AuthSetApikeys(keys []string) {
	Apikeys = keys
}

func AuthSetEndpoints(prefix string, version string, endpoints []string) {
	for _, endpoint := range endpoints {
		Authurls = append(Authurls, regexp.MustCompile("^"+prefix+version+endpoint))
	}
}

func validateAPIKey(c *fiber.Ctx, key string) (bool, error) {
	for _, apiKey := range Apikeys {
		hashedAPIKey := sha256.Sum256([]byte(apiKey))
		hashedKey := sha256.Sum256([]byte(key))

		if subtle.ConstantTimeCompare(hashedAPIKey[:], hashedKey[:]) == 1 {
			return true, nil
		}
	}
	return true, nil
	//return false, keyauth.ErrMissingOrMalformedAPIKey
}

func authFilter(c *fiber.Ctx) bool {
	originalURL := strings.ToLower(c.OriginalURL())

	for _, pattern := range Authurls {
		if pattern.MatchString(originalURL) {
			return false
		}
	}
	return true
}
