// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package gnet

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"time"
)

type GitHubConfiguration struct {
	Url          string `yaml:"url"`
	Token        string `yaml:"token"`
	Organization string `yaml:"organization"`
	PerPage      int    `yaml:"per_page"`
	TimeOut      int    `yaml:"time_out"`
}

var Conf GitHubConfiguration

func init() {
	f, err := os.Open("config/config.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	values, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	yaml.Unmarshal(values, &Conf)
}

func MakeGitHubRequest(query string, token string) *Request {
	if token != "" {
		Conf.Token = token
	}

	if Conf.Token == "" {
		panic("Missing GitHub token.")
	}

	return &Request{Conf.Url,
		map[string]string{"Authorization": "bearer " + Conf.Token},
		query,
		time.Duration(Conf.TimeOut) * time.Second}
}
