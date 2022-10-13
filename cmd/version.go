package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = ""

func newVersionCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:          "version",
		Short:        "Displays d4sva binary version",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", version)
		},
	}
}

func init() {
	rootCmd.AddCommand(newVersionCmd(version))
}
