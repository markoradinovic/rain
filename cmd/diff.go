package cmd

import (
	"fmt"
	"strings"

	"github.com/aws-cloudformation/rain/diff"
	"github.com/aws-cloudformation/rain/util"
	"github.com/awslabs/aws-cloudformation-template-formatter/parse"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff [left] [right]",
	Short: "Compare templates with other templates or stacks",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		leftFn, rightFn := args[0], args[1]

		left, err := parse.ReadFile(leftFn)
		if err != nil {
			util.Die(err)
		}

		right, err := parse.ReadFile(rightFn)
		if err != nil {
			util.Die(err)
		}

		output := diff.Format(diff.Compare(left, right))

		for _, line := range strings.Split(output, "\n") {
			colour := util.None

			switch {
			case strings.HasPrefix(line, ">>> "):
				colour = util.Green
			case strings.HasPrefix(line, "<<< "):
				colour = util.Red
			case strings.HasPrefix(line, "||| "):
				colour = util.Orange
			}

			fmt.Println(util.Text{line, colour})
		}
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
