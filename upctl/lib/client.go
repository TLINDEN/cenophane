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
	//"github.com/alecthomas/repr"
	"github.com/imroc/req/v3"
	"github.com/schollz/progressbar/v3"
	"github.com/tlinden/up/upctl/cfg"
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type Response struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Request struct {
	R   *req.Request
	Url string
}

type ListParams struct {
	Apicontext string `json:"apicontext"`
}

type Upload struct {
	Id       string    `json:"id"`
	Expire   string    `json:"expire"`
	File     string    `json:"file"`    // final filename (visible to the downloader)
	Members  []string  `json:"members"` // contains multiple files, so File is an archive
	Uploaded Timestamp `json:"uploaded"`
	Context  string    `json:"context"`
	Url      string    `json:"url"`
}

type Uploads struct {
	Entries []*Upload `json:"uploads"`
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Code    int       `json:"code"`
}

const Maxwidth = 10

func Setup(c *cfg.Config, path string) *Request {
	client := req.C()
	if c.Debug {
		client.DevMode()
	}

	client.SetUserAgent("upctl-" + cfg.VERSION)

	R := client.R()

	if c.Retries > 0 {
		// Enable retry and set the maximum retry count.
		R.SetRetryCount(c.Retries).
			//  Set  the  retry  sleep   interval  with  a  commonly  used
			//   algorithm:  capped   exponential   backoff  with   jitter
			// (https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/).
			SetRetryBackoffInterval(1*time.Second, 5*time.Second).
			AddRetryHook(func(resp *req.Response, err error) {
				req := resp.Request.RawRequest
				if c.Debug {
					fmt.Println("Retrying endpoint request:", req.Method, req.URL, err)
				}
			})
	}

	if len(c.Apikey) > 0 {
		client.SetCommonBearerAuthToken(c.Apikey)
	}

	return &Request{Url: c.Endpoint + path, R: R}

}

func GatherFiles(rq *Request, args []string) error {
	for _, file := range args {
		info, err := os.Stat(file)

		if os.IsNotExist(err) {
			return err
		}

		if info.IsDir() {
			err := filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() {
					rq.R.SetFile("upload[]", path)
				}
				return nil
			})

			if err != nil {
				return err
			}
		} else {
			rq.R.SetFile("upload[]", file)
		}
	}

	return nil
}

func UploadFiles(c *cfg.Config, args []string) error {
	// setup url, req.Request, timeout handling etc
	rq := Setup(c, "/file/")

	// collect files to upload from @argv
	if err := GatherFiles(rq, args); err != nil {
		return err
	}

	// progres bar
	bar := progressbar.Default(100)
	var left float64

	// actual post w/ settings
	resp, err := rq.R.
		SetFormData(map[string]string{
			"expire": c.Expire,
		}).
		SetUploadCallbackWithInterval(func(info req.UploadInfo) {
			left = float64(info.UploadedSize) / float64(info.FileSize) * 100.0
			bar.Add(int(left))
		}, 10*time.Millisecond).
		Post(rq.Url)

	if err != nil {
		return err
	}

	return RespondExtended(resp)
}

func HandleResponse(c *cfg.Config, resp *req.Response) error {
	// we expect a json response, extract the error, if any
	r := Response{}

	if err := json.Unmarshal([]byte(resp.String()), &r); err != nil {
		// text output!
		r.Message = resp.String()
	}

	if c.Debug {
		trace := resp.Request.TraceInfo()
		fmt.Println(trace.Blame())
		fmt.Println("----------")
		fmt.Println(trace)
	}

	if !r.Success {
		if len(r.Message) == 0 {
			if resp.Err != nil {
				return resp.Err
			} else {
				return errors.New("Unknown error")
			}
		} else {
			return errors.New(r.Message)
		}
	}

	// all right
	if r.Message != "" {
		fmt.Println(r.Message)
	}

	return nil
}

func List(c *cfg.Config, args []string) error {
	rq := Setup(c, "/list/")

	params := &ListParams{Apicontext: c.Apicontext}
	resp, err := rq.R.
		SetBodyJsonMarshal(params).
		Get(rq.Url)

	if err != nil {
		return err
	}

	return RespondTable(resp)
}

func Delete(c *cfg.Config, args []string) error {
	for _, id := range args {
		rq := Setup(c, "/file/"+id+"/")

		resp, err := rq.R.Delete(rq.Url)

		if err != nil {
			return err
		}

		if err := HandleResponse(c, resp); err != nil {
			return err
		}

		fmt.Printf("Upload %s successfully deleted.\n", id)
	}

	return nil
}

func Describe(c *cfg.Config, args []string) error {
	id := args[0] // we describe only 1 object

	rq := Setup(c, "/upload/"+id+"/")
	resp, err := rq.R.Get(rq.Url)

	if err != nil {
		return err
	}

	return RespondExtended(resp)
}

func Download(c *cfg.Config, args []string) error {
	id := args[0]

	// progres bar
	bar := progressbar.Default(100)

	callback := func(info req.DownloadInfo) {
		if info.Response.Response != nil {
			bar.Add(1)
		}
	}

	rq := Setup(c, "/file/"+id+"/")
	resp, err := rq.R.
		SetOutputFile(id).
		SetDownloadCallback(callback).
		Get(rq.Url)

	if err != nil {
		return err
	}

	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil {
		os.Remove(id)
		return err
	}

	filename := params["filename"]
	if filename == "" {
		os.Remove(id)
		return fmt.Errorf("No filename provided!")
	}

	cleanfilename, _ := Untaint(filename, regexp.MustCompile(`[^a-zA-Z0-9\-\._]`))

	if err := os.Rename(id, cleanfilename); err != nil {
		os.Remove(id)
		return fmt.Errorf("\nUnable to rename file: " + err.Error())
	}

	fmt.Printf("%s successfully downloaded to file %s.", id, cleanfilename)

	return nil
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
