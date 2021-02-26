//
// Copyright (c) 2019 Sony Mobile Communications Inc.
// SPDX-License-Identifier: MIT
//

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
	usersTableFields = &UsersFields{
		// &Field{"Id", -1}, // default for table as key for map of Records
		Login:        Field{"Login", 0},
		Name:         Field{"Name", 1},
		Admin:        Field{"Admin", 2},
		MFA:          Field{"2FA", 3},
		Email:        Field{"Email", 4},
		Company:      Field{"Company", 5},
		Url:          Field{"Url", 6},
		Updated:      Field{"Updated", 7},
		Repositories: Field{"Accessible Repositories", 8}}
	usersTableFieldNames = namesOf(usersTableFields.asList())
)

type UsersFields struct {
	Login        Field
	Name         Field
	Admin        Field
	MFA          Field
	Email        Field
	Company      Field
	Url          Field
	Updated      Field
	Repositories Field
}

func (f *UsersFields) asList() []Field {
	return []Field{usersTableFields.Login,
		usersTableFields.Name,
		usersTableFields.Admin,
		usersTableFields.MFA,
		usersTableFields.Email,
		usersTableFields.Company,
		usersTableFields.Url,
		usersTableFields.Updated,
		usersTableFields.Repositories}
}

func (f *UsersFields) DisplayNames() []string {
	return usersTableFieldNames
}

type UserStatus struct {
	Message string `json:"message"`
}

type UserRepoNames struct {
	Name string `json:"nameWithOwner"`
}

type UserRepos struct {
	TotalCount int `json:"totalCount"`
}

type User struct {
	Id        string    `json:"id"`
	Login     string    `json:"login"`
	Name      *string   `json:"name,omitempty"`
	Email     *string   `json:"email,omitempty"`
	Company   *string   `json:"company,omitempty"`
	Url       string    `json:"url"`
	UpdatedAt time.Time `json:"updatedAt"`
	Repos     UserRepos `json:"repositories"`
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
	return MakeTable(usersTableFields.asList())
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

func (r *UsersResponse) String() string {
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
		name := ""
		if user.Member.Name != nil {
			name = *user.Member.Name
		}
		email := ""
		if user.Member.Email != nil {
			email = *user.Member.Email
		}
		company := ""
		if user.Member.Company != nil {
			company = *user.Member.Company
		}
		repos := fmt.Sprintf("%d", user.Member.Repos.TotalCount)

		c.AddKey(user.Member.Id)
		c.Records[user.Member.Id] = []string{user.Member.Login,
			name,
			user.Role,
			fmt.Sprintf("%t", user.Has2FA),
			email,
			company,
			user.Member.Url,
			user.Member.UpdatedAt.String(),
			repos}
	}
}

func (r *UsersResponse) GetFields() Fields {
	return usersTableFields
}

func (r *UsersResponse) HasField(s string) bool {
	return utils.StringInSlice(s, usersTableFieldNames)
}

func (r *UsersResponse) GetCsvFile() string {
	return usersCsv
}
