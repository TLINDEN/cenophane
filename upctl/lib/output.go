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
	"github.com/olekukonko/tablewriter"
	"github.com/tlinden/cenophane/common"
	"io"
	"strings"
	"time"
)

// make a human readable version of the expire setting
func prepareExpire(expire string, start common.Timestamp) string {
	switch expire {
	case "asap":
		return "On first access"
	default:
		return time.Unix(start.Unix()+int64(common.Duration2int(expire)), 0).Format("2006-01-02 15:04:05")
	}

	return ""
}

// generic table writer
func WriteTable(w io.Writer, headers []string, data [][]string) {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetHeader(headers)
	table.AppendBulk(data)

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	table.Render()

	fmt.Fprintln(w, tableString.String())
}

/* Print output like psql \x

   Prints  all  Uploads  and  Forms which  exist  in  common.Response,
   however, we expect only one kind  of them to be actually filled, so
   the function can be used for forms and uploads.
*/
func WriteExtended(w io.Writer, response *common.Response) {
	format := fmt.Sprintf("%%%ds: %%s\n", Maxwidth)

	// we shall only have 1 element, however, if we ever support more, here we go
	for _, entry := range response.Uploads {
		expire := prepareExpire(entry.Expire, entry.Created)
		fmt.Fprintf(w, format, "Id", entry.Id)
		fmt.Fprintf(w, format, "Expire", expire)
		fmt.Fprintf(w, format, "Context", entry.Context)
		fmt.Fprintf(w, format, "Created", entry.Created)
		fmt.Fprintf(w, format, "Filename", entry.File)
		fmt.Fprintf(w, format, "Url", entry.Url)
		fmt.Fprintln(w)
	}

	// FIXME: add response.Forms loop here
}

// extract an common.Uploads{} struct from json response
func GetResponse(resp *req.Response) (*common.Response, error) {
	response := common.Response{}

	if err := json.Unmarshal([]byte(resp.String()), &response); err != nil {
		return nil, errors.New("Could not unmarshall JSON response: " + err.Error())
	}

	if !response.Success {
		return nil, errors.New(response.Message)
	}

	return &response, nil
}

// turn the Uploads{} struct into a table and print it
func UploadsRespondTable(w io.Writer, resp *req.Response) error {
	response, err := GetResponse(resp)
	if err != nil {
		return err
	}

	if response.Message != "" {
		fmt.Fprintln(w, response.Message)
	}

	// tablewriter
	data := [][]string{}
	for _, entry := range response.Uploads {
		data = append(data, []string{
			entry.Id, entry.Expire, entry.Context, entry.Created.Format("2006-01-02 15:04:05"),
		})
	}

	WriteTable(w, []string{"ID", "EXPIRE", "CONTEXT", "CREATED"}, data)

	return nil
}

// turn the Uploads{} struct into xtnd output and print it
func RespondExtended(w io.Writer, resp *req.Response) error {
	response, err := GetResponse(resp)
	if err != nil {
		return err
	}

	if response.Message != "" {
		fmt.Fprintln(w, response.Message)
	}

	WriteExtended(w, response)

	return nil
}
