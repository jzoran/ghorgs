// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package commands

import (
	"ghorgs/cache"
	"ghorgs/entities"
)

type Command interface {
	AddCache(c map[string]*cache.Table)
	Do(entityMap map[string]entities.Entity) error
}
