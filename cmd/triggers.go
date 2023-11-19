// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"time"

	"github.com/alecthomas/kong"
	"github.com/lucebac/winreg-tasks/triggers"
	"github.com/lucebac/winreg-tasks/utils"
)

type triggersCommand struct {
	Dump   bool   `help:"Dumps the contents of the Triggers Value in hex." optional:"" default:"false" short:"d"`
	TaskId string `help:"The UUID of the Task" arg:""`
}

func (t *triggersCommand) Run(ctx *context) error {
	triggersRaw, err := ctx.provider.GetTriggers(t.TaskId)
	if err != nil {
		return fmt.Errorf("cannot get triggers for task (%v)", err)
	}

	if t.Dump {
		hex := utils.Hexdump(triggersRaw, 16)
		fmt.Println(hex)
	}

	triggers, err := triggers.FromBytes(triggersRaw, time.Local)
	if err != nil {
		return fmt.Errorf("cannot parse triggers (%v)", err)
	}

	fmt.Println("Header:")
	fmt.Printf("\tVersion: %d", triggers.Header.Version)
	fmt.Printf("\tStartBoundary: %s", triggers.Header.StartBoundary.String())
	fmt.Printf("\tEndBoundary: %s", triggers.Header.EndBoundary.String())

	fmt.Println("JobBucket:")
	fmt.Printf("\tFlags: %08x", triggers.JobBucket.Flags)
	fmt.Printf("\tCRC32: %08x", triggers.JobBucket.Crc32)
	fmt.Printf("\tPrincipal ID: %s", triggers.JobBucket.PrincipalId)
	fmt.Printf("\tDisplay Name: %s", triggers.JobBucket.DisplayName)
	fmt.Printf("\tUser: %s", triggers.JobBucket.UserInfo.UserToString())

	fmt.Println("Triggers:")
	if len(triggers.Triggers) == 0 {
		fmt.Println("\t<no triggers>")
		return nil
	}

	for _, trigger := range triggers.Triggers {
		fmt.Println("\t" + trigger.String())
	}

	return nil
}

func init() {
	registerSubcommand(kong.DynamicCommand("triggers", "Dump the Triggers of a given Task.", "", &triggersCommand{}))
}
