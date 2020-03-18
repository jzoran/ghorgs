// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package cmd

import (
	"fmt"
	"ghorgs/cache"
	//	"ghorgs/entities"
	cmds "github.com/spf13/cobra"
)

type archiveT struct {
	n         int
	since     string
	names     string
	outFolder string
	data      map[string]*cache.Table
}

var a = &archiveT{}

var archiveCmd = &cmds.Command{
	Use:   "archive",
	Short: "Archive GitHub repositories according to given criteria.",
	Long:  `Remove GitHub repositories according to given criteria and archive to a given folder.`,
	Args:  a.validateArgs,
	Run:   a.run,
}

func init() {
	archiveCmd.Flags().IntP("n",
		"n",
		1,
		"Number of repositories to archive.")

	archiveCmd.Flags().StringP("since",
		"s",
		"",
		"Remove repositories inactive since this date.")

	archiveCmd.Flags().StringP("repos",
		"r",
		"",
		"Comma separated list of repositories to archive.")

	archiveCmd.Flags().StringP("out",
		"O",
		"",
		"Output folder where archives of repositories are recorded..")

	rootCmd.AddCommand(archiveCmd)

}

func (a *archiveT) addCache(c map[string]*cache.Table) {
	a.data = c
}

func (a *archiveT) validateArgs(c *cmds.Command, args []string) error {
	return nil
}

func (a *archiveT) run(c *cmds.Command, args []string) {
	fmt.Println("TODO: implement archive...")
}
