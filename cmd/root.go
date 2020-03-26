package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:                   "devups",
	Long:                  `A command-line assistant for developers`,
	Args:                  cobra.MinimumNArgs(0),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ResponseHandler(args))
	},
}

// func init() {
// }

// Execute the commands
func Execute() {
	rootCmd.Execute()
}
