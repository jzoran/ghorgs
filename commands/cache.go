// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package commands

import (
	"fmt"
	"ghorgs/cache"
	"ghorgs/gnet"
	"ghorgs/protocols"
	"ghorgs/utils"
	"log"
)

func Cache(request []string, using map[string]protocols.Protocol) map[string]*cache.Table {
	result := make(map[string]*cache.Table, len(request))
	for _, protoName := range request {
		proto := using[protoName]
		fmt.Printf("Caching %s...", protoName)

		t := proto.MakeTable()
		req := proto.MakeQuery(gnet.Conf.Organization)
		gitHubRequest := gnet.MakeGitHubRequest(req.GetGraphQlJson(), gnet.Conf.Token)
		resp := gitHubRequest.Fetch()

		proto.FromJsonBuffer(resp)
		proto.AppendTable(t)

		counter := req.GetCount()
		if utils.Debug.Verbose {
			log.Print(proto.ToString())
		} else {
			if counter <= proto.GetTotal() {
				fmt.Printf("\n%s: %d/%d", proto.GetName(), counter, proto.GetTotal())
			} else {
				fmt.Printf("\n%s: %d/%d", proto.GetName(), proto.GetTotal(), proto.GetTotal())
			}
		}

		for proto.HasNext() {
			req.GetNext(proto.GetNext())

			gitHubRequest = gnet.MakeGitHubRequest(req.GetGraphQlJson(), gnet.Conf.Token)
			resp = gitHubRequest.Fetch()

			proto.FromJsonBuffer(resp)
			proto.AppendTable(t)
			if utils.Debug.Verbose {
				log.Print(proto.ToString())
			} else {
				counter += req.GetCount()
				if counter <= proto.GetTotal() {
					fmt.Printf("\r%s: %d/%d", proto.GetName(), counter, proto.GetTotal())
				} else {
					fmt.Printf("\r%s: %d/%d", proto.GetName(), proto.GetTotal(), proto.GetTotal())
				}
			}
		}
		result[protoName] = t
	}

	return result
}
