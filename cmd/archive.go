// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package cmd

import (
	"fmt"
	"ghorgs/model"
	"ghorgs/utils"
	cmds "github.com/spf13/cobra"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

type archiver struct {
	n         int
	since     string
	names     []string
	outFolder string
	data      map[string]*model.Table
}

var (
	a          = &archiver{}
	archiveCmd = &cmds.Command{
		Use:   "archive",
		Short: "Archive GitHub repositories according to given criteria.",
		Long:  `Remove GitHub repositories according to given criteria and archive to a given folder.`,
		Args:  a.validateArgs,
		Run:   a.run,
	}
	repos       = model.Repos
	reposFields = model.Repos.GetFields().(*model.RepositoryFields)
)

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

func (a *archiver) addCache(c map[string]*model.Table) {
	a.data = c
}

func (a *archiver) validateArgs(c *cmds.Command, args []string) error {
	var err error

	// Verify that the number of repos is a positive integer.
	a.n, err = c.Flags().GetInt("n")
	if err != nil {
		panic(err)
	}

	if a.n <= 0 {
		return fmt.Errorf("Insert --n greater than 0.")
	}

	// Verify that the date is in format YYYY-MM-DD, starting from 1900-01-01
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

	// Verify that repos are a comma separated list of alphanumerics and
	// special characters '.', '_' and '-'.
	// Ignore number of repos to archive.
	repos, err := c.Flags().GetString("repos")
	if err != nil {
		panic(err)
	}
	if repos != "" {
		matched, err := regexp.MatchString(`^[\.|\-|\_|[:alnum:]]+(\,[\.|\-|\_|[:alnum:]]+)*$`, repos)
		if err != nil {
			return err
		}
		if !matched {
			return fmt.Errorf("--repos can only contain a comma separated list of repository names written in ascii alpha-numeric characters ([._-] are allowed.).")
		}

		a.names = strings.Split(repos, ",")
		a.n = 0
	}

	// Verify path of out folder.
	a.outFolder, err = c.Flags().GetString("out")
	if err != nil {
		panic(err)
	}
	if _, err := os.Stat(a.outFolder); os.IsNotExist(err) {
		return err
	}

	return nil
}

func (a *archiver) run(c *cmds.Command, args []string) {
	fmt.Println("TODO: implement archive...")

	// 0. get cache for repos
	a.addCache(Cache([]model.Entity{repos}))

	// 2. if --repos set, get cache projection to --repos,
	var projection *model.Table
	var err error
	if a.names != nil {
		projection, err = dataProjectionByName()
		if err != nil {
			fmt.Println(err.Error())
			if projection == nil {
				// nothing to work with so just return
				return
			}
		}
	} else {
		projection = a.data[repos.GetName()]
	}

	// 1. sort by `last updated`
	_, err = projection.SortByField(reposFields.Updated.Name)
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
		projection, err = projection.GreaterThanByField(reposFields.Updated.Name, a.since)
		if err != nil {
			panic(err)
		}
	}

	// 4. display the result to the user and request confirmation
	fmt.Println("\nThe following repositories will be removed from GitHub and archived:")
	fmt.Println(fmt.Sprintf("%s\n", projection.ToString()))
	fmt.Println("Are you sure you want to continue? (y/N):")

	var y string
	n, err := fmt.Scanf("%s", &y)
	if err != nil {
		fmt.Println(n)
		fmt.Println(err.Error())
	}
	if y != "y" && y != "yes" && y != "yep" && y != "Y" && y != "Yes" && y != "Sure thing, mate! Please do carry on." {
		fmt.Println("OK, aborting...")
		return
	}

	// 5. iterate over result to:
	for _, key := range projection.Keys {
		//   5.0 git clone from url into -O
		url := projection.Records[key][reposFields.Url.Index]
		fmt.Println(fmt.Sprintf("Cloning `%s` to `%s` ...", url, a.outFolder))
		repoName := projection.Records[key][reposFields.Name.Index]
		err = utils.GitClone(url, a.outFolder, repoName)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		//   5.1 tar.gz the clone in -O
		clonePath := path.Join(a.outFolder, repoName)
		err = utils.TarGz(repoName, clonePath, a.outFolder)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		//   5.2 compare tar -tvf with clone (compare size?)
		//   5.3 rm clone in -O
		fmt.Println(fmt.Sprintf("Removing %s...", clonePath))
		os.RemoveAll(path.Join(a.outFolder, repoName))
		//   5.4 rm repo in GitHub
	}
}

func dataProjectionByName() (*model.Table, error) {
	projection := repos.MakeTable()
	_, err := a.data[repos.GetName()].SortByField(reposFields.Name.Name)
	if err != nil {
		panic(err)
	}

	// find all a.names in the a.data set and add to projection
	ok := true
	for _, name := range a.names {
		t, err := a.data[repos.GetName()].FindByField(reposFields.Name.Name, name)
		if err == nil {
			key := t.Keys[0]
			record := t.Records[key]

			projection.AddKey(key)
			projection.AddRecord(key, record)
		} else {
			ok = false
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
