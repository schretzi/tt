package tt

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func checkArguments(actual int, target int, cmd *cobra.Command) {
	if actual != target {
		fmt.Println(CharError, "Wrong number of Arguments")
		fmt.Println("")
		cmd.Help()
		os.Exit(1)
	}
}
