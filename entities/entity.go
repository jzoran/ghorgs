// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package entities

import "ghorgs/cache"

// Entity interface represents methods to create and
// execute query and create and execute a response to
// that query
type Entity interface {
	GetName() string

	MakeTable() *cache.Table
	AppendTable(c *cache.Table)
	GetTableFields() []string
	HasField(s string) bool
	GetCsvFile() string

	MakeQuery(org string) IQuery
	FromJsonBuffer(buff []byte)
	GetTotal() int
	HasNext() bool
	GetNext() string

	ToString() string
}
