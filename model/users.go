// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package model

import (
	"encoding/json"
	"fmt"
	"ghorgs/utils"
	"time"
)

const (
	usersGraphQlJson     = "config/users.json"
	usersNextGraphQlJson = "config/users_next.json"
	usersCsv             = "users.csv"
	usersName            = "users"
)

var (
	usersTableFields = []Field{
		//&Field{"Id", -1}, // default for table as key for map of Records
		Field{"Login", 0},
		Field{"Name", 1},
		Field{"Admin", 2},
		Field{"2FA", 3},
		Field{"Email", 4},
		Field{"Company", 5},
		Field{"Url", 6},
		Field{"Bio", 7},
		Field{"Status", 8},
		Field{"Updated", 9},
		Field{"Repositories Contributed To", 10}}
	usersTableFieldNames = namesOf(usersTableFields)
)

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
	QueryBase
}

func makeUsersQuery(organization string) *UsersQuery {
	return &UsersQuery{makeQuery(usersGraphQlJson, organization)}
}

func (q *UsersQuery) GetCount() int {
	return q.QueryBase.Count
}

func (q *UsersQuery) GetNext(after string) {
	q.QueryBase.getNext(usersNextGraphQlJson, after)
}

func (r *UsersResponse) GetName() string {
	return usersName
}

func (r *UsersResponse) MakeTable() *Table {
	return MakeTable(usersTableFields)
}

func (r *UsersResponse) MakeQuery(org string) Query {
	return makeUsersQuery(org)
}

func (r *UsersResponse) FromJsonBuffer(buff []byte) {
	err := json.Unmarshal(buff, &r)
	if err != nil {
		panic(err)
	}
}

func (r *UsersResponse) GetTotal() int {
	return r.Data.Org.Members.Total
}

func (r *UsersResponse) HasNext() bool {
	return r.Data.Org.Members.PageInfo.HasNext
}

func (r *UsersResponse) GetNext() string {
	return r.Data.Org.Members.PageInfo.End
}

func (r *UsersResponse) ToString() string {
	s, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(s)
}

func (r *UsersResponse) AppendTable(c *Table) {
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

		c.AddKey(user.Member.Id)
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

func (r *UsersResponse) GetTableFields() []Field {
	return usersTableFields
}

func (r *UsersResponse) GetTableFieldNames() []string {
	return usersTableFieldNames
}

func (r *UsersResponse) HasField(s string) bool {
	return utils.StringInSlice(s, usersTableFieldNames)
}

func (r *UsersResponse) GetCsvFile() string {
	return usersCsv
}
