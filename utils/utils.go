package utils

import (
	"fmt"

	"github.com/spf13/cobra"
)

func DefaultHelpCmd(cmd *cobra.Command) string {
	return fmt.Sprintf("%v\n\n%v", cmd.Long, cmd.UsageString())
}
