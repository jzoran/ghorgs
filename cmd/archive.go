// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package cmd

import (
	"fmt"
	"ghorgs/gnet"
	"ghorgs/model"
	"ghorgs/utils"
	cmds "github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

type archiver struct {
	quiet     bool
	n         int
	since     string
	names     []string
	backup    bool
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
	archiveCmd.Flags().BoolP("quiet",
		"q",
		false,
		"DO NOT ask user for confirmation."+
			"(Use with care, e.g. in scripts where interaction is minimal or impossible.)")

	archiveCmd.Flags().IntP("n",
		"n",
		0,
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

* Name can contain alphanumeric and special characters '_', '.' and '-'.
* If --repos is used with --since, then the result is:
  "archive the repositories from --repos list if they have been inactive --since this point in time.

NOTE: --n will be ignored if used with --repos.
`)

	archiveCmd.Flags().BoolP("backup",
		"b",
		false,
		"Only backup the repositories. DO NOT REMOVE them.")

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
	a.quiet, err = c.Flags().GetBool("quiet")
	if err != nil {
		panic(err)
	}

	// Verify that the number of repos is a positive integer.
	a.n, err = c.Flags().GetInt("n")
	if err != nil {
		panic(err)
	}

	if a.n < 0 {
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

	if a.n == 0 && a.since == "" && len(a.names) == 0 {
		return fmt.Errorf("No criteria for archiving provided. Exiting.")
	}

	a.backup, err = c.Flags().GetBool("backup")
	if err != nil {
		panic(err)
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

	if gnet.Conf.User == "" || gnet.Conf.Token == "" {
		fmt.Println("Error! Invalid credentials.")
		return
	}

	// 0. get cache for repos
	ca, err := Cache([]model.Entity{repos})
	if err != nil {
		fmt.Println("Error! %s", err.Error())
		return
	}

	a.addCache(ca)

	// 2. if --repos set, get cache projection to --repos,
	var projection *model.Table
	if a.names != nil {
		projection, err = a.dataProjectionByName()
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

	// 1. sort by `last updated` if "last n since" is requested,
	//    otherwise, keep unsorted, i.e. in order of original
	//    request from cli, e.g. for
	//       `ghorgs archive -r repo1,repo3,repo2`
	//    present as:
	//       id_repo1 repo1 type_repo1 ...
	//       id_repo3 repo3 type_repo3 ...
	//       id_repo2 repo2 type_repo2 ...
	if a.n > 0 || a.since != "" {
		_, err = projection.SortByField(reposFields.Updated.Name)
		if err != nil {
			panic(err)
		}
	}

	if a.n > 0 {
		// if --n set, get copy of cache with --n least active
		projection, err = projection.First(a.n)
		if err != nil {
			panic(err)
		}
	}

	// 3. get copy of cache cut by --since flag
	if a.since != "" {
		projection, err = projection.LessThanByField(reposFields.Updated.Name, a.since)
		if err != nil {
			panic(err)
		}
	}

	if len(projection.Keys) == 0 {
		fmt.Println("There are no repositories with requested criteria.Exiting.")
		return
	}

	// 4. display the result to the user and request confirmation
	msg := "\nThe following repositories will be "
	if a.backup {
		msg += "backed up "
	} else {
		msg += "removed from GitHub and archived "
	}

	fmt.Println(msg +
		fmt.Sprintf("(%d):",
			len(projection.Keys)))
	fmt.Println(fmt.Sprintf("%s\n", projection.ToString()))

	if !a.quiet && !utils.GetUserConfirmation() {
		return
	}

	// 5. iterate over result to:
	for _, key := range projection.Keys {
		//   5.0 git clone from url into -O
		rawurl := projection.Records[key][reposFields.Url.Index]
		url, err := utils.Url(rawurl,
			gnet.Conf.User,
			gnet.Conf.Token)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(fmt.Sprintf("Cloning `%s` to `%s` ...", rawurl, a.outFolder))
		repoName := projection.Records[key][reposFields.Name.Index]
		err = utils.GitClone(url, a.outFolder, repoName)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		//   5.1 tar.gz the clone in -O
		clonePath := path.Join(a.outFolder, repoName)
		destArchive := path.Join(repoName+".tar.gz", a.outFolder)
		fmt.Println(fmt.Sprintf("Creating archive '%s' in '%s'...",
			destArchive, a.outFolder))
		err = utils.TarGz(destArchive, clonePath)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		//   5.2 compare tar -tvf with clone (compare size?)
		fmt.Println(fmt.Sprintf("Archive '%s' created. Verifying...", destArchive))
		err = utils.TargzVerify(destArchive, clonePath)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		//   5.3 rm clone in -O
		fmt.Println(fmt.Sprintf("Removing %s...", clonePath))
		os.RemoveAll(path.Join(a.outFolder, repoName))

		//   5.4 if only backup, that's it, we're done
		if a.backup {
			continue
		}

		// 5.5 otherwise, rm repo in GitHub
		rmRequest := gnet.MakeGitHubV3Request(http.MethodDelete,
			path.Join(repos.GetName(),
				gnet.Conf.Organization,
				repoName),
			gnet.Conf.Token)
		if utils.Debug.DryRun {
			fmt.Println(fmt.Sprintf("Executing: %s %s ", rmRequest.Url, rmRequest.Method))
		} else {
			resp, status := rmRequest.Execute()
			if utils.Debug.Verbose {
				log.Print(resp)
			}
			// check response for error:
			// - `Status: 204 No Content` is OK
			// - `Status: 403 Forbidden` - abort since Token doesn't have Delete rights
			// - Any other code, continue
			if status.Code == http.StatusForbidden {
				fmt.Println("Error! HttpResponse:", status.Status)
				fmt.Println("Token is not allowed to delete repository.")
				return
			}
			if status.Code != http.StatusOK && status.Code != http.StatusNoContent {
				fmt.Println("Error! HttpResponse:", status.Status)
				continue
			}
		} // dry run
	} // for _, key := range projection.Keys {
}

func (a *archiver) dataProjectionByName() (*model.Table, error) {
	return a.data[repos.GetName()].FindAllByFieldValues(reposFields.Name.Name, a.names)
}
