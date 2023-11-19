// SPDX-License-Identifier: MIT

package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/lucebac/winreg-tasks/actions"
	"github.com/lucebac/winreg-tasks/utils"
)

type actionsCommand struct {
	Dump   bool   `help:"Dumps the contents of the Actions Value in hex." optional:"" default:"false" short:"d"`
	TaskId string `help:"The UUID of the Task" arg:""`
}

func (a *actionsCommand) Run(ctx *context) error {
	actionsRaw, err := ctx.provider.GetActions(a.TaskId)
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

	fmt.Println("Context: " + actions.Context)
	fmt.Println(`Actions:`)

	if len(actions.Properties) == 0 {
		fmt.Println("\t<no actions>")
		return nil
	}

	for _, props := range actions.Properties {
		fmt.Println("\t" + props.String())
	}

	return nil
}

func init() {
	registerSubcommand(kong.DynamicCommand("actions", "Dump the Actions of a given Task.", "", &actionsCommand{}))
}
