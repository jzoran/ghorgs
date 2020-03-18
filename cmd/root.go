// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

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
	gnet.Conf.Token = flags.GetString("token")
	gnet.Conf.Organization = flags.GetString("organization")
	utils.Debug.Verbose = flags.GetBool("verbose")
}

func init() {
	cmds.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolP("verbose",
		"v",
		false,
		"Toggle debug printouts.")
	flags.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().StringP("token",
		"t",
		"",
		"Security token used on Github.\n"+
			"  Required GitHub scopes covered by token are:\n"+
			"    - user,\n"+
			"    - public_repo,\n"+
			"    - repo,\n"+
			"    - repo_deployment,\n"+
			"    - repo:status,\n"+
			"    - read:repo_hook,\n"+
			"    - read:org,\n"+
			"    - read:public_key,\n"+
			"    - read:gpg_key")
	flags.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))

	rootCmd.PersistentFlags().StringP("organization",
		"o",
		"",
		"Organizational account on GitHub analyzed.")

	flags.BindPFlag("organization", rootCmd.PersistentFlags().Lookup("organization"))
}

func initConfig() {
	flags.SetConfigName("config")
	flags.SetConfigType("yaml")
	flags.AddConfigPath("./config")

	if err := flags.ReadInConfig(); err != nil {
		if _, ok := err.(flags.ConfigFileNotFoundError); ok {
			// ignore, issue warning
			fmt.Printf("Warning: config file not found.")
		} else {
			panic(fmt.Errorf("Fatal config error: %s", err))
		}
	} else {
		err = flags.Unmarshal(&gnet.Conf)
		if err != nil {
			panic(fmt.Errorf("Fatal config error: %s", err))
		}
	}
}

func Execute() error {
	return rootCmd.Execute()
}
