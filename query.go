package main

import (
	"fmt"
	"io/ioutil"
)

type Query struct {
	Count            int
	GraphQlQueryJson string
}

func makeQuery(jsonFile string, count int) Query {
	if count == 0 {
		count = githubConfig.PerPage
	}

	bytes, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}
	return Query{count, fmt.Sprintf(string(bytes), count)}
}

func (q *Query) getNext(jsonFile string, after string) {
	bytes, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}

	q.GraphQlQueryJson = fmt.Sprintf(string(bytes), q.Count, after)
}
