package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"time"
)

type GitHubConfig struct {
	Url          string `yaml:"url"`
	Token        string `yaml:"token"`
	Organization string `yaml:"organization"`
	PerPage      int    `yaml:"per_page"`
	TimeOut      int    `yaml:"time_out"`
}

var githubConfig GitHubConfig

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

	yaml.Unmarshal(values, &githubConfig)
}

func makeGitHubRequest(query string, token string) *Request {
	if token != "" {
		githubConfig.Token = token
	}

	if githubConfig.Token == "" {
		panic("Missing GitHub token.")
	}

	return &Request{githubConfig.Url,
		map[string]string{"Authorization": "bearer " + githubConfig.Token},
		query,
		time.Duration(githubConfig.TimeOut) * time.Second}
}
