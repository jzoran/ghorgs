// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package cmd

import (
	"fmt"
	"ghorgs/cache"
	"ghorgs/entities"
	"ghorgs/gnet"
	"ghorgs/utils"
	"log"
)

func Cache(request []string, using map[string]entities.Entity) map[string]*cache.Table {
	result := make(map[string]*cache.Table, len(request))
	for _, entityName := range request {
		entity := using[entityName]
		fmt.Printf("\nCaching %s...", entityName)

		t := entity.MakeTable()
		req := entity.MakeQuery(gnet.Conf.Organization)
		if utils.Debug.Verbose {
			log.Print(req)
		}

		gitHubRequest := gnet.MakeGitHubRequest(req.GetGraphQlJson(), gnet.Conf.Token)
		resp := gitHubRequest.Fetch()
		if utils.Debug.Verbose {
			log.Print(resp)
		}

		entity.FromJsonBuffer(resp)
		entity.AppendTable(t)

		counter := req.GetCount()
		if utils.Debug.Verbose {
			log.Print(entity.ToString())
		} else {
			if counter <= entity.GetTotal() {
				fmt.Printf("\n%s: %d/%d", entity.GetName(), counter, entity.GetTotal())
			} else {
				fmt.Printf("\n%s: %d/%d", entity.GetName(), entity.GetTotal(), entity.GetTotal())
			}
		}

		for entity.HasNext() {
			req.GetNext(entity.GetNext())

			gitHubRequest = gnet.MakeGitHubRequest(req.GetGraphQlJson(), gnet.Conf.Token)
			resp = gitHubRequest.Fetch()

			entity.FromJsonBuffer(resp)
			entity.AppendTable(t)
			if utils.Debug.Verbose {
				log.Print(entity.ToString())
			} else {
				counter += req.GetCount()
				if counter <= entity.GetTotal() {
					fmt.Printf("\r%s: %d/%d", entity.GetName(), counter, entity.GetTotal())
				} else {
					fmt.Printf("\r%s: %d/%d", entity.GetName(), entity.GetTotal(), entity.GetTotal())
				}
			}
		}
		result[entityName] = t
	}

	return result
}
