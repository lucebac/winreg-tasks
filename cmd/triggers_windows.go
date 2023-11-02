// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/alecthomas/kong"
	"github.com/lucebac/winreg-tasks/triggers"
	"github.com/lucebac/winreg-tasks/utils"
)

type triggersCommand struct {
	Dump   bool   `help:"Dumps the contents of the Triggers Value in hex." optional:"" default:"false" short:"d"`
	TaskId string `help:"The UUID of the Task" arg:""`
}

func (t *triggersCommand) Run(ctx *kong.Context) error {
	key := openTaskKey(t.TaskId)
	if key == 0 {
		return fmt.Errorf("cannot open task key")
	}
	defer key.Close()

	triggersRaw, _, err := key.GetBinaryValue("Triggers")
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

	log.Println("Header:")
	log.Printf("\tVersion: %d", triggers.Header.Version)
	log.Printf("\tStartBoundary: %s", triggers.Header.StartBoundary.String())
	log.Printf("\tEndBoundary: %s", triggers.Header.EndBoundary.String())

	log.Println("JobBucket:")
	log.Printf("\tFlags: %08x", triggers.JobBucket.Flags)
	log.Printf("\tCRC32: %08x", triggers.JobBucket.Crc32)
	log.Printf("\tPrincipal ID: %s", triggers.JobBucket.PrincipalId)
	log.Printf("\tDisplay Name: %s", triggers.JobBucket.DisplayName)
	log.Printf("\tUser: %s", triggers.JobBucket.UserInfo.UserToString())

	log.Println("Triggers:")
	if len(triggers.Triggers) == 0 {
		log.Println("\t<no triggers>")
		return nil
	}

	for _, trigger := range triggers.Triggers {
		log.Println("\t" + trigger.String())
	}

	return nil
}

func init() {
	registerSubcommand(kong.DynamicCommand("triggers", "Dump the Triggers of a given Task.", "", &triggersCommand{}))
}
