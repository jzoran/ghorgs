// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package commands

import (
	"ghorgs/cache"
	"ghorgs/protocols"
)

type Command interface {
	AddCache(c map[string]*cache.Table)
	Do(protoMap map[string]protocols.Protocol) error
}
