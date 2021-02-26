//
// Copyright (c) 2019 Sony Mobile Communications Inc.
// SPDX-License-Identifier: MIT
//

package model

import (
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

	String() string
}

var (
	// EntityMap contains a map of:
	//
	//   key = entity name
	//   value = struct implementing entity interface
	//
	// for all entities allowed to be used from interactive commands.
	EntityMap map[string]Entity
	Repos     *ReposResponse
	Users     *UsersResponse
	Teams     *TeamsResponse
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
}

// EntityNamesList returns the list of names of entities
// (use e.g. to dump the list)
func EntityNamesList() []string {
	return keysOf(EntityMap)
}

// ValidateEntities checks that a comma separated string of entities, e,
// is correct and returns the list of Entity objects (and nil error).
// In case of error returns empty Entity list and  "Unknown entity" error.
func ValidateEntities(e string) ([]Entity, error) {
	var activeEntities = make([]Entity, 0, len(EntityMap))
	if e == "" || e == "all" {
		for _, entity := range EntityMap {
			activeEntities = append(activeEntities, entity)
		}
	} else {
		var slices = strings.Split(e, ",")
		for _, s := range slices {
			entity, ok := EntityMap[s]
			if !ok {
				return []Entity{}, fmt.Errorf("Unknown entity: %s\n", s)
			}
			activeEntities = append(activeEntities, entity)
		}
	}

	return activeEntities, nil
}

// ValidateEntityField checks that a given field exists in the list of
// activeEntities.
func ValidateEntityField(field string, activeEntities []Entity) error {
	if field == "" || field == ID.Name {
		return nil
	}

	for _, entity := range activeEntities {
		if !entity.HasField(field) {
			return fmt.Errorf("Field `%s` not found in `%s`. Choose one of: %s.\n",
				field,
				entity.GetName(),
				strings.Join(entity.GetFields().DisplayNames(), ", "))
		}
	}

	return nil
}

func keysOf(m map[string]Entity) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
