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
var protocols []Protocol

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

	protocols = []Protocol{&ReposResponse{}, &UsersResponse{}}
}

func main() {
	if args.Help {
		flag.Usage()
		return
	}

	for _, proto := range protocols {
		csv := proto.makeCsv()
		req := proto.makeQuery(args.Organization)
		gitHubRequest := makeGitHubRequest(req.getGraphQlJson(), args.Token)
		resp := gitHubRequest.fetch()

		proto.fromJsonBuffer(resp)
		proto.appendCsv(csv)

		counter := req.getCount()
		if debug.Verbose {
			log.Print(proto.toString())
		} else {
			if counter <= proto.getTotal() {
				fmt.Printf("\n%s: %d/%d", proto.getName(), counter, proto.getTotal())
			} else {
				fmt.Printf("\n%s: %d/%d", proto.getName(), proto.getTotal(), proto.getTotal())
			}
		}

		for proto.hasNext() {
			req.getNext(proto.getAfter())

			gitHubRequest = makeGitHubRequest(req.getGraphQlJson(), args.Token)
			resp = gitHubRequest.fetch()

			proto.fromJsonBuffer(resp)
			proto.appendCsv(csv)
			if debug.Verbose {
				log.Print(proto.toString())
			} else {
				counter += req.getCount()
				if counter <= proto.getTotal() {
					fmt.Printf("\r%s: %d/%d", proto.getName(), counter, proto.getTotal())
				} else {
					fmt.Printf("\r%s: %d/%d", proto.getName(), proto.getTotal(), proto.getTotal())
				}
			}
		}
		csv.flush(proto.getCsvTitle())
	}
}
