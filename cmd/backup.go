//
// Copyright (c) 2019 Sony Mobile Communications Inc.
// SPDX-License-Identifier: MIT
//

package cmd

import (
	"fmt"
	"ghorgs/gnet"
	"ghorgs/model"
	"ghorgs/utils"
	cmds "github.com/spf13/cobra"
	"os"
	"path"
	"regexp"
	"strings"
)

type backuper struct {
	quiet     bool
	n         int
	since     string
	names     []string
	outFolder string
	data      map[string]*model.Table
}

var (
	b         = &backuper{}
	backupCmd = &cmds.Command{
		Use:   "backup",
		Short: "Backup GitHub repositories according to given criteria.",
		Long: "Download GitHub repositories according to given criteria" +
			" and save a tar.gz file to a given folder.",
		Args: b.validateArgs,
		Run:  b.run,
	}
	backRepos       = model.Repos
	backReposFields = model.Repos.GetFields().(*model.RepositoryFields)
)

func init() {
	backupCmd.Flags().BoolP("quiet",
		"q",
		false,
		"DO NOT ask user for confirmation."+
			"(Use with care, e.g. in scripts where interaction is minimal or impossible.)")

	backupCmd.Flags().IntP("n",
		"n",
		0,
		`Number of repositories to backup.

* If --n is used together with --since, then the result is:
  "the number --n of repositories to backup --since point in time - whichever comes first."
* If used alone, then the result is:
  "the most active number of repositories to backup".

NOTE: It will be ignored if used with --repos.
`)

	backupCmd.Flags().StringP("since",
		"s",
		"",
		`Backup repositories active since this date (YYYY-MM-DD).

* If --since is used together with --n, then the result is:
  "the number --n of repositories to backup --since point in time - whichever comes first."
* If --since is used together with --repos, then the result is:
  "backup the repositories from --repos list if they have been active --since this point in time".
`)

	backupCmd.Flags().StringP("repos",
		"r",
		"",
		`Comma separated list of repositories to backup.

* Name can contain alphanumeric and special characters '_', '.' and '-'.
* If --repos is used with --since, then the result is:
  "back up the repositories from --repos list if they have been active --since this point in time.

NOTE: --n will be ignored if used with --repos.
`)

	backupCmd.Flags().StringP("out",
		"O",
		".",
		"Output folder where archives of repositories are recorded.")

	rootCmd.AddCommand(backupCmd)
}

func (b *backuper) addCache(c map[string]*model.Table) {
	b.data = c
}

func (b *backuper) validateArgs(c *cmds.Command, args []string) error {
	var err error
	b.quiet, err = c.Flags().GetBool("quiet")
	if err != nil {
		panic(err)
	}

	// Verify that the number of repos is a positive integer.
	b.n, err = c.Flags().GetInt("n")
	if err != nil {
		panic(err)
	}

	if b.n < 0 {
		return fmt.Errorf("Insert --n greater than 0.")
	}

	// Verify that the date is in format YYYY-MM-DD, starting from 1900-01-01
	b.since, err = c.Flags().GetString("since")
	if err != nil {
		panic(err)
	}
	if b.since != "" {
		matched, err := regexp.MatchString(`^(19|[2-9]\d)\d\d-(0?[1-9]|1[0-2])-(0?[1-9]|[12]\d|3[01])$`,
			b.since)
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
	backRepos, err := c.Flags().GetString("repos")
	if err != nil {
		panic(err)
	}
	if backRepos != "" {
		matched, err := regexp.MatchString(`^[\.|\-|\_|[:alnum:]]+(\,[\.|\-|\_|[:alnum:]]+)*$`,
			backRepos)
		if err != nil {
			return err
		}
		if !matched {
			return fmt.Errorf("--repos can only contain a comma separated list of repository" +
				" names written in ascii alpha-numeric characters ([._-] are allowed.)")
		}

		b.names = strings.Split(backRepos, ",")
		b.n = 0
	}

	if b.n == 0 && b.since == "" && len(b.names) == 0 {
		return fmt.Errorf("No criteria for archiving provided. Exiting.")
	}

	// Verify path of out folder.
	b.outFolder, err = c.Flags().GetString("out")
	if err != nil {
		panic(err)
	}
	if _, err := os.Stat(b.outFolder); os.IsNotExist(err) {
		return err
	}

	return nil
}

func (b *backuper) run(c *cmds.Command, args []string) {
	if gnet.Conf.User == "" || gnet.Conf.Token == "" {
		fmt.Println("Error! Invalid credentials.")
		return
	}

	// 0. get cache for repos
	ca, err := Cache([]model.Entity{backRepos})
	if err != nil {
		fmt.Println("Error!", err.Error())
		return
	}

	b.addCache(ca)

	// 2. if --repos set, get cache projection to --repos,
	var projection *model.Table
	if b.names != nil {
		projection, err = b.dataProjectionByName()
		if err != nil {
			fmt.Println(err.Error())
			if projection == nil {
				// nothing to work with so just return
				return
			}
		}
	} else {
		projection = b.data[backRepos.GetName()]
	}

	// 1. sort by `last updated` if "last n since" is requested,
	//    otherwise, keep unsorted, i.e. in order of original
	//    request from cli, e.g. for
	//       `ghorgs backup -r repo1,repo3,repo2`
	//    present as:
	//       id_repo1 repo1 type_repo1 ...
	//       id_repo3 repo3 type_repo3 ...
	//       id_repo2 repo2 type_repo2 ...
	if b.n > 0 || b.since != "" {
		_, err = projection.SortByField(backReposFields.Updated.Name)
		if err != nil {
			panic(err)
		}
	}

	if b.n > 0 {
		// if --n set, get copy of cache with --n most active
		projection, err = projection.Last(b.n)
		if err != nil {
			panic(err)
		}
	}

	// 3. get copy of cache cut by --since flag
	if b.since != "" {
		projection, err = projection.GreaterThanByField(backReposFields.Updated.Name, b.since)
		if err != nil {
			panic(err)
		}
	}

	if len(projection.Keys) == 0 {
		fmt.Println("There are no repositories with requested criteria. Exiting.")
		return
	}

	// 4. display the result to the user and request confirmation
	msg := "\nThe following repositories will be backed up "
	fmt.Printf(msg+"(%d):\n", len(projection.Keys))
	fmt.Printf("%s\n", projection)

	if !b.quiet && !utils.GetUserConfirmation() {
		return
	}

	// 5. iterate over result to:
	for _, key := range projection.Keys {
		//   5.0 git clone from url into -O
		rawurl := projection.Records[key][backReposFields.Url.Index]
		url, err := utils.Url(rawurl,
			gnet.Conf.User,
			gnet.Conf.Token)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Printf("Cloning `%s` to `%s` ...\n", rawurl, b.outFolder)
		repoName := projection.Records[key][backReposFields.Name.Index]
		err = utils.GitClone(url, b.outFolder, repoName)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		//   5.1 tar.gz the clone in -O
		clonePath := path.Join(b.outFolder, repoName)
		fmt.Printf("Creating archive '%s' in '%s'...\n",
			repoName+".tar.gz", b.outFolder)
		err = utils.TarGz(repoName, clonePath)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		//   5.2 compare tar -tvf with clone (compare size?)
		fmt.Printf("Archive '%s' created. Verifying...\n", repoName+".tar.gz")
		err = utils.TargzVerify(repoName, clonePath)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		//   5.3 rm clone in -O
		fmt.Printf("Removing %s...\n", clonePath)
		os.RemoveAll(path.Join(b.outFolder, repoName))
	} // for _, key := range projection.Keys {
}

func (b *backuper) dataProjectionByName() (*model.Table, error) {
	return b.data[backRepos.GetName()].FindAllByFieldValues(backReposFields.Name.Name, b.names)
}
