// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package main

import (
	"log"
)

type Table struct {
	Records map[string][]string
	Keys    []string
}

func makeTable() *Table {
	keys := make([]string, 0)
	return &Table{nil, keys}
}

func (t *Table) addKey(key string) {
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
