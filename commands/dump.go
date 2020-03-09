// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package commands

import (
	"fmt"
	"ghorgs/cache"
	"ghorgs/protocols"
)

type Dump struct {
	By   string
	data map[string]*cache.Table
}

func (d *Dump) AddCache(c map[string]*cache.Table) {
	d.data = c
}

func (d *Dump) Do(protoMap map[string]protocols.Protocol) error {
	for name, t := range d.data {
		proto := protoMap[name]
		filename := proto.GetCsvFile()

		fmt.Printf("\nDumping %s...", filename)
		if d.By != "" {
			_, err := t.SortByField(d.By)
			if err != nil {
				panic(err)
			}
		}
		csv := &cache.Csv{filename, t}
		csv.Flush(proto.GetTableFields())
	}

	return nil
}
