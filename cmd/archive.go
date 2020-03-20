// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package cmd

import (
	"fmt"
	"ghorgs/cache"
	//	"ghorgs/entities"
	cmds "github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

type archiveT struct {
	n         int
	since     string
	names     []string
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
		`Number of repositories to archive.

* If --n is used together with --since, then the result is:
  "the number --n of repositories to archive --since point in time - whichever comes first."
* If used alone, then the result is:
  "the least active number of repositories to archive".

NOTE: It will be ignored if used with --repos.
`)

	archiveCmd.Flags().StringP("since",
		"s",
		"",
		`Remove repositories inactive since this date (YYYY-MM-DD).

* If --since is used together with --n, then the result is:
  "the number --n of repositories to archive --since point in time - whichever comes first."
* If --since is used together with --repos, then the result is:
  "archive the repositories from --repos list if they have been inactive --since this point in time".
`)

	archiveCmd.Flags().StringP("repos",
		"r",
		"",
		`Comma separated list of repositories to archive.

* Name can only contain alphanumeric characters.
* If --repos is used with --since, then the result is:
  "archive the repositories from --repos list if they have been inactive --since this point in time.

NOTE: --n will be ignored if used with --repos.
`)

	archiveCmd.Flags().StringP("out",
		"O",
		".",
		"Output folder where archives of repositories are recorded.")

	rootCmd.AddCommand(archiveCmd)

}

func (a *archiveT) addCache(c map[string]*cache.Table) {
	a.data = c
}

func (a *archiveT) validateArgs(c *cmds.Command, args []string) error {
	var err error

	// verify that the number of repos is a positive integer
	a.n, err = c.Flags().GetInt("n")
	if err != nil {
		panic(err)
	}

	if a.n <= 0 {
		return fmt.Errorf("Insert --n greater than 0.")
	}

	// verify that the date is in format YYYY-MM-DD, starting from 1900-01-01
	a.since, err = c.Flags().GetString("since")
	if err != nil {
		panic(err)
	}
	if a.since != "" {
		matched, err := regexp.MatchString(`^(19|[2-9]\d)\d\d-(0?[1-9]|1[0-2])-(0?[1-9]|[12]\d|3[01])$`,
			a.since)
		if err != nil {
			return err
		}
		if !matched {
			return fmt.Errorf("The date --since does not match format: YYYY-MM-DD " +
				"(starting from the 1900s)...")
		}
	}

	// verify that repos are a comma separated list of alphanumerics and
	// ignore number of repos to archive
	repos, err := c.Flags().GetString("repos")
	if err != nil {
		panic(err)
	}
	if repos != "" {
		// matched, err := regexp.MatchString(`^[[:alnum:]]+(\,[[:alnum:]]+)*$`, repos)
		matched, err := regexp.MatchString(`^[\.|\-|\_|[:alnum:]]+(\,[\.|\-|\_|[:alnum:]]+)*$`, repos)
		if err != nil {
			return err
		}
		if !matched {
			return fmt.Errorf("--repos can only contain a comma separated list of repository names written in ascii alpha-numeric characters.")
		}

		a.names = strings.Split(repos, ",")
		a.n = 0
	}

	// verify path
	a.outFolder, err = c.Flags().GetString("out")
	if err != nil {
		panic(err)
	}
	if _, err := os.Stat(a.outFolder); os.IsNotExist(err) {
		return err
	}

	return nil
}

func (a *archiveT) run(c *cmds.Command, args []string) {
	fmt.Println("TODO: implement archive...")
}
