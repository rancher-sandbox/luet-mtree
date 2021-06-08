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
	"github.com/itxaka/luet-mtree/pkg/action"
	"github.com/itxaka/luet-mtree/pkg/log"
	"github.com/spf13/cobra"
	"os"
)

// checkCmd represents the check command
func newCheckCmd() *cobra.Command {
	var format string
	var exclude []string
	cmd := &cobra.Command{
		Use:           "check [file or dir] [validation file]",
		Short:         "Check a file or dir against a validation file",
		SilenceUsage:  true,
		SilenceErrors: true,
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

			checkAction := action.NewCheckAction(args[0], args[1], format, exclude)
			out, err := checkAction.Run()

			if err != nil {
				if out != "" {
					// The called decides to write it to standard out or not instead of the package
					_, _ = os.Stdout.Write([]byte(out))
				}
				log.Log(err.Error())
				os.Exit(1)
			}
			log.Log("Check for %s with validation file %s done!", args[0], args[1])
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&format, "format", "f", "bsd", "Format for output. Choices are bsd, path and json.")
	f.StringSliceVarP(&exclude, "exclude", "x", make([]string, 0), "Exclude paths from check. Checks against the path prefix, so 'oem/' will cover both 'oem/' and 'oem/features/' paths.")

	return cmd
}
