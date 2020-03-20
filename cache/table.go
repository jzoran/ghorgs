// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package cache

import (
	"errors"
	"fmt"
	"log"
	"sort"
)

type Table struct {
	Records   map[string][]string
	Keys      []string
	Fields    []string
	sortField int
}

func MakeTable(fields []string) *Table {
	keys := make([]string, 0)
	return &Table{Records: nil, Keys: keys, Fields: fields}
}

func (t *Table) AddKey(key string) {
	t.Keys = append(t.Keys, key)
}

func (t *Table) ToString() string {
	var s string
	for id, key := range t.Keys {
		s += key + "\t"
		line := t.Records[key]
		for i, cell := range line {
			s += cell
			if i < len(line)-1 {
				s += "\t"
			}
		}
		if id < len(t.Keys)-1 {
			s += "\n"
		}
	}

	return s
}

func (t *Table) Log() {
	log.Print(t.ToString() + "\n")
}

// sort interface + method
type By Table

func (a By) Len() int { return len(a.Keys) }
func (a By) Less(i, j int) bool {
	if a.sortField == -1 {
		return a.Keys[i] < a.Keys[j]
	}
	return a.Records[a.Keys[i]][a.sortField] < a.Records[a.Keys[j]][a.sortField]
}
func (a By) Swap(i, j int) {
	a.Keys[i], a.Keys[j] = a.Keys[j], a.Keys[i]
}

func (t *Table) SortByField(field string) (*Table, error) {
	t.sortField = -1
	for i, val := range t.Fields {
		if val == field {
			t.sortField = i - 1 // Fields include "Id" at index 0
			break
		}
	}
	if t.sortField == -2 {
		return nil, errors.New(fmt.Sprintf("Invalid sort field: %s\n", field))
	}

	tt := *t
	sort.Sort(By(tt))
	return &tt, nil
}
