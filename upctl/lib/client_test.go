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
	//"github.com/alecthomas/repr"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/tlinden/ephemerup/upctl/cfg"
	"net/http"
	"testing"
)

const endpoint string = "http://localhost:8080/api/v1"

type Unit struct {
	name     string
	apikey   string   // set to something else than "token" to fail auth
	wantfail bool     // true: expect to fail
	files    []string // path relative to ./t/
	sendcode int      // for httpmock
	sendjson string   // struct to respond with
	route    string   // dito
	method   string   // method to use
}

// simulate our ephemerup server
func Intercept(tt Unit) {
	httpmock.RegisterResponder(tt.method, endpoint+tt.route,
		func(request *http.Request) (*http.Response, error) {
			respbody := fmt.Sprintf(tt.sendjson)
			resp := httpmock.NewStringResponse(tt.sendcode, respbody)
			resp.Header.Set("Content-Type", "application/json; charset=utf-8")
			return resp, nil
		})
}

func TestUploadFiles(t *testing.T) {
	conf := &cfg.Config{
		Mock:     true,
		Apikey:   "token",
		Endpoint: endpoint,
		Silent:   true,
	}

	tests := []Unit{
		{
			name:     "upload-file",
			apikey:   "token",
			wantfail: false,
			route:    "/file/",
			sendcode: 200,
			sendjson: `{"success": true}`,
			files:    []string{"../t/t1"}, // pwd is lib/ !
			method:   "POST",
		},
		{
			name:     "upload-nonexistent-file",
			apikey:   "token",
			wantfail: true,
			route:    "/file/",
			sendcode: 200,
			sendjson: `{"success": false}`,
			files:    []string{"../t/none"},
			method:   "POST",
		},
		{
			name:     "upload-unauth",
			apikey:   "token",
			wantfail: true,
			route:    "/file/",
			sendcode: 403,
			sendjson: `{"success": false}`,
			files:    []string{"../t/t1"},
			method:   "POST",
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("UploadFiles-%s-%t", tt.name, tt.wantfail)
		Intercept(tt)
		err := UploadFiles(conf, tt.files)

		if err != nil && !tt.wantfail {
			t.Errorf("%s failed! wantfail: %t, error: %s", testname, tt.wantfail, err.Error())
		}
	}
}

func TestList(t *testing.T) {
	conf := &cfg.Config{
		Mock:     true,
		Apikey:   "token",
		Endpoint: endpoint,
		Silent:   true,
	}

	listing := `{"uploads":[{"id":"c8dh","expire":"asap","file":"t1","members":["t1"],"uploaded":1679318969.6434112,"context":"foo","url":""}],"success":true,"message":"","code":200}`
	tests := []Unit{
		{
			name:     "list",
			apikey:   "token",
			wantfail: false,
			route:    "/list/",
			sendcode: 200,
			sendjson: listing,
			files:    []string{},
			method:   "GET",
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("List-%s-%t", tt.name, tt.wantfail)
		Intercept(tt)
		err := List(conf, []string{})

		if err != nil && !tt.wantfail {
			t.Errorf("%s failed! wantfail: %t, error: %s", testname, tt.wantfail, err.Error())
		}
	}

}
