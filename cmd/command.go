//
// Copyright (c) 2019 Sony Mobile Communications Inc.
// SPDX-License-Identifier: MIT
//

package cmd

import (
	"ghorgs/model"
	cmds "github.com/spf13/cobra"
)

type commander interface {
	addCache(c map[string]*model.Table)
	validateArgs(c *cmds.Command, args []string) error
	run(c *cmds.Command, args []string)
}
