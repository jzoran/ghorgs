package main

import (
	"flag"
	"fmt"
	"log"
)

type Args struct {
	Help    bool
	Verbose bool
	Token   string
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
	flag.Parse()
	debug.Verbose = args.Verbose
}

func main() {
	if args.Help {
		flag.Usage()
		return
	}

	var req = makeReposQuery(0)
	if debug.Verbose {
		log.Print(req.GraphQlQueryJson)
	}

	gitHubRequest := makeGitHubRequest(req.GraphQlQueryJson, args.Token)
	resp := gitHubRequest.fetch()

	var dat ReposResponse
	dat.fromJsonBuffer(resp)
	dat.appendCsv()

	counter := req.Count
	if debug.Verbose {
		log.Print(dat.toString())
	} else {
		if counter <= dat.Data.Org.Repos.Total {
			fmt.Printf("repos: %d/%d", counter, dat.Data.Org.Repos.Total)
		} else {
			fmt.Printf("repos: %d/%d", dat.Data.Org.Repos.Total, dat.Data.Org.Repos.Total)
		}
	}

	for dat.Data.Org.Repos.PageInfo.HasNext {
		req.getNext(dat.Data.Org.Repos.PageInfo.End)
		if debug.Verbose {
			log.Print(req.GraphQlQueryJson)
		}
		gitHubRequest = makeGitHubRequest(req.GraphQlQueryJson, args.Token)
		resp = gitHubRequest.fetch()

		dat.fromJsonBuffer(resp)
		dat.appendCsv()
		if debug.Verbose {
			log.Print(dat.toString())
		} else {
			counter += req.Count
			if counter <= dat.Data.Org.Repos.Total {
				fmt.Printf("\rrepos: %d/%d", counter, dat.Data.Org.Repos.Total)
			} else {
				fmt.Printf("\rrepos: %d/%d", dat.Data.Org.Repos.Total, dat.Data.Org.Repos.Total)
			}
		}
	}

	var ureq = makeUsersQuery(0)
	if debug.Verbose {
		log.Print(ureq.GraphQlQueryJson)
	}

	gitHubRequest = makeGitHubRequest(ureq.GraphQlQueryJson, args.Token)
	resp = gitHubRequest.fetch()

	var udat UsersResponse
	udat.fromJsonBuffer(resp)
	udat.appendCsv()

	counter = ureq.Count
	if debug.Verbose {
		log.Print(udat.toString())
	} else {
		if counter <= udat.Data.Org.Members.Total {
			fmt.Printf("\nusers: %d/%d", counter, udat.Data.Org.Members.Total)
		} else {
			fmt.Printf("\nusers: %d/%d", udat.Data.Org.Members.Total, udat.Data.Org.Members.Total)
		}
	}

	for udat.Data.Org.Members.PageInfo.HasNext {
		ureq.getNext(udat.Data.Org.Members.PageInfo.End)
		if debug.Verbose {
			log.Print(ureq.GraphQlQueryJson)
		}
		gitHubRequest = makeGitHubRequest(ureq.GraphQlQueryJson, args.Token)
		resp = gitHubRequest.fetch()

		udat.fromJsonBuffer(resp)
		udat.appendCsv()
		if debug.Verbose {
			log.Print(udat.toString())
		} else {
			counter += ureq.Count
			if counter <= udat.Data.Org.Members.Total {
				fmt.Printf("\rusers: %d/%d", counter, udat.Data.Org.Members.Total)
			} else {
				fmt.Printf("\rusers: %d/%d", udat.Data.Org.Members.Total, udat.Data.Org.Members.Total)
			}
		}
	}
}
