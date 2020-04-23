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
	"path"
	"regexp"
	"strings"
)

type remover struct {
	quiet   bool
	mfa     bool
	company bool
	access  bool
	names   []string
	data    map[string]*model.Table
}

var (
	r         = &remover{}
	removeCmd = &cmds.Command{
		Use:   "remove",
		Short: "Remove GitHub users according to given criteria.",
		Long:  `Remove GitHub users according to given criteria.`,
		Args:  r.validateArgs,
		Run:   r.run,
	}
	users       = model.Users
	usersFields = model.Users.GetFields().(*model.UsersFields)
)

func init() {
	removeCmd.Flags().BoolP("quiet",
		"q",
		false,
		"DO NOT ask user for confirmation. "+
			"(Use with care, e.g. in scripts where interaction is minimal or impossible.)")

	removeCmd.Flags().BoolP("2FA",
		"A",
		false,
		"Remove users without 2FA set up.")

	removeCmd.Flags().BoolP("company",
		"c",
		false,
		"Remove users without company affiliation.")

	removeCmd.Flags().BoolP("access",
		"a",
		false,
		"Remove users without access to any repository owned by the organization.")

	removeCmd.Flags().StringP("users",
		"U",
		"",
		"Comma separated list of users to remove. "+
			"Name can contain alphanumeric and special characters '_', '.' and '-'.")

	rootCmd.AddCommand(removeCmd)

}

func (r *remover) addCache(c map[string]*model.Table) {
	r.data = c
}

func (r *remover) validateArgs(c *cmds.Command, args []string) error {
	var err error
	r.quiet, err = c.Flags().GetBool("quiet")
	if err != nil {
		panic(err)
	}

	r.mfa, err = c.Flags().GetBool("2FA")
	if err != nil {
		panic(err)
	}

	r.company, err = c.Flags().GetBool("company")
	if err != nil {
		panic(err)
	}

	r.access, err = c.Flags().GetBool("access")
	if err != nil {
		panic(err)
	}

	// Verify that users are a comma separated list of alphanumerics and
	// special characters '.', '_' and '-'.
	// Ignore other criteria.
	users, err := c.Flags().GetString("users")
	if err != nil {
		panic(err)
	}
	if users != "" {
		matched, err := regexp.MatchString(`^[\.|\-|\_|[:alnum:]]+(\,[\.|\-|\_|[:alnum:]]+)*$`, users)
		if err != nil {
			return err
		}
		if !matched {
			return fmt.Errorf("--users can only contain a comma separated list of usernames " +
				"written in ascii alpha-numeric characters ([._-] are allowed.).")
		}

		r.names = strings.Split(users, ",")
		r.mfa = false
		r.company = false
		r.access = false
	}

	return nil
}

func (r *remover) run(c *cmds.Command, args []string) {

	if gnet.Conf.User == "" || gnet.Conf.Token == "" {
		fmt.Println("Error! Invalid credentials.")
		return
	}

	// 0. get cache for users
	ca, err := Cache([]model.Entity{users})
	if err != nil {
		fmt.Println("Error! %s", err.Error())
		return
	}

	r.addCache(ca)

	// 2. if --users set, get cache projection to --users,
	var projection *model.Table
	if r.names != nil {
		projection, err = r.dataProjectionByName()
		if err != nil {
			fmt.Println(err.Error())
			if projection == nil {
				// nothing to work with so just return
				return
			}
		}
	} else {
		projection = r.data[users.GetName()]
	}

	// 2FA, Company affiliation and Accessible repositories
	// criteria are combined with AND operation.
	// (Note: if r.names == true,
	//       then r.mfa == r.company == r.access == false)

	// 1. check by 2FA
	if r.mfa {
		tmp, err := projection.FindAllByField(usersFields.MFA.Name, "false")
		if err != nil {
			fmt.Println(err.Error())
			// allow partial results, so don't return
		}
		if tmp == nil {
			// nothing to work with so return here
			return
		}

		projection = tmp
	}

	// 2. check by company affiliation
	if r.company {
		tmp, err := projection.FindAllByField(usersFields.Company.Name, "")
		if err != nil {
			fmt.Println(err.Error())
			// allow partial results, so don't return
		}
		if tmp == nil {
			// nothing to work with so return here
			return
		}

		projection = tmp
	}

	// 3. check by accessible repositories
	if r.access {
		tmp, err := projection.FindAllByField(usersFields.Repositories.Name, "0")
		if err != nil {
			fmt.Println(err.Error())
			// allow partial results, so don't return
		}
		if tmp == nil {
			// nothing to work with so return here
			return
		}

		projection = tmp
	}

	if projection == nil {
		// nothing to work with so just return
		return
	}

	// 4. display the result to the user and request confirmation
	fmt.Println(
		fmt.Sprintf("\nThe following users will be removed from the organization (%d):",
			len(projection.Keys)))
	fmt.Println(fmt.Sprintf("%s\n", projection.ToString()))

	if !r.quiet && !utils.GetUserConfirmation() {
		return
	}

	// 5. iterate over the result to remove the users
	for _, key := range projection.Keys {
		userLogin := projection.Records[key][usersFields.Login.Index]
		// create GitHub v3 request to delete a user:
		//     DELETE /orgs/:org/members/:username
		rmRequest := gnet.MakeGitHubV3Request(http.MethodDelete,
			path.Join("orgs",
				gnet.Conf.Organization,
				"members",
				userLogin),
			gnet.Conf.Token)
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
	}
}

func (r *remover) dataProjectionByName() (*model.Table, error) {
	return r.data[users.GetName()].FindAllByFieldValues(usersFields.Login.Name, r.names)
}
