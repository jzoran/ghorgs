// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package main

// Protocol interface represents methods to create and
// execute query and create and execture a response to
// that query
type Protocol interface {
	getName() string

	makeCsv() *Csv
	appendCsv(c *Csv)
	getCsvTitle() []string

	makeQuery(org string) IQuery

	fromJsonBuffer(buff []byte)
	getTotal() int
	hasNext() bool
	getNext() string

	toString() string
}
