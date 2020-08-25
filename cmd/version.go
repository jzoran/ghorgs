//
// Copyright (c) 2019 Sony Mobile Communications Inc.
// SPDX-License-Identifier: MIT
//

package cmd

import (
	"fmt"
	cmds "github.com/spf13/cobra"
)

const version = "1.5.1"

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
