// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package cache

import (
	"os"
	"reflect"
)

type Csv struct {
	FileName string
	Data     *Table
}

func MakeCsv(filename string) *Csv {
	data := MakeTable([]string{})
	return &Csv{filename, data}
}

func (c *Csv) Flush(title []string) {
	var f *os.File

	if _, err := os.Stat(c.FileName); err == nil || os.IsNotExist(err) {
		f, err = os.OpenFile(c.FileName,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
	} else {
		panic(err)
	}

	s := lineToString(title)
	if _, err := f.WriteString(s + "\n"); err != nil {
		panic(err)
	}

	for _, k := range c.Data.Keys {
		// if by any chance we've loaded a dump from disk
		// and it includes title line, just skip it now
		if k == title[0] && reflect.DeepEqual(c.Data.Records[k], title[1:]) {
			continue
		}

		s = k + "\t" + lineToString(c.Data.Records[k])
		if _, err := f.WriteString(s + "\n"); err != nil {
			panic(err)
		}
	}
}

func lineToString(line []string) string {
	s := ""
	for i, cell := range line {
		s = s + cell
		if i < len(line)-1 {
			s = s + "\t"
		}
	}

	return s
}

func (c *Csv) Log() {
	c.Data.Log()
}
