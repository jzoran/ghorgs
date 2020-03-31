// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package model

import (
	"fmt"
	"ghorgs/gnet"
	"io/ioutil"
)

type QueryBase struct {
	Organization     string
	Count            int
	GraphQlQueryJson string
}

type Query interface {
	GetGraphQlJson() string
	GetNext(after string)
	GetCount() int
}

func makeQuery(jsonFile string, organization string) QueryBase {
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
	return QueryBase{organization,
		gnet.Conf.PerPage,
		fmt.Sprintf(string(bytes), organization, gnet.Conf.PerPage)}
}

func (q *QueryBase) GetCount() int {
	return q.Count
}

func (q *QueryBase) GetGraphQlJson() string {
	return q.GraphQlQueryJson
}

func (q *QueryBase) getNext(jsonFile string, after string) {
	bytes, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}

	q.GraphQlQueryJson = fmt.Sprintf(string(bytes), q.Organization, q.Count, after)
}
