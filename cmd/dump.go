// SPDX-License-Identifier: MIT

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/lucebac/winreg-tasks/task"
	"github.com/rs/zerolog/log"
)

type dumpCommand struct {
	Outfile  string `optional:"" help:"Path to file to write output to. Otherwise the output is written to stdout." short:"o"`
	Timezone string `optional:"" help:"The timezone of the system where the Registry Hive is from. Default: Local timezone of this machine." aliases:"tz" short:"t"`
}

func (d *dumpCommand) Run(ctx *context) error {
	var tz *time.Location = time.Local
	var err error

	if d.Timezone != "" {
		if tz, err = parseLocationRelativeUTC(d.Timezone); err != nil {
			return fmt.Errorf(`cannot parse UTC offset (%v)`, err)
		}
	}

	taskIdList, err := ctx.provider.GetTaskIdList()
	if err != nil {
		return err
	}

	var parsedTasks []task.Task

	for _, taskId := range taskIdList {
		task := task.NewTask(taskId, ctx.provider)

		if err := task.ParseAll(tz); err != nil {
			log.Error().Err(err).Str("taskId", taskId).Msg(`cannot parse task`)
			continue
		}

		parsedTasks = append(parsedTasks, task)
	}

	var f *os.File = os.Stdout

	if d.Outfile != "" && d.Outfile != "-" {
		f, err = os.OpenFile(d.Outfile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		if err != nil {
			return fmt.Errorf("cannot open output file (%v)", err)
		}
		defer f.Close()
	}

	m := json.NewEncoder(f)

	for _, t := range parsedTasks {
		if err := m.Encode(t); err != nil {
			return fmt.Errorf(`cannot serialize task %s (%v)`, t.Id, err)
		}
	}

	return nil
}

func init() {
	registerSubcommand(kong.DynamicCommand("dump", "Dump the Task list to a file", "", &dumpCommand{}))
}
