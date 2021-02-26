//
// Copyright (c) 2019 Sony Mobile Communications Inc.
// SPDX-License-Identifier: MIT
//

package cmd

import (
	"fmt"
	"ghorgs/gnet"
	"ghorgs/utils"
	cmds "github.com/spf13/cobra"
	flags "github.com/spf13/viper"
)

var rootCmd = &cmds.Command{
	Use:              "ghorgs",
	Short:            "ghorgs = GitHub ORGanizationS",
	Long:             "ghorgs = GitHub ORGanizationS\nA simple cli tool to manage organizational accounts on GitHub.",
	PersistentPreRun: initFlags,
}

func initFlags(c *cmds.Command, args []string) {
	gnet.Conf.User = flags.GetString("user")
	gnet.Conf.Token = flags.GetString("token")
	gnet.Conf.Organization = flags.GetString("organization")
	utils.Debug.Verbose = flags.GetBool("verbose")
	utils.Debug.DryRun = flags.GetBool("dry-run")
}

func init() {
	cmds.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolP("dry-run",
		"d",
		false,
		"Perform a dry run of the command without actually executing it in the end.")
	flags.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run")) // nolint

	rootCmd.PersistentFlags().BoolP("verbose",
		"v",
		false,
		"Toggle debug printouts.")
	flags.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")) // nolint

	rootCmd.PersistentFlags().StringP("user",
		"u",
		"",
		"User name of the owner of token. (Needed with 'git clone'.)")
	flags.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user")) // nolint

	rootCmd.PersistentFlags().StringP("token",
		"t",
		"",
		"Security token used on Github. Overrides the token from configuration file.\n"+
			"  Required GitHub scopes covered by a single token in the config file are:\n"+
			"    - user,\n"+
			"    - delete_repo,\n"+
			"    - public_repo,\n"+
			"    - repo,\n"+
			"    - repo_deployment,\n"+
			"    - repo:status,\n"+
			"    - read:repo_hook,\n"+
			"    - read:org,\n"+
			"    - read:public_key,\n"+
			"    - read:gpg_key.\n"+
			" Individual commands don't require all the scopes, so different tokens can be "+
			" used in the command line for different commands.")
	flags.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token")) // nolint

	rootCmd.PersistentFlags().StringP("organization",
		"o",
		"",
		"Organizational account on GitHub analyzed.")

	flags.BindPFlag("organization", rootCmd.PersistentFlags().Lookup("organization")) // nolint
}

func initConfig() {
	flags.AddConfigPath(gnet.ConfigPath)
	flags.SetConfigName(gnet.ConfigName)
	flags.SetConfigType(gnet.ConfigType)

	if err := flags.ReadInConfig(); err != nil {
		if _, ok := err.(flags.ConfigFileNotFoundError); ok {
			// ignore, issue warning
			fmt.Printf("Warning: config file not found.")
			return
		}

		panic(fmt.Errorf("Fatal config error: %s", err))
	}

	if err := flags.Unmarshal(&gnet.Conf); err != nil {
		panic(fmt.Errorf("Fatal config error: %s", err))
	}
}

func Execute() error {
	return rootCmd.Execute()
}
