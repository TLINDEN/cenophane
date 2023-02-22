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

package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/tlinden/up/upctl/cfg"
	//"path/filepath"
	"time"
)

type Response struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func Runclient(c *cfg.Config, args []string) error {
	if len(args) == 0 {
		return errors.New("No files specified to upload.")
	}

	client := req.C()
	if c.Debug {
		client.DevMode()
	}

	client.SetUserAgent("upctl-" + cfg.VERSION)

	url := c.Endpoint + "/putfile"

	rq := client.R()
	for _, file := range args {
		rq.SetFile("upload[]", file)
	}

	if c.Retries > 0 {
		// Enable retry and set the maximum retry count.
		rq.SetRetryCount(c.Retries).
			//  Set  the  retry  sleep   interval  with  a  commonly  used
			//   algorithm:  capped   exponential   backoff  with   jitter
			// (https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/).
			SetRetryBackoffInterval(1*time.Second, 5*time.Second).
			AddRetryHook(func(resp *req.Response, err error) {
				req := resp.Request.RawRequest
				if c.Debug {
					fmt.Println("Retrying endpoint request:", req.Method, req.URL)
				}
			})
	}

	resp, err := rq.
		SetFormData(map[string]string{
			"expire": "1d",
		}).
		SetUploadCallbackWithInterval(func(info req.UploadInfo) {
			fmt.Printf("\r%q uploaded %.2f%%", info.FileName, float64(info.UploadedSize)/float64(info.FileSize)*100.0)
		}, 10*time.Millisecond).
		Post(url)

	fmt.Println("")

	if err != nil {
		return err
	}

	r := Response{}

	json.Unmarshal([]byte(resp.String()), &r)

	fmt.Println(r)

	if c.Debug {
		trace := resp.TraceInfo()  // Use `resp.Request.TraceInfo()` to avoid unnecessary struct copy in production.
		fmt.Println(trace.Blame()) // Print out exactly where the http request is slowing down.
		fmt.Println("----------")
		fmt.Println(trace)
	}

	return nil
}
