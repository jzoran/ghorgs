// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package protocols

import "ghorgs/cache"

// Protocol interface represents methods to create and
// execute query and create and execture a response to
// that query
type Protocol interface {
	GetName() string

	MakeTable() *cache.Table
	AppendTable(c *cache.Table)
	GetCsvFile() string
	GetCsvTitle() []string

	MakeQuery(org string) IQuery

	FromJsonBuffer(buff []byte)
	GetTotal() int
	HasNext() bool
	GetNext() string

	ToString() string
}
