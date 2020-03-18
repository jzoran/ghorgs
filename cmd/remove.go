// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package cmd

import (
	"fmt"
	"ghorgs/cache"
	//	"ghorgs/entities"
	cmds "github.com/spf13/cobra"
)

type removeT struct {
	n     int
	since string
	names string
	data  map[string]*cache.Table
}

var r = &removeT{}

var removeCmd = &cmds.Command{
	Use:   "remove",
	Short: "Remove GitHub users according to given criteria.",
	Long:  `Remove GitHub users according to given criteria.`,
	Args:  r.validateArgs,
	Run:   r.run,
}

func init() {
	removeCmd.Flags().IntP("n",
		"n",
		1,
		"Number of users to remove.")

	removeCmd.Flags().StringP("since",
		"s",
		"",
		"Remove users inactive since this date.")

	removeCmd.Flags().StringP("users",
		"u",
		"",
		"Comma separated list of users to remove.")

	rootCmd.AddCommand(removeCmd)

}

func (r *removeT) addCache(c map[string]*cache.Table) {
	r.data = c
}

func (r *removeT) validateArgs(c *cmds.Command, args []string) error {
	return nil
}

func (r *removeT) run(c *cmds.Command, args []string) {
	fmt.Println("TODO: implement remove...")
}
