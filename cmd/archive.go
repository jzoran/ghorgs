// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package cmd

import (
	"fmt"
	"ghorgs/cache"
	"ghorgs/entities"
	"ghorgs/utils"
	cmds "github.com/spf13/cobra"
	"log"
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

	// 0. get cache for repos
	a.addCache(Cache([]string{"repos"}, entities.EntityMap))

	// 2. if --repos set, get cache projection to --repos,
	var projection *cache.Table
	var err error
	if a.names != nil {
		projection, err = dataProjectionByName()
		if err != nil && projection == nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		projection = a.data["repos"]
	}

	// 1. sort by `last updated`
	_, err = projection.SortByField("Updated")
	if err != nil {
		panic(err)
	}

	if a.n > 0 {
		// if --n set, get copy of cache with --n least active
		projection, err = projection.Last(a.n)
		if err != nil {
			panic(err)
		}
	}

	// 3. get copy of cache cut by --since flag
	if a.since != "" {
		projection, err = projection.GreaterThanByField("Updated", a.since)
		if err != nil {
			panic(err)
		}
	}

	// 4. display the result to the user and request confirmation
	fmt.Println("\nThe following repositories will be removed from GitHub and archived:")
	fmt.Println(fmt.Sprintf("%s\n", projection.ToString()))
	// 5. iterate over result to:
	//   5.0 git clone from url into -O
	//   5.1 tar.gz the clone in -O
	//   5.2 compare tar -tvf with clone (compare size?)
	//   5.3 rm clone in -O
	//   5.4 rm repo in GitHub
}

func dataProjectionByName() (*cache.Table, error) {
	reposEntity := entities.EntityMap["repos"]
	projection := cache.MakeTable(reposEntity.GetTableFields())
	_, err := a.data["repos"].SortByField("Name")
	if err != nil {
		panic(err)
	}

	// find all a.names in the a.data set and add to projection
	ok := false
	for _, name := range a.names {
		t, err := a.data["repos"].FindByField("Name", name)
		if err == nil {
			ok = true
			key := t.Keys[0]
			record := t.Records[key]

			projection.AddKey(key)
			projection.AddRecord(key, record)
		} else {
			if utils.Debug.Verbose {
				log.Println(err.Error())
			}
			fmt.Println(fmt.Sprintf("`%s` not found in GitHub repositories.", name))
		}
	}
	if ok {
		return projection, nil
	}

	if len(projection.Keys) == 0 {
		// no requested repo was found
		return nil, fmt.Errorf("Errors found!")
	}

	// there were some repos found and some not
	return projection, fmt.Errorf("Errors found!")
}
