// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"log"

	"github.com/alecthomas/kong"
	"github.com/lucebac/winreg-tasks/actions"
	"github.com/lucebac/winreg-tasks/utils"
)

type actionsCommand struct {
	Dump   bool   `help:"Dumps the contents of the Actions Value in hex." optional:"" default:"false" short:"d"`
	TaskId string `help:"The UUID of the Task" arg:""`
}

func (a *actionsCommand) Run(ctx *kong.Context) error {
	key := openTaskKey(a.TaskId)
	if key == 0 {
		return fmt.Errorf("cannot open task key")
	}
	defer key.Close()

	actionsRaw, _, err := key.GetBinaryValue("Actions")
	if err != nil {
		return fmt.Errorf("cannot get actions for task (%v)", err)
	}

	if a.Dump {
		hex := utils.Hexdump(actionsRaw, 16)
		fmt.Println(hex)
	}

	actions, err := actions.FromBytes(actionsRaw)
	if err != nil {
		return fmt.Errorf("cannot parse actions (%v)", err)
	}

	log.Println("Context: " + actions.Context)
	log.Println(`Actions:`)

	if len(actions.Properties) == 0 {
		log.Println("\t<no actions>")
		return nil
	}

	for _, props := range actions.Properties {
		log.Println("\t" + props.String())
	}

	return nil
}

func init() {
	registerSubcommand(kong.DynamicCommand("actions", "Dump the Actions of a given Task.", "", &actionsCommand{}))
}
