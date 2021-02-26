//
// Copyright (c) 2019 Sony Mobile Communications Inc.
// SPDX-License-Identifier: MIT
//

package model

import (
	"encoding/json"
	"fmt"
	"ghorgs/utils"
)

const (
	teamsGraphQlJson     = "config/teams.json"
	teamsNextGraphQlJson = "config/teams_next.json"
	teamsCsv             = "teams.csv"
	teamsName            = "teams"
)

var (
	teamsTableFields = &TeamsFields{
		// &Field{"Id", -1}, // default for table as key for map of Records
		Name:         Field{"Name", 0},
		Url:          Field{"Url", 1},
		ParentId:     Field{"ParentId", 2},
		ParentName:   Field{"ParentName", 3},
		Children:     Field{"Children", 4},
		Repositories: Field{"Repositories", 5},
		Members:      Field{"Members", 6},
		Invitations:  Field{"Invitations", 7}}
	teamsTableFieldNames = namesOf(teamsTableFields.asList())
)

type TeamsFields struct {
	Name         Field
	Url          Field
	ParentId     Field
	ParentName   Field
	Children     Field
	Repositories Field
	Members      Field
	Invitations  Field
}

func (f *TeamsFields) asList() []Field {
	return []Field{teamsTableFields.Name,
		teamsTableFields.Url,
		teamsTableFields.ParentId,
		teamsTableFields.ParentName,
		teamsTableFields.Children,
		teamsTableFields.Repositories,
		teamsTableFields.Members,
		teamsTableFields.Invitations}
}

func (f *TeamsFields) DisplayNames() []string {
	return teamsTableFieldNames
}

type Team struct {
	Id          string       `json:"id"`
	Name        *string      `json:"name,omitempty"`
	Url         string       `json:"url"`
	Parent      *TeamParent  `json:"parentTeam,omitempty"`
	Children    TeamChildren `json:"childTeams"`
	Repos       TeamRepos    `json:"repositories"`
	Members     TeamMembers  `json:"members"`
	Invitations TeamInvites  `json:"invitations"`
}

type TeamParent struct {
	Id   string  `json:"id"`
	Name *string `json:"name,omitempty"`
}

type TeamChildren struct {
	TotalCount int `json:"totalCount"`
}

type TeamRepos struct {
	TotalCount int `json:"totalCount"`
}

type TeamMembers struct {
	TotalCount int `json:"totalCount"`
}

type TeamInvites struct {
	TotalCount int `json:"totalCount"`
}

type OrgTeams struct {
	Nodes    []Team `json:"nodes"`
	PageInfo Paging `json:"pageInfo"`
	Total    int    `json:"totalCount"`
}

type TeamsOrganization struct {
	Teams OrgTeams `json:"teams"`
}

type TeamsDataMap struct {
	Org TeamsOrganization `json:"organization"`
}

type TeamsResponse struct {
	Data TeamsDataMap `json:"data"`
}

type TeamsQuery struct {
	QueryBase
}

func makeTeamsQuery(organization string) *TeamsQuery {
	return &TeamsQuery{makeQuery(teamsGraphQlJson, organization)}
}

func (q *TeamsQuery) GetCount() int {
	return q.QueryBase.Count
}

func (q *TeamsQuery) GetNext(after string) {
	q.QueryBase.getNext(teamsNextGraphQlJson, after)
}

func (r *TeamsResponse) GetName() string {
	return teamsName
}

func (r *TeamsResponse) MakeTable() *Table {
	return MakeTable(teamsTableFields.asList())
}

func (r *TeamsResponse) MakeQuery(org string) Query {
	return makeTeamsQuery(org)
}

func (r *TeamsResponse) FromJsonBuffer(buff []byte) {
	err := json.Unmarshal(buff, &r)
	if err != nil {
		panic(err)
	}
}

func (r *TeamsResponse) GetTotal() int {
	return r.Data.Org.Teams.Total
}

func (r *TeamsResponse) HasNext() bool {
	return r.Data.Org.Teams.PageInfo.HasNext
}

func (r *TeamsResponse) GetNext() string {
	return r.Data.Org.Teams.PageInfo.End
}

func (r *TeamsResponse) String() string {
	s, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(s)
}

func (r *TeamsResponse) AppendTable(c *Table) {
	if c.Records == nil {
		c.Records = make(map[string][]string)
		c.Keys = make([]string, 0)
	}

	for _, team := range r.Data.Org.Teams.Nodes {
		name := ""
		if team.Name != nil {
			name = *team.Name
		}

		parentId := ""
		parentName := ""
		if team.Parent != nil {
			parentId = team.Parent.Id
			if team.Parent.Name != nil {
				parentName = *team.Parent.Name
			}
		}

		children := fmt.Sprintf("%d", team.Children.TotalCount)
		repos := fmt.Sprintf("%d", team.Repos.TotalCount)
		members := fmt.Sprintf("%d", team.Members.TotalCount)
		invites := fmt.Sprintf("%d", team.Invitations.TotalCount)

		c.AddKey(team.Id)
		c.Records[team.Id] = []string{name,
			team.Url,
			parentId,
			parentName,
			children,
			repos,
			members,
			invites}
	}
}

func (r *TeamsResponse) GetFields() Fields {
	return teamsTableFields
}

func (r *TeamsResponse) HasField(s string) bool {
	return utils.StringInSlice(s, teamsTableFieldNames)
}

func (r *TeamsResponse) GetCsvFile() string {
	return teamsCsv
}
