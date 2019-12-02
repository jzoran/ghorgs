package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const UsersGraphQlJson = "config/users.json"
const UsersNextGraphQlJson = "config/users_next.json"
const UsersCsv = "users.csv"

type UserStatus struct {
	Message string `json:"message"`
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

func makeUsersQuery(count int) UsersQuery {
	return UsersQuery{makeQuery(UsersGraphQlJson, count)}
}

func (q *UsersQuery) getNext(after string) {
	q.Query.getNext(UsersNextGraphQlJson, after)
}

func (r *UsersResponse) fromJsonBuffer(buff []byte) {
	err := json.Unmarshal(buff, &r)
	if err != nil {
		panic(err)
	}
}

func (r *UsersResponse) toString() string {
	s, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(s)
}

func (r *UsersResponse) appendCsv() {
	var f *os.File

	if _, err := os.Stat(UsersCsv); err == nil {
		f, err = os.OpenFile(UsersCsv,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
	} else if os.IsNotExist(err) {
		f, err = os.OpenFile(UsersCsv,
			os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if _, err := f.WriteString("Id\tLogin\tName\tAdmin\t2FA\tEmail\tCompany\tUrl\tBio\tStatus\tUpdated\n"); err != nil {
			panic(err)
		}
	} else {
		panic(err)
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

		bio := ""
		if user.Member.Bio != nil {
			bio = *user.Member.Bio
		}

		msg := ""
		if user.Member.Status != nil {
			msg = user.Member.Status.Message
		}

		s := fmt.Sprintf("%s\t%s\t%s\t%s\t%t\t%s\t%s\t%s\t%s\t%s\t%s\n",
			user.Member.Id,
			user.Member.Login,
			name,
			user.Role,
			user.Has2FA,
			email,
			company,
			user.Member.Url,
			bio,
			msg,
			user.Member.UpdatedAt)
		if _, err := f.WriteString(s); err != nil {
			panic(err)
		}
	}
}
