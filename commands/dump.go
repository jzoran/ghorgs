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

		tt, err := t.SortByColumn(d.By)
		if err != nil {
			panic(err)
		}

		fmt.Printf("\nDumping %s...", filename)
		csv := &cache.Csv{filename, tt}
		csv.Flush(proto.GetCsvTitle())
	}

	return nil
}
