// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Args struct {
	Help         bool
	Verbose      bool
	Token        string
	Organization string
}

var args Args

func init() {
	flag.BoolVar(&args.Help, "h", false, "Prints this help.")
	flag.BoolVar(&args.Verbose, "v", false, "Prints verbose debug prints.")
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
	debug.Verbose = args.Verbose
	log.SetOutput(os.Stdout)
}

func main() {
	if args.Help {
		flag.Usage()
		return
	}

	reposCsv := makeCsv(ReposCsv)
	req := makeReposQuery(args.Organization)
	if debug.Verbose {
		log.Print(req.GraphQlQueryJson)
	}

	gitHubRequest := makeGitHubRequest(req.GraphQlQueryJson, args.Token)
	resp := gitHubRequest.fetch()

	var repos ReposResponse
	repos.fromJsonBuffer(resp)
	repos.appendCsv(reposCsv)

	counter := req.Count
	if debug.Verbose {
		log.Print(repos.toString())
	} else {
		if counter <= repos.Data.Org.Repos.Total {
			fmt.Printf("repos: %d/%d", counter, repos.Data.Org.Repos.Total)
		} else {
			fmt.Printf("repos: %d/%d", repos.Data.Org.Repos.Total, repos.Data.Org.Repos.Total)
		}
	}

	for repos.Data.Org.Repos.PageInfo.HasNext {
		req.getNext(repos.Data.Org.Repos.PageInfo.End)
		if debug.Verbose {
			log.Print(req.GraphQlQueryJson)
		}
		gitHubRequest = makeGitHubRequest(req.GraphQlQueryJson, args.Token)
		resp = gitHubRequest.fetch()

		repos.fromJsonBuffer(resp)
		repos.appendCsv(reposCsv)
		if debug.Verbose {
			log.Print(repos.toString())
		} else {
			counter += req.Count
			if counter <= repos.Data.Org.Repos.Total {
				fmt.Printf("\rrepos: %d/%d", counter, repos.Data.Org.Repos.Total)
			} else {
				fmt.Printf("\rrepos: %d/%d", repos.Data.Org.Repos.Total, repos.Data.Org.Repos.Total)
			}
		}
	}
	reposCsv.flush(ReposCsvTitle)

	usersCsv := makeCsv(UsersCsv)
	ureq := makeUsersQuery(args.Organization)
	if debug.Verbose {
		log.Print(ureq.GraphQlQueryJson)
	}

	gitHubRequest = makeGitHubRequest(ureq.GraphQlQueryJson, args.Token)
	resp = gitHubRequest.fetch()

	var users UsersResponse
	users.fromJsonBuffer(resp)
	users.appendCsv(usersCsv)

	counter = ureq.Count
	if debug.Verbose {
		log.Print(users.toString())
	} else {
		if counter <= users.Data.Org.Members.Total {
			fmt.Printf("\nusers: %d/%d", counter, users.Data.Org.Members.Total)
		} else {
			fmt.Printf("\nusers: %d/%d", users.Data.Org.Members.Total, users.Data.Org.Members.Total)
		}
	}

	for users.Data.Org.Members.PageInfo.HasNext {
		ureq.getNext(users.Data.Org.Members.PageInfo.End)
		if debug.Verbose {
			log.Print(ureq.GraphQlQueryJson)
		}
		gitHubRequest = makeGitHubRequest(ureq.GraphQlQueryJson, args.Token)
		resp = gitHubRequest.fetch()

		users.fromJsonBuffer(resp)
		users.appendCsv(usersCsv)
		if debug.Verbose {
			log.Print(users.toString())
		} else {
			counter += ureq.Count
			if counter <= users.Data.Org.Members.Total {
				fmt.Printf("\rusers: %d/%d", counter, users.Data.Org.Members.Total)
			} else {
				fmt.Printf("\rusers: %d/%d", users.Data.Org.Members.Total, users.Data.Org.Members.Total)
			}
		}
		usersCsv.flush(UsersCsvTitle)
	}
}
