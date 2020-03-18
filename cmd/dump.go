// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package cmd

import (
	"fmt"
	"ghorgs/cache"
	"ghorgs/entities"
	cmds "github.com/spf13/cobra"
)

type dumpT struct {
	entities []string
	by       string
	data     map[string]*cache.Table
}

var d = &dumpT{}

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
			"    "+sliceToStr(entities.EntityList)+".")

	dumpCmd.Flags().StringP("by",
		"b",
		"",
		"Name of the entity field to use for sorting the result of the dump.\n"+
			"If empty, default sort on GitHub is creation date.")

	rootCmd.AddCommand(dumpCmd)
}

func (d *dumpT) addCache(c map[string]*cache.Table) {
	d.data = c
}

func (d *dumpT) validateArgs(c *cmds.Command, args []string) error {
	ents, err := c.Flags().GetString("entities")
	if err != nil {
		panic(err)
	}
	d.entities, err = entities.ValidateEntities(ents)
	if err != nil {
		return err
	}

	d.by, err = c.Flags().GetString("by")
	if err != nil {
		panic(err)
	}
	err = entities.ValidateEntityField(d.by, d.entities)
	if err != nil {
		return err
	}

	return nil
}

func (d *dumpT) run(c *cmds.Command, args []string) {
	d.addCache(Cache(d.entities, entities.EntityMap))
	for name, t := range d.data {
		entity := entities.EntityMap[name]
		filename := entity.GetCsvFile()

		fmt.Printf("\nDumping %s...", filename)
		if d.by != "" {
			_, err := t.SortByField(d.by)
			if err != nil {
				panic(err)
			}
		}
		csv := &cache.Csv{filename, t}
		csv.Flush(entity.GetTableFields())
	}
}

func sliceToStr(sl []string) string {
	var str = ""
	for _, item := range sl {
		str = str + item + ", "
	}
	return str[:len(str)-2]
}
