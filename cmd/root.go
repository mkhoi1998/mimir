package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:                   "devups",
	Long:                  `A command-line assistant for developers`,
	Args:                  cobra.MinimumNArgs(0),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		ans := ResponseHandler(args)
		fmt.Println(ans)
	},
}

func init() {
}

// Execute the commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
