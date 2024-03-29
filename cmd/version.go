package cmd

import (
	"os"

	"github.com/goliatone/lgr/pkg/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version tag of lgr",
	// Long:  `All software has versions. This is Hugo's`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return version.Print(os.Stdout)
	},
}
