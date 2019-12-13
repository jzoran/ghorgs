// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"time"
)

const ReposGraphQlJson = "config/repos.json"
const ReposNextGraphQlJson = "config/repos_next.json"
const ReposCsv = "repos.csv"

var ReposCsvTitle = []string{"Id", "Name", "Url", "DiskUsage (kB)", "Updated", "Last Push"}

type Repository struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	DiskUsage int       `json:"diskUsage"`
	UpdatedAt time.Time `json:"updatedAt"`
	PushedAt  time.Time `json:"pushedAt"`
}

type Paging struct {
	HasNext bool   `json:"hasNextPage"`
	End     string `json:"endCursor"`
}

type Repositories struct {
	Nodes    []Repository `json:"nodes"`
	PageInfo Paging       `json:"pageInfo"`
	Total    int          `json:"totalCount"`
}

type Organization struct {
	Repos Repositories `json:"repositories"`
}

type DataMap struct {
	Org Organization `json:"organization"`
}

type ReposResponse struct {
	Data DataMap `json:"data"`
}

type ReposQuery struct {
	Query
}

func makeReposQuery(organization string) ReposQuery {
	return ReposQuery{makeQuery(ReposGraphQlJson, organization)}
}

func (q *ReposQuery) getNext(after string) {
	q.Query.getNext(ReposNextGraphQlJson, after)
}

func (r *ReposResponse) fromJsonBuffer(buff []byte) {
	err := json.Unmarshal(buff, &r)
	if err != nil {
		panic(err)
	}
}

func (r *ReposResponse) toString() string {
	s, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(s)
}

func (r *ReposResponse) appendCsv(c Csv) {
	if c.Records == nil {
		c.Records = make(map[string][]string)
		c.Keys = make([]string, 0)
	}

	for _, repo := range r.Data.Org.Repos.Nodes {
		c.addKey(repo.Id)
		c.Records[repo.Id] = []string{repo.Name,
			repo.Url,
			fmt.Sprintf("%d", repo.DiskUsage),
			fmt.Sprintf("%s", repo.UpdatedAt),
			fmt.Sprintf("%s", repo.PushedAt)}
	}
}
