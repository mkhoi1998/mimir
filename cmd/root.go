package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mkhoi1998/devsup/utils"
)

var rootCmd = &cobra.Command{
	Use:  "devups",
	Long: `A command-line assistant for developers`,
	Run:  handleRootCmd,
}

func init() {
}

func handleRootCmd(cmd *cobra.Command, args []string) {
	fmt.Println(utils.DefaultHelpCmd(cmd))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
