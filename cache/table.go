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
	Records map[string][]string
	Keys    []string
	Columns []string
	sortCol int
}

func MakeTable(columns []string) *Table {
	keys := make([]string, 0)
	return &Table{Records: nil, Keys: keys, Columns: columns}
}

func (t *Table) AddKey(key string) {
	t.Keys = append(t.Keys, key)
}

func (t *Table) log() {
	var s string
	for id, line := range t.Records {
		s = id + "\t"
		for i, cell := range line {
			s = s + cell
			if i < len(line)-1 {
				s = s + "\t"
			}
		}
		log.Print(s + "\n")
	}
}

// sort interface + method
type By Table

func (a By) Len() int { return len(a.Keys) }
func (a By) Less(i, j int) bool {
	if a.sortCol == -1 {
		return a.Keys[i] < a.Keys[j]
	}
	return a.Records[a.Keys[i]][a.sortCol] < a.Records[a.Keys[j]][a.sortCol]
}
func (a By) Swap(i, j int) {
	a.Keys[i], a.Keys[j] = a.Keys[j], a.Keys[i]
}

func (t *Table) SortByColumn(column string) (*Table, error) {
	t.sortCol = -1
	for i, val := range t.Columns {
		if val == column {
			t.sortCol = i - 1 // Columns include "Id" at index 0
			break
		}
	}
	if t.sortCol == -2 {
		return nil, errors.New(fmt.Sprintf("Invalid sort column: %s\n", column))
	}

	tt := *t
	sort.Sort(By(tt))
	return &tt, nil
}
