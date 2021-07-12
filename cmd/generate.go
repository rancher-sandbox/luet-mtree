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
	"errors"
	"github.com/itxaka/luet-mtree/pkg/action"
	"github.com/itxaka/luet-mtree/pkg/log"
	"github.com/spf13/cobra"
	"os"
)

// generateCmd represents the generate command
func newGenerateCmd() *cobra.Command {
	var outputFile string
	var keywords []string
	var overrideKeywords []string

	cmd := &cobra.Command{
		Use:   "generate [file or dir]",
		Short: "generate a checksum file for the file or dir provided",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				_ = cmd.Usage()
				return nil
			}

			if len(keywords) > 0 && len(overrideKeywords) > 0 {
				log.Log("Cannot use both -k and -K flags. Either add keywords or override the default ones.")
				return errors.New("cannot use both -k and -K flags. Either add keywords or override the default ones")
			}
			generateAction := action.NewGenerateAction(args[0], outputFile, keywords, overrideKeywords)
			err := generateAction.Run()
			if err != nil {
				log.Log(err.Error())
				os.Exit(1)
			}
			log.Log("Generation for %s done!", args[0])
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVarP(&outputFile, "output", "o", "", "Name for output file, otherwise it defaults to stdout")
	f.StringSliceVarP(&keywords, "keywords", "k", []string{}, "Add keywords to default ones (type, sha512digest)")
	f.StringSliceVarP(&overrideKeywords, "overridekeywords", "K", []string{}, "Override the default keyworkds (type, sha512digest)")
	return cmd
}
