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
	"log"
	"net/http"
)

func Cache(request []model.Entity) (map[string]*model.Table, error) {
	result := make(map[string]*model.Table, len(request))
	for _, entity := range request {
		fmt.Printf("\nCaching %s...", entity.GetName())

		t := entity.MakeTable()
		req := entity.MakeQuery(gnet.Conf.Organization)
		if utils.Debug.Verbose {
			log.Print(req)
		}

		gitHubRequest := gnet.MakeGitHubV4Request(req.GetGraphQlJson(), gnet.Conf.Token)
		resp, status := gitHubRequest.Execute()
		if status.Code != http.StatusOK {
			return result, fmt.Errorf("Error! HttpResponse: %s", status.Status)
		}

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

			gitHubRequest = gnet.MakeGitHubV4Request(req.GetGraphQlJson(), gnet.Conf.Token)
			resp, status = gitHubRequest.Execute()
			if status.Code != http.StatusOK {
				return result, fmt.Errorf("Error! HttpResponse: %s", status.Status)
			}

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
		result[entity.GetName()] = t
		if !utils.Debug.Verbose {
			fmt.Printf("\n")
		}
	}

	return result, nil
}
