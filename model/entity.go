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
	GetFields() Fields
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
	EntityMap       map[string]Entity
	EntityNamesList []string

	Repos *ReposResponse
	Users *UsersResponse
	Teams *TeamsResponse
)

func init() {
	Repos = &ReposResponse{}
	Users = &UsersResponse{}
	Teams = &TeamsResponse{}
	EntityMap = map[string]Entity{
		Repos.GetName(): Repos,
		Users.GetName(): Users,
		Teams.GetName(): Teams,
	}
	EntityNamesList = keysOf(EntityMap)
}

// Check that a comma separated list of entities, e,
// is correct and set ActiveEntities to requested
// subset of entities.
// In case of error return "Unknown entity" error,
// otherwise nil.
func ValidateEntities(e string) ([]Entity, error) {
	var activeEntities = make([]Entity, 0, len(EntityNamesList))
	if e == "" || e == "all" {
		for _, entity := range EntityMap {
			activeEntities = append(activeEntities, entity)
		}
	} else {
		var slices = strings.Split(e, ",")
		for _, s := range slices {
			entity, ok := EntityMap[s]
			if !ok {
				return []Entity{}, errors.New(fmt.Sprintf("Unknown entity: %s\n", s))
			}
			activeEntities = append(activeEntities, entity)
		}
	}

	return activeEntities, nil
}

// Check that a given field exists in the list
// of active entities ActiveEntities.
func ValidateEntityField(field string, activeEntities []Entity) error {
	if field == "" || field == ID.Name {
		return nil
	}

	for _, entity := range activeEntities {
		if !entity.HasField(field) {
			return errors.New(fmt.Sprintf("Field `%s` not found in `%s`. Choose one of: %s.\n",
				field,
				entity.GetName(),
				strings.Join(entity.GetFields().DisplayNames(), ", ")))
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
