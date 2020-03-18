// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package cmd

import (
	"ghorgs/cache"
	cmds "github.com/spf13/cobra"
)

type commandI interface {
	addCache(c map[string]*cache.Table)
	validateArgs(c *cmds.Command, args []string) error
	run(c *cmds.Command, args []string)
}
