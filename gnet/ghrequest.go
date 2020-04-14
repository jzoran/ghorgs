// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package gnet

import (
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

var Conf GitHubConfiguration

func MakeGitHubRequest(query string, token string) *Request {
	if Conf.Token == "" {
		panic("Missing GitHub token.")
	}

	return &Request{Conf.Url,
		map[string]string{"Authorization": "bearer " + Conf.Token},
		query,
		time.Duration(Conf.TimeOut) * time.Second}
}
