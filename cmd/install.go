package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use: "install",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("see https://github.com/jupyter/echo_kernel/blob/master/echo_kernel/install.py")
	},
}

func init() {
	RootCmd.AddCommand(installCmd)

}
