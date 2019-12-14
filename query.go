// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

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

type IQuery interface {
	getGraphQlJson() string
	getNext(after string)
	getCount() int
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
	return Query{organization,
		githubConfig.PerPage,
		fmt.Sprintf(string(bytes), organization, githubConfig.PerPage)}
}

func (q *Query) getCount() int {
	return q.Count
}

func (q *Query) getGraphQlJson() string {
	return q.GraphQlQueryJson
}

func (q *Query) getNext(jsonFile string, after string) {
	bytes, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}

	q.GraphQlQueryJson = fmt.Sprintf(string(bytes), q.Organization, q.Count, after)
}
