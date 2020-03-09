// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package commands

import (
	"ghorgs/cache"
	"ghorgs/entities"
)

type Remove struct {
	N     int
	Since string
	Names string
	data  map[string]*cache.Table
}

func (r *Remove) AddCache(c map[string]*cache.Table) {
	r.data = c
}

func (r *Remove) Do(entityMap map[string]entities.Entity) error {
	return nil
}
