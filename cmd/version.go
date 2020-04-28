package cmd

import (
	"fmt"
	cmds "github.com/spf13/cobra"
)

const version = "1.2.1"

var versionCmd = &cmds.Command{
	Use:   "version",
	Short: "prints version of ghorgs tool",
	Long:  `print version of ghorgs tool`,
	Args:  cmds.NoArgs,
	Run: func(c *cmds.Command, args []string) {
		fmt.Println(version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
