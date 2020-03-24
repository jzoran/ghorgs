// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package entities

import (
	"encoding/json"
	"fmt"
	"ghorgs/cache"
	"ghorgs/fields"
	"ghorgs/utils"
	"time"
)

const (
	reposGraphQlJson     = "config/repos.json"
	reposNextGraphQlJson = "config/repos_next.json"
	reposCsv             = "repos.csv"
	reposName            = "repos"
)

var (
	reposTableFields = []fields.Field{
		//&fields.Field{"Id", -1} // default for table as key for map of Records
		fields.Field{"Name", 0},
		fields.Field{"Type", 1},
		fields.Field{"Url", 2},
		fields.Field{"DiskUsage (kB)", 3},
		fields.Field{"Updated", 4},
		fields.Field{"Last Push", 5}}
	reposTableFieldNames = fields.NamesOf(reposTableFields)
)

type Repository struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Private   bool      `json:"isPrivate"`
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

func makeReposQuery(organization string) *ReposQuery {
	return &ReposQuery{makeQuery(reposGraphQlJson, organization)}
}

func (q *ReposQuery) GetCount() int {
	return q.Query.Count
}

func (q *ReposQuery) GetNext(after string) {
	q.Query.getNext(reposNextGraphQlJson, after)
}

func (r *ReposResponse) GetName() string {
	return reposName
}

func (r *ReposResponse) MakeTable() *cache.Table {
	return cache.MakeTable(reposTableFields)
}

func (r *ReposResponse) MakeQuery(org string) IQuery {
	return makeReposQuery(org)
}

func (r *ReposResponse) FromJsonBuffer(buff []byte) {
	err := json.Unmarshal(buff, &r)
	if err != nil {
		panic(err)
	}
}

func (r *ReposResponse) GetTotal() int {
	return r.Data.Org.Repos.Total
}

func (r *ReposResponse) HasNext() bool {
	return r.Data.Org.Repos.PageInfo.HasNext
}

func (r *ReposResponse) GetNext() string {
	return r.Data.Org.Repos.PageInfo.End
}

func (r *ReposResponse) ToString() string {
	s, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(s)
}

func (r *ReposResponse) AppendTable(c *cache.Table) {
	if c.Records == nil {
		c.Records = make(map[string][]string)
	}

	for _, repo := range r.Data.Org.Repos.Nodes {
		c.AddKey(repo.Id)
		isPrivate := "PUBLIC"
		if repo.Private {
			isPrivate = "PRIVATE"
		}
		c.Records[repo.Id] = []string{repo.Name,
			isPrivate,
			repo.Url,
			fmt.Sprintf("%d", repo.DiskUsage),
			fmt.Sprintf("%s", repo.UpdatedAt),
			fmt.Sprintf("%s", repo.PushedAt)}
	}
}

func (r *ReposResponse) GetTableFields() []fields.Field {
	return reposTableFields
}

func (r *ReposResponse) GetTableFieldNames() []string {
	return reposTableFieldNames
}

func (r *ReposResponse) HasField(s string) bool {
	return utils.StringInSlice(s, reposTableFieldNames)
}

func (r *ReposResponse) GetCsvFile() string {
	return reposCsv
}
