package cmd

import (
	"fmt"
	"github.com/itxaka/luet-mtree/internal/version"
	"github.com/spf13/cobra"
)


func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "version",
		Args: cobra.ExactArgs(0),
		Short: "Print the version",
		RunE: func(cmd *cobra.Command, args []string) error {
			v := version.Get()
			if cmd.Flag("long").Changed {
				fmt.Printf("%#v", v)
			} else {
				fmt.Printf("%s+g%s", v.Version, v.GitCommit[:7])
			}

			return nil
		},
	}
	f := cmd.Flags()
	f.Bool("long", false,"Show long version info")
	return cmd
}
