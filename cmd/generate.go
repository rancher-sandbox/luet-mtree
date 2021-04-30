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
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
func newGenerateCmd() *cobra.Command {
	var outputFile string
	var keywords []string

	cmd := &cobra.Command{
		Use:   "generate [file or dir]",
		Short: "generate a checksum file for the file or dir provided",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				_ = cmd.Usage()
				return nil
			}
			generateAction := action.NewGenerateAction(args[0], outputFile, keywords)
			err := generateAction.Run()
			if err != nil {
				return err
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVarP(&outputFile, "output", "o", "", "Name for output file, otherwise it defaults to stdout")
	f.StringSliceVarP(&keywords, "keywords", "k", []string{}, "Keywords to use to generate the tree (sha256 will automatically be added)")
	return cmd
}
