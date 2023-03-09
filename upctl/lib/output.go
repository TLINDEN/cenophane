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
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"time"
)

func WriteTable(headers []string, data [][]string) error {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader(headers)
	table.AppendBulk(data)

	// for _, row := range data.entries {
	//      table.Append(trimRow(row))
	// }

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

	return nil
}

func prepareExpire(expire string, start Timestamp) string {
	switch expire {
	case "asap":
		return "On first access"
	default:
		return time.Unix(start.Unix()+int64(duration2int(expire)), 0).Format("2006-01-02 15:04:05")
	}

	return ""
}

func WriteExtended(uploads *Uploads) {
	format := fmt.Sprintf("%%%ds: %%s\n", Maxwidth)

	// we shall only have 1 element, however, if we ever support more, here we go
	for _, entry := range uploads.Entries {
		expire := prepareExpire(entry.Expire, entry.Uploaded)
		fmt.Printf(format, "Id", entry.Id)
		fmt.Printf(format, "Expire", expire)
		fmt.Printf(format, "Context", entry.Context)
		fmt.Printf(format, "Uploaded", entry.Uploaded)
		fmt.Printf(format, "Filename", entry.File)
		fmt.Println()
	}
}
