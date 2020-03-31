// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package model

import (
	"errors"
	"fmt"
	"strings"
)

// Entity interface represents methods to create and
// execute query and create and execute a response to
// that query
type Entity interface {
	GetName() string

	MakeTable() *Table
	AppendTable(c *Table)
	GetTableFields() []Field
	GetTableFieldNames() []string
	HasField(s string) bool
	GetCsvFile() string

	MakeQuery(org string) Query
	FromJsonBuffer(buff []byte)
	GetTotal() int
	HasNext() bool
	GetNext() string

	ToString() string
}

var (
	EntityMap  map[string]Entity
	EntityList []string

	repos *ReposResponse
	users *UsersResponse
)

func init() {
	repos = &ReposResponse{}
	users = &UsersResponse{}
	EntityMap = map[string]Entity{
		repos.GetName(): repos,
		users.GetName(): users,
	}
	EntityList = keysOf(EntityMap)
}

// Check that a comma separated list of entities, e,
// is correct and set ActiveEntities to requested
// subset of entities.
// In case of error return "Unknown entity" error,
// otherwise nil.
func ValidateEntities(e string) ([]string, error) {
	var activeEntities = make([]string, 0, len(EntityList))
	if e == "" || e == "all" {
		for name, _ := range EntityMap {
			activeEntities = append(activeEntities, name)
		}
	} else {
		var slices = strings.Split(e, ",")
		for _, s := range slices {
			_, ok := EntityMap[s]
			if !ok {
				return []string{}, errors.New(fmt.Sprintf("Unknown entity: %s\n", s))
			}
			activeEntities = append(activeEntities, s)
		}
	}

	return activeEntities, nil
}

// Check that a given field exists in the list
// of active entities ActiveEntities.
func ValidateEntityField(field string, activeEntities []string) error {
	if field == "" || field == ID.Name {
		return nil
	}

	for _, entityName := range activeEntities {
		entity := EntityMap[entityName]
		if !entity.HasField(field) {
			return errors.New(fmt.Sprintf("Field `%s` not found in `%s`. Choose one of: %s.\n",
				field,
				entityName,
				strings.Join(entity.GetTableFieldNames(), ", ")))
		}
	}

	return nil
}

func keysOf(m map[string]Entity) []string {
	keys := make([]string, 0, len(m))
	for key, _ := range m {
		keys = append(keys, key)
	}
	return keys
}
