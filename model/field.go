//
// Copyright (c) 2019 Sony Mobile Communications Inc.
// SPDX-License-Identifier: MIT
//

package model

type Field struct {
	Name  string
	Index int
}

type Fields interface {
	asList() []Field
	DisplayNames() []string
}

var (
	// Since a single record in Table.Records is
	// a slice of strings with slice indices mapping
	// exactly to fields (excluding "Id"), these
	// fields are indexed from 0 to len(record).
	// "Id" field itself is a key for that record
	// (and not part of the records slice), so we
	// give "Id" index -1 by convention.
	INVALID_FIELD Field = Field{"INVALID_FIELD", -2}
	ID            Field = Field{"Id", -1}
)

func namesOf(fields []Field) []string {
	names := make([]string, 0)
	names = append(names, ID.Name)
	for _, field := range fields {
		names = append(names, field.Name)
	}
	return names
}
