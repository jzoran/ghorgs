// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package main

import (
	"errors"
	"flag"
	"fmt"
	"ghorgs/commands"
	"ghorgs/entities"
	"ghorgs/gnet"
	"ghorgs/utils"
	"log"
	"os"
	"strings"
)

type Args struct {
	Help         bool
	Verbose      bool
	Cmd          string
	Entity       string
	By           string
	N            int
	Since        string
	Names        string
	Token        string
	Organization string
}

var args Args
var entityMap map[string]entities.Entity

func init() {
	repos := &entities.ReposResponse{}
	users := &entities.UsersResponse{}
	entityMap = map[string]entities.Entity{
		repos.GetName(): repos,
		users.GetName(): users,
	}

	flag.BoolVar(&args.Help, "h", false, "Prints this help.")
	flag.BoolVar(&args.Verbose, "v", false, "Prints verbose debug prints.")
	flag.StringVar(&args.Cmd, "c", "dump", "Command to execute on GitHub. Can be one of:\n"+
		"    - dump\n"+
		"        Dumps the data into csv files.\n"+
		"            -e = \"all\" for full dump or comma separated list of one or more of:\n"+
		"                   "+keysOfMap(entityMap)+"\n"+
		"            -b = Name of the table field to use for sorting the result of dump.\n"+
		"                 If empty, default sort on GitHub is creation date.\n"+
		"        Dump by time of creation is the default command.\n"+
		"    - archive\n"+
		"        Removes GitHub repositories according to:\n"+
		"            - n = least active n\n"+
		"            - s = inactive since <date>\n"+
		"            - r = comma separated list of repository names\n"+
		"        downloads the repository, creates a tarball and stores it in:\n"+
		"            - o = output folder\n"+
		"    - remove\n"+
		"        Removes GitHub user from the organizational account according to:\n"+
		"            - n = least active n\n"+
		"            - s = inactive since <date>\n"+
		"            - r = comma separated list of repository names\n")
	flag.StringVar(&args.Entity, "e", "all", "List of comma separated tables from the GitHub database"+
		" to apply the command on. Can be:\n"+
		"    - all = all the tables,\n"+
		"    - comma separated list of one or more of:\n"+
		"        "+keysOfMap(entityMap)+"\n")
	flag.StringVar(&args.By, "b", "", "Name of the table field to use for sorting the result of dump. "+
		"If empty, default sort on GitHub is creation date.\n")
	flag.StringVar(&args.Token, "t", gnet.Conf.Token, "Security token used on Github.\n"+
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
	flag.StringVar(&args.Organization, "o", gnet.Conf.Organization, "Organizational account being analyzed.")
	flag.Parse()

	// command line params trump config file
	gnet.Conf.Token = args.Token
	gnet.Conf.Organization = args.Organization

	// set debug/log
	utils.Debug.Verbose = args.Verbose
	log.SetOutput(os.Stdout)
}

func main() {
	if args.Help {
		flag.Usage()
		return
	}

	activeEntities, err := validateEntities()
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		return
	}

	var cmd commands.Command
	switch args.Cmd {
	case "dump":
		if args.By != "" {
			err = validateEntitySortField(activeEntities)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		cmd = &commands.Dump{By: args.By}
	case "remove":
		cmd = &commands.Remove{N: args.N, Since: args.Since, Names: args.Names}
	default:
		fmt.Println(fmt.Sprintf("Command `%s` not defined.\n", args.Cmd))
		flag.Usage()
		return
	}

	cmd.AddCache(commands.Cache(activeEntities, entityMap))
	cmd.Do(entityMap)
}

func validateEntities() ([]string, error) {
	var activeEntities = make([]string, 0, len(entityMap))
	if args.Entity == "all" {
		for entityName, _ := range entityMap {
			activeEntities = append(activeEntities, entityName)
		}
	} else {
		var slices = strings.Split(args.Entity, ",")
		for _, s := range slices {
			_, ok := entityMap[s]
			if !ok {
				return []string{}, errors.New(fmt.Sprintf("Unknown table: %s\n", s))
			}
			activeEntities = append(activeEntities, s)
		}
	}
	return activeEntities, nil
}

func validateEntitySortField(activeEntities []string) error {
	for _, entityName := range activeEntities {
		entity := entityMap[entityName]
		if !entity.HasField(args.By) {
			return errors.New(fmt.Sprintf("Field `%s` not found in `%s`. Choose one of: %s.\n",
				args.By,
				entityName,
				strings.Join(entity.GetTableFields(), ", ")))
		}
	}

	return nil
}

func keysOfMap(m map[string]entities.Entity) string {
	var keys = ""
	for key, _ := range m {
		keys = keys + key + ", "
	}
	return keys[:len(keys)-2]
}
