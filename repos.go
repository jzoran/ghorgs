package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const ReposGql = "config/repos.json"
const ReposNextGql = "config/repos_next.json"
const ReposCsv = "repos.csv"

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
	Count int
	Query string
}

func makeReposQuery(count int) ReposQuery {
	if count == 0 {
		count = githubConfig.PerPage
	}

	bytes, err := ioutil.ReadFile(ReposGql)
	if err != nil {
		panic(err)
	}
	return ReposQuery{count, fmt.Sprintf(string(bytes), count)}
}

func (r *ReposQuery) getNext(after string) {
	bytes, err := ioutil.ReadFile(ReposNextGql)
	if err != nil {
		panic(err)
	}

	r.Query = fmt.Sprintf(string(bytes), r.Count, after)
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

func (r *ReposResponse) appendCsv() {
	var f *os.File

	if _, err := os.Stat(ReposCsv); err == nil {
		f, err = os.OpenFile(ReposCsv,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
	} else if os.IsNotExist(err) {
		f, err = os.OpenFile(ReposCsv,
			os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if _, err := f.WriteString("Id\tName\tUrl\tDiskUsage (kB)\tUpdated\tLast Push\n"); err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}
	for _, repo := range r.Data.Org.Repos.Nodes {
		s := fmt.Sprintf("%s\t%s\t%s\t%d\t%s\t%s\n",
			repo.Id,
			repo.Name,
			repo.Url,
			repo.DiskUsage,
			repo.UpdatedAt,
			repo.PushedAt)
		if _, err := f.WriteString(s); err != nil {
			panic(err)
		}
	}
}
