/*
Copyright © 2021 SUSE LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/


package cmd

import (
	"github.com/spf13/cobra"
	"github.com/itxaka/luet-mtree/pkg/action"
)

// checkCmd represents the check command
func newCheckCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "check [file or dir] [validation file]",
		Short: "Check a file or dir against a validation file",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				_ = cmd.Usage()
				return nil
			}

			// Just checking if the value is valid
			switch format {
				case "bsd", "json", "path":
				default:
					_ = cmd.Usage()
					return nil
			}

			checkAction := action.NewCheckAction(args[0], args[1], format)
			err := checkAction.Run()
			if err != nil {
				return err
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&format, "format", "f","bsd", "Format for output. Choices are bsd, path and json.")

 	return cmd
}