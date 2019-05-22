package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws-cloudformation/rain/client/cfn"
	"github.com/aws-cloudformation/rain/util"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm [stack name]",
	Short: "Delete a stack",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stackName := args[0]

		err := cfn.DeleteStack(stackName)
		if err != nil {
			util.Die(err)
		}

		for {
			stack, err := cfn.GetStack(stackName)
			if err != nil {
				util.Die(err)
			}

			outputStack(stack, true)

			message := ""

			status := string(stack.StackStatus)

			switch {
			case status == "DELETE_COMPLETE":
				message = "Successfully deleted " + stackName
			case strings.HasSuffix(status, "_COMPLETE") || strings.HasSuffix(status, "_FAILED"):
				message = "Failed to delete " + stackName
			}

			if message != "" {
				fmt.Println()
				fmt.Println(message)
				return
			}

			time.Sleep(2 * time.Second)
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
