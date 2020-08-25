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

type GitHubConfiguration struct {
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
	Conf       GitHubConfiguration
)

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

func MakeGitHubV4Request(query string, token string) *Request {
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
