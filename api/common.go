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
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/tlinden/up/upd/cfg"
	"regexp"
	"strconv"
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

// stores 1 upload object, gets into db
type Upload struct {
	Id       string    `json:"id"`
	Expire   string    `json:"expire"`
	File     string    `json:"file"`    // final filename (visible to the downloader)
	Members  []string  `json:"members"` // contains multiple files, so File is an archive
	Uploaded Timestamp `json:"uploaded"`
	Context  string    `json:"context"`
	Url      string    `json:"url"`
}

// this one is also used for marshalling to the client
type Uploads struct {
	Entries []*Upload `json:"uploads"`

	// integrate the Result struct so we can signal success
	Result
}

// incoming id
type Id struct {
	Id string `json:"name" xml:"name" form:"name"`
}

// vaious helbers
func Log(format string, values ...any) {
	fmt.Printf("[DEBUG] "+format+"\n", values...)
}

func Ts() string {
	t := time.Now()
	return t.Format("2006-01-02-15-04-")
}

/*
   We could use time.ParseDuration(), but this doesn't support days.

   We  could also  use github.com/xhit/go-str2duration/v2,  which does
   the job,  but it's  just another dependency,  just for  this little
   gem. And  we don't need a  time.Time value.

   Convert a  duration into  seconds (int).
   Valid  time units  are "s", "m", "h" and "d".
*/
func duration2int(duration string) int {
	re := regexp.MustCompile(`(\d+)([dhms])`)
	seconds := 0

	for _, match := range re.FindAllStringSubmatch(duration, -1) {
		if len(match) == 3 {
			v, _ := strconv.Atoi(match[1])
			switch match[2][0] {
			case 'd':
				seconds += v * 86400
			case 'h':
				seconds += v * 3600
			case 'm':
				seconds += v * 60
			case 's':
				seconds += v
			}
		}
	}

	return seconds
}

/*
   Calculate   if  time   is   up  based   on   start  time.Time   and
   duration. Returns  true if time  is expired. Start time  comes from
   the database.

aka:
   if(now - start) >= duration { time is up}
*/
func IsExpired(conf *cfg.Config, start time.Time, duration string) bool {
	var expiretime int // seconds

	now := time.Now()

	if duration == "asap" {
		expiretime = conf.DefaultExpire
	} else {
		expiretime = duration2int(duration)
	}

	if now.Unix()-start.Unix() >= int64(expiretime) {
		return true
	}

	return false
}

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

/*
   Retrieve the  API Context  name from the  session, assuming  is has
   been  successfully  authenticated. However,  if  there  are no  api
   contexts     defined,    we'll     use     'default'    (set     in
   auth.validateAPIKey()).

   If there's no apicontext in the session, assume unauth user, return ""
*/
func GetApicontext(c *fiber.Ctx) (string, error) {
	sess, err := Sessionstore.Get(c)
	if err != nil {
		return "", fmt.Errorf("Unable to initialize session store from context: " + err.Error())
	}

	apicontext := sess.Get("apicontext")
	if apicontext != nil {
		return apicontext.(string), nil
	}

	return "", nil
}
