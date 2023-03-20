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
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tlinden/cenophane/upctl/cfg"
	"github.com/tlinden/cenophane/upctl/lib"
)

func ListCommand(conf *cfg.Config) *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "list [options] [file ..]",
		Short: "List uploads",
		Long:  `List uploads.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// errors at this stage do not cause the usage to be shown
			cmd.SilenceUsage = true

			return lib.List(conf, args)
		},
	}

	// options
	listCmd.PersistentFlags().StringVarP(&conf.Apicontext, "apicontext", "", "", "Filter by given API context")

	listCmd.Aliases = append(listCmd.Aliases, "ls")
	listCmd.Aliases = append(listCmd.Aliases, "l")

	return listCmd
}
