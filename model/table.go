// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package model

import (
	"errors"
	"fmt"
	"log"
	"sort"
)

type Table struct {
	Records    map[string][]string
	Keys       []string
	Fields     []Field
	pivotField Field
}

func MakeTable(fields []Field) *Table {
	keys := make([]string, 0)
	return &Table{Records: nil, Keys: keys, Fields: fields}
}

func (t *Table) FieldNames() []string {
	return namesOf(t.Fields)
}

func (t *Table) AddKey(key string) {
	t.Keys = append(t.Keys, key)
}

func (t *Table) AddRecord(key string, record []string) {
	if t.Records == nil {
		t.Records = make(map[string][]string)
	}
	t.Records[key] = record
}

func (t *Table) ToString() string {
	var s string
	for id, key := range t.Keys {
		s += key + "\t"
		line := t.Records[key]
		for i, cell := range line {
			s += cell
			if i < len(line)-1 {
				s += "\t"
			}
		}
		if id < len(t.Keys)-1 {
			s += "\n"
		}
	}

	return s
}

func (t *Table) Log() {
	log.Print(t.ToString() + "\n")
}

// sort interface + method
type By Table

func (a By) Len() int { return len(a.Keys) }
func (a By) Less(i, j int) bool {
	if a.pivotField.Index == ID.Index {
		return a.Keys[i] < a.Keys[j]
	}
	return a.Records[a.Keys[i]][a.pivotField.Index] < a.Records[a.Keys[j]][a.pivotField.Index]
}
func (a By) Swap(i, j int) {
	a.Keys[i], a.Keys[j] = a.Keys[j], a.Keys[i]
}

func (t *Table) SortByField(field string) (*Table, error) {
	err := t.setPivotField(field)
	if err != nil {
		return nil, err
	}

	tt := *t
	sort.Sort(By(tt))
	return &tt, nil
}

func (t *Table) FindByField(field string, val string) (*Table, error) {
	err := t.setPivotField(field)
	if err != nil {
		return nil, err
	}

	keyI := sort.Search(len(t.Keys), func(i int) bool {
		return val <= t.Records[t.Keys[i]][t.pivotField.Index]
	})

	if keyI < len(t.Keys) && val == t.Records[t.Keys[keyI]][t.pivotField.Index] {
		key := t.Keys[keyI]
		keys := []string{key}
		ret := &Table{Records: nil, Keys: keys, Fields: t.Fields}
		ret.AddRecord(key, t.Records[key])
		return ret, nil
	}

	return nil, fmt.Errorf("%s not found in field %s.", val, field)
}

// func (t *Table) LessThanByField(field string, val string) (*Table, error) {
// 	err := t.setPivotField(field)
//  if err != nil {
//      return nil, err
//  }
//
// 	keys := make([]string, 0)
// 	ret := &Table{Records: nil, Keys: keys, Fields: t.Fields}
// 	for _, key := range t.Keys {
// 		if t.Records[key][t.pivotField.Index] < val {
// 			ret.AddKey(key)
// 			ret.AddRecord(key, t.Records[key])
// 		}
// 	}

// 	return ret, nil
// }

func (t *Table) GreaterThanByField(field string, val string) (*Table, error) {
	err := t.setPivotField(field)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0)
	ret := &Table{Records: nil, Keys: keys, Fields: t.Fields}
	for _, key := range t.Keys {
		if t.Records[key][t.pivotField.Index] > val {
			ret.AddKey(key)
			ret.AddRecord(key, t.Records[key])
		}
	}

	return ret, nil
}

func (t *Table) Last(n int) (*Table, error) {
	if n < 1 || n > len(t.Keys) {
		return nil, fmt.Errorf("Out of range error.")
	}

	keys := make([]string, 0)
	ret := &Table{Records: nil, Keys: keys, Fields: t.Fields}
	for _, key := range t.Keys[len(t.Keys)-n:] {
		ret.AddKey(key)
		ret.AddRecord(key, t.Records[key])
	}

	return ret, nil
}

// func (t *Table) First(n int) (*Table, error) {
// 	if n <= 0 || n > len(t.Keys) {
// 		return nil, fmt.Errorf("Out of range error.")
// 	}

// 	keys := make([]string, 0)
// 	ret := &Table{Records: nil, Keys: keys, Fields: t.Fields}
// 	for _, key := range t.Keys[:len(t.Keys)-n] {
// 		ret.AddKey(key)
// 		ret.AddRecord(key, t.Records[key])
// 	}

// 	return ret, nil
// }

func (t *Table) setPivotField(fieldName string) error {
	if fieldName == ID.Name {
		t.pivotField = ID
		return nil
	}

	t.pivotField = INVALID_FIELD
	for _, field := range t.Fields {
		if field.Name == fieldName {
			t.pivotField = field
			break
		}
	}

	if t.pivotField == INVALID_FIELD {
		return errors.New(fmt.Sprintf("Invalid search field: %s\n", fieldName))
	}

	return nil
}
