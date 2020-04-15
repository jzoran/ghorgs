// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package cmd

import (
	"fmt"
	"ghorgs/model"
	"ghorgs/view"
	cmds "github.com/spf13/cobra"
)

type dumper struct {
	entities []model.Entity
	by       string
	data     map[string]*model.Table
}

var d = &dumper{}

var dumpCmd = &cmds.Command{
	Use:   "dump",
	Short: "Dumps the requested entities into a csv file.",
	Long:  `Dumps the requested entities into a csv file.`,
	Args:  d.validateArgs,
	Run:   d.run,
}

func init() {
	dumpCmd.Flags().StringP("entities",
		"e",
		"all",
		"'all' for full dump or comma separated list of one or more of:\n"+
			"    "+sliceToStr(model.EntityNamesList)+".")

	dumpCmd.Flags().StringP("by",
		"b",
		"",
		"Name of the entity field to use for sorting the result of the dump.\n"+
			"If empty, default sort on GitHub is creation date.")

	rootCmd.AddCommand(dumpCmd)
}

func (d *dumper) addCache(c map[string]*model.Table) {
	d.data = c
}

func (d *dumper) validateArgs(c *cmds.Command, args []string) error {
	ents, err := c.Flags().GetString("entities")
	if err != nil {
		panic(err)
	}
	d.entities, err = model.ValidateEntities(ents)
	if err != nil {
		return err
	}

	d.by, err = c.Flags().GetString("by")
	if err != nil {
		panic(err)
	}
	err = model.ValidateEntityField(d.by, d.entities)
	if err != nil {
		return err
	}

	return nil
}

func (d *dumper) run(c *cmds.Command, args []string) {
	ca, err := Cache(d.entities)
	if err != nil {
		fmt.Println("Error! %s", err.Error())
		return
	}

	d.addCache(ca)
	for name, t := range d.data {
		entity := model.EntityMap[name]
		filename := entity.GetCsvFile()

		fmt.Printf("\nDumping %s...", filename)
		if d.by != "" {
			_, err := t.SortByField(d.by)
			if err != nil {
				panic(err)
			}
		}
		csv := &view.Csv{filename, t}
		csv.Flush()
	}
}

func sliceToStr(sl []string) string {
	var str = ""
	for _, item := range sl {
		str = str + item + ", "
	}
	return str[:len(str)-2]
}
