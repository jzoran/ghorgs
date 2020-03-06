// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package main

import (
	"errors"
	"flag"
	"fmt"
	"ghorgs/commands"
	"ghorgs/protocols"
	"ghorgs/utils"
	"log"
	"os"
	"strings"
)

type Args struct {
	Help         bool
	Verbose      bool
	Cmd          string
	Proj         string
	By           string
	N            int
	Since        string
	Names        string
	Token        string
	Organization string
}

var args Args
var protoMap map[string]protocols.Protocol

func init() {
	repos := &protocols.ReposResponse{}
	users := &protocols.UsersResponse{}
	protoMap = map[string]protocols.Protocol{
		repos.GetName(): repos,
		users.GetName(): users,
	}

	flag.BoolVar(&args.Help, "h", false, "Prints this help.")
	flag.BoolVar(&args.Verbose, "v", false, "Prints verbose debug prints.")
	flag.StringVar(&args.Cmd, "c", "dump", "Command to execute on GitHub. Can be one of:\n"+
		"    - dump\n"+
		"        Dumps the data into csv files.\n"+
		"            -d = \"all\" for full dump or comma separated list of one or more of:\n"+
		"                   "+keysOfMap(protoMap)+"\n"+
		"            -b = dump sorted by this column (default is Id).\n"+
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
	flag.StringVar(&args.Proj, "d", "all", "List of comma separated tables from the database"+
		" to apply the command on. Can be:\n"+
		"    - all = all the tables,\n"+
		"    - comma separated list of one or more of:\n"+
		"        "+keysOfMap(protoMap)+"\n")
	flag.StringVar(&args.By, "b", "Id", "Name of the column to use for sorting the result of dump.\n")
	flag.StringVar(&args.Token, "t", "", "Security token used on Github.\n"+
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
	flag.StringVar(&args.Organization, "o", "", "Organizational account being analyzed.")
	flag.Parse()
	utils.Debug.Verbose = args.Verbose
	log.SetOutput(os.Stdout)
}

func main() {
	if args.Help {
		flag.Usage()
		return
	}

	activeProtos, err := validateProtocols()
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		return
	}

	var cmd commands.Command
	switch args.Cmd {
	case "dump":
		err = validateProtocolSortBy(activeProtos)
		if err != nil {
			fmt.Println(err)
			return
		}
		cmd = &commands.Dump{By: args.By}
	case "remove":
		cmd = &commands.Remove{N: args.N, Since: args.Since, Names: args.Names}
	default:
		fmt.Println(fmt.Sprintf("Command `%s` not defined.\n", args.Cmd))
		flag.Usage()
		return
	}

	cmd.AddCache(commands.Cache(activeProtos, protoMap))
	cmd.Do(protoMap)
}

func validateProtocols() ([]string, error) {
	var activeProtos = make([]string, 0, len(protoMap))
	if args.Proj == "all" {
		for protoName, _ := range protoMap {
			activeProtos = append(activeProtos, protoName)
		}
	} else {
		var slices = strings.Split(args.Proj, ",")
		for _, s := range slices {
			_, ok := protoMap[s]
			if !ok {
				return []string{}, errors.New(fmt.Sprintf("Unknown table: %s\n", s))
			}
			activeProtos = append(activeProtos, s)
		}
	}
	return activeProtos, nil
}

func validateProtocolSortBy(activeProtos []string) error {
	res := true
	for _, protoName := range activeProtos {
		proto := protoMap[protoName]
		res = res && utils.StringInSlice(args.By, proto.GetCsvTitle())
	}
	if res {
		return nil
	}

	return errors.New("Sort by column not found.\n")
}

func keysOfMap(m map[string]protocols.Protocol) string {
	var keys = ""
	for key, _ := range m {
		keys = keys + key + " ,"
	}
	return keys[:len(keys)-2]
}
