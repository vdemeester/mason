package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print out the version of mason",
	Long:  `Print out the version of mason.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("0.0.0-dev")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
