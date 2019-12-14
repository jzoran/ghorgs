// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"time"
)

const UsersGraphQlJson = "config/users.json"
const UsersNextGraphQlJson = "config/users_next.json"
const UsersCsv = "users.csv"
const UsersName = "users"

var UsersCsvTitle = []string{"Id", "Login", "Name", "Admin", "2FA", "Email", "Company", " Url", "Bio", "Status", "Updated", "Repositories Contributed To"}

type UserStatus struct {
	Message string `json:"message"`
}

type UserRepoNames struct {
	Name string `json:"nameWithOwner"`
}

type UserRepos struct {
	TotalCount int             `json:"totalCount"`
	ReposList  []UserRepoNames `json:"nodes"`
}

type User struct {
	Id        string      `json:"id"`
	Login     string      `json:"login"`
	Name      *string     `json:"name",omitempty`
	Email     *string     `json:"email",omitempty`
	Company   *string     `json:"company",omitempty`
	Url       string      `json:"url"`
	Bio       *string     `json:"bio",omitempty`
	Status    *UserStatus `json:"status",omitempty`
	UpdatedAt time.Time   `json:"updatedAt"`
	Repos     UserRepos   `json:"repositoriesContributedTo"`
}

type OrgMember struct {
	Has2FA bool   `json:"hasTwoFactorEnabled"`
	Role   string `json:"role"`
	Member User   `json:"node"`
}

type OrgMembers struct {
	Nodes    []OrgMember `json:"edges"`
	PageInfo Paging      `json:"pageInfo"`
	Total    int         `json:"totalCount"`
}

type UsersOrganization struct {
	Members OrgMembers `json:"membersWithRole"`
}

type UsersDataMap struct {
	Org UsersOrganization `json:"organization"`
}

type UsersResponse struct {
	Data UsersDataMap `json:"data"`
}

type UsersQuery struct {
	Query
}

func makeUsersQuery(organization string) *UsersQuery {
	return &UsersQuery{makeQuery(UsersGraphQlJson, organization)}
}

func (q *UsersQuery) getCount() int {
	return q.Query.Count
}

func (q *UsersQuery) getNext(after string) {
	q.Query.getNext(UsersNextGraphQlJson, after)
}

func (r *UsersResponse) getName() string {
	return UsersName
}

func (r *UsersResponse) makeCsv() *Csv {
	return makeCsv(UsersCsv)
}

func (r *UsersResponse) makeQuery(org string) IQuery {
	return makeUsersQuery(org)
}

func (r *UsersResponse) fromJsonBuffer(buff []byte) {
	err := json.Unmarshal(buff, &r)
	if err != nil {
		panic(err)
	}
}

func (r *UsersResponse) getTotal() int {
	return r.Data.Org.Members.Total
}

func (r *UsersResponse) hasNext() bool {
	return r.Data.Org.Members.PageInfo.HasNext
}

func (r *UsersResponse) getAfter() string {
	return r.Data.Org.Members.PageInfo.End
}

func (r *UsersResponse) toString() string {
	s, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(s)
}

func (r *UsersResponse) appendCsv(c *Csv) {
	if c.Records == nil {
		c.Records = make(map[string][]string)
		c.Keys = make([]string, 0)
	}

	for _, user := range r.Data.Org.Members.Nodes {
		name := "-"
		if user.Member.Name != nil {
			name = *user.Member.Name
		}
		email := "-"
		if user.Member.Email != nil {
			email = *user.Member.Email
		}
		company := "-"
		if user.Member.Company != nil {
			company = *user.Member.Company
		}
		bio := "-"
		if user.Member.Bio != nil {
			bio = *user.Member.Bio
		}
		msg := "-"
		if user.Member.Status != nil {
			msg = user.Member.Status.Message
		}
		repos := "-"
		if user.Member.Repos.TotalCount > 0 {
			for i, repo := range user.Member.Repos.ReposList {
				if i == 0 {
					repos = repo.Name
				} else {
					repos = repos + ", " + repo.Name
				}
			}
		}

		c.addKey(user.Member.Id)
		c.Records[user.Member.Id] = []string{user.Member.Login,
			name,
			user.Role,
			fmt.Sprintf("%t", user.Has2FA),
			email,
			company,
			user.Member.Url,
			bio,
			msg,
			fmt.Sprintf("%s", user.Member.UpdatedAt),
			repos}
	}
}

func (r *UsersResponse) getCsvTitle() []string {
	return UsersCsvTitle
}
