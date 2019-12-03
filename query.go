package main

import (
	"fmt"
	"io/ioutil"
)

type Query struct {
	Organization     string
	Count            int
	GraphQlQueryJson string
}

func makeQuery(jsonFile string, organization string) Query {
	if organization == "" {
		organization = githubConfig.Organization
	}

	if githubConfig.Organization == "" {
		panic("Missing GitHub Organization.")
	}

	bytes, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}
	return Query{organization, githubConfig.PerPage, fmt.Sprintf(string(bytes), organization, githubConfig.PerPage)}
}

func (q *Query) getNext(jsonFile string, after string) {
	bytes, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}

	q.GraphQlQueryJson = fmt.Sprintf(string(bytes), q.Organization, q.Count, after)
}
