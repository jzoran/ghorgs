//
// Copyright (c) 2019 Sony Mobile Communications Inc.
// SPDX-License-Identifier: MIT
//

package gnet

import (
	"net/url"
	"path"
	"time"
)

const (
	ConfigPath = "./config"
	ConfigName = "config"
	ConfigType = "yaml"
)

type gitHubConfiguration struct {
	Url          string `mapstructure:"url"`
	User         string `mapstructure:"user"`
	Token        string `mapstructure:"token"`
	Organization string `mapstructure:"organization"`
	PerPage      int    `mapstructure:"per_page"`
	TimeOut      int    `mapstructure:"time_out"`
}

var (
	v4Path     = "/graphql"
	postMethod = "POST"
	Conf       gitHubConfiguration
)

// MakeGitHubV3Request creates a Request object to access and
// execute GitHub REST like API v3. See https://docs.github.com/en/rest
// for reference (and in particular for the meaning of method and
// query parameters).
//
//    method = HTTP verb [HEAD, GET, POST, PATCH, PUT, DELETE]
//    query = specific resource path on GitHub API endpoint, e.g.
//            /user/repos and similar.
//    token = API authorization token
func MakeGitHubV3Request(method, query, token string) *Request {
	if Conf.Token == "" {
		panic("Missing GitHub token.")
	}

	u, err := url.Parse(Conf.Url)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, query)

	return &Request{u.String(),
		method,
		map[string]string{"Authorization": "bearer " + Conf.Token},
		query,
		time.Duration(Conf.TimeOut) * time.Second}
}

// MakeGitHubV4Request creates a Request object to access and
// execute GitHub GraphQL API v4. See https://docs.github.com/en/graphql
// for reference (and in particular about the meaning of the query
// parameter).
//
//     query = json representation of the graphql query
//     token = API authorization token
func MakeGitHubV4Request(query, token string) *Request {
	if Conf.Token == "" {
		panic("Missing GitHub token.")
	}

	u, err := url.Parse(Conf.Url)
	if err != nil {
		panic(err)
	}

	u.Path = path.Join(u.Path, v4Path)

	return &Request{u.String(),
		postMethod,
		map[string]string{"Authorization": "bearer " + Conf.Token},
		query,
		time.Duration(Conf.TimeOut) * time.Second}
}
