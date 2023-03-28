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
	//"errors"
	"github.com/spf13/cobra"
	"github.com/tlinden/cenophane/upctl/cfg"
	"github.com/tlinden/cenophane/upctl/lib"
	"os"
)

func FormCommand(conf *cfg.Config) *cobra.Command {
	var formCmd = &cobra.Command{
		Use:   "form {create|delete|modify|list}",
		Short: "Form commands",
		Long:  `Manage upload forms.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// errors at this stage do not cause the usage to be shown
			//cmd.SilenceUsage = true
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
			return nil
		},
	}

	formCmd.Aliases = append(formCmd.Aliases, "frm")
	formCmd.Aliases = append(formCmd.Aliases, "f")

	formCmd.AddCommand(FormCreateCommand(conf))

	return formCmd
}

func FormCreateCommand(conf *cfg.Config) *cobra.Command {
	var formCreateCmd = &cobra.Command{
		Use:   "create [options]",
		Short: "Create a new form",
		Long:  `Create a new form for consumers so they can upload something.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// errors at this stage do not cause the usage to be shown
			cmd.SilenceUsage = true

			return lib.CreateForm(os.Stdout, conf)
		},
	}

	// options
	formCreateCmd.PersistentFlags().StringVarP(&conf.Expire, "expire", "e", "", "Expire setting: asap or duration (accepted shortcuts: dmh)")
	formCreateCmd.PersistentFlags().StringVarP(&conf.Description, "description", "D", "", "Description of the form")
	formCreateCmd.PersistentFlags().StringVarP(&conf.Notify, "notify", "n", "", "Email address to get notified when consumer has uploaded files")

	formCreateCmd.Aliases = append(formCreateCmd.Aliases, "add")
	formCreateCmd.Aliases = append(formCreateCmd.Aliases, "+")

	return formCreateCmd
}
