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
	"github.com/go-resty/resty/v2"
	"github.com/tlinden/up/upctl/cfg"
	"path/filepath"
)

type Response struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func Runclient(cfg *cfg.Config, args []string) error {
	if len(args) == 0 {
		return errors.New("No files specified to upload.")
	}

	client := resty.New()
	client.SetDebug(cfg.Debug)

	url := cfg.Endpoint + "/putfile"

	postfiles := make(map[string]string)

	for _, file := range args {
		postfiles[filepath.Base(file)] = file
	}

	resp, err := client.R().
		SetFiles(postfiles).
		SetFormData(map[string]string{"expire": "1d"}).
		Post(url)

	if err != nil {
		return err
	}

	if cfg.Debug {
		fmt.Println("Response Info:")
		fmt.Println("  Error      :", err)
		fmt.Println("  Status Code:", resp.StatusCode())
		fmt.Println("  Status     :", resp.Status())
		fmt.Println("  Proto      :", resp.Proto())
		fmt.Println("  Time       :", resp.Time())
		fmt.Println("  Received At:", resp.ReceivedAt())
		fmt.Println("  Body       :\n", resp)
		fmt.Println()
	}

	r := Response{}

	json.Unmarshal([]byte(resp.String()), &r)

	fmt.Println(r)

	return nil
}
