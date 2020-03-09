// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package commands

import (
	"fmt"
	"ghorgs/cache"
	"ghorgs/entities"
)

type Dump struct {
	By   string
	data map[string]*cache.Table
}

func (d *Dump) AddCache(c map[string]*cache.Table) {
	d.data = c
}

func (d *Dump) Do(entityMap map[string]entities.Entity) error {
	for name, t := range d.data {
		entity := entityMap[name]
		filename := entity.GetCsvFile()

		fmt.Printf("\nDumping %s...", filename)
		if d.By != "" {
			_, err := t.SortByField(d.By)
			if err != nil {
				panic(err)
			}
		}
		csv := &cache.Csv{filename, t}
		csv.Flush(entity.GetTableFields())
	}

	return nil
}
