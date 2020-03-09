// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package entities

import (
	"fmt"
	"ghorgs/gnet"
	"io/ioutil"
)

type Query struct {
	Organization     string
	Count            int
	GraphQlQueryJson string
}

type IQuery interface {
	GetGraphQlJson() string
	GetNext(after string)
	GetCount() int
}

func makeQuery(jsonFile string, organization string) Query {
	if organization == "" {
		organization = gnet.Conf.Organization
	}

	if gnet.Conf.Organization == "" {
		panic("Missing GitHub Organization.")
	}

	bytes, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}
	return Query{organization,
		gnet.Conf.PerPage,
		fmt.Sprintf(string(bytes), organization, gnet.Conf.PerPage)}
}

func (q *Query) GetCount() int {
	return q.Count
}

func (q *Query) GetGraphQlJson() string {
	return q.GraphQlQueryJson
}

func (q *Query) getNext(jsonFile string, after string) {
	bytes, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}

	q.GraphQlQueryJson = fmt.Sprintf(string(bytes), q.Organization, q.Count, after)
}
