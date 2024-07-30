package tt

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version        = "dev"
	CommitHash     = "n/a"
	BuildTimestamp = "n/a"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version of tt",
	Long:  `Display the version of tt`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s-%s (%s)\n", Version, CommitHash, BuildTimestamp)
	},
}
