// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"log"

	"github.com/alecthomas/kong"
	"github.com/lucebac/winreg-tasks/dynamicinfo"
	"github.com/lucebac/winreg-tasks/utils"
)

type dynamicInfoCommand struct {
	Dump   bool   `help:"Dumps the contents of the DynamicInfo Value in hex." optional:"" default:"false" short:"d"`
	TaskId string `help:"The UUID of the Task" arg:""`
}

func (d *dynamicInfoCommand) Run(ctx *kong.Context) error {
	key := openTaskKey(d.TaskId)
	if key == 0 {
		return fmt.Errorf("cannot open task key")
	}
	defer key.Close()

	dynamicInfoRaw, _, err := key.GetBinaryValue("DynamicInfo")
	if err != nil {
		return fmt.Errorf("cannot get dynamic info for task (%v)", err)
	}

	if d.Dump {
		hex := utils.Hexdump(dynamicInfoRaw, 16)
		fmt.Println(hex)
	}

	dynamicInfo, err := dynamicinfo.FromBytes(dynamicInfoRaw)
	if err != nil {
		return fmt.Errorf("cannot parse DynamicInfo (%v)", err)
	}

	log.Printf("Creation Time: %s", dynamicInfo.CreationTime.String())
	log.Printf("Last Run Time: %s", dynamicInfo.LastRunTime.String())
	log.Printf("Task State: 0x%08x", dynamicInfo.TaskState)
	log.Printf("Last Error Code: 0x%08x", dynamicInfo.LastErrorCode)
	log.Printf("Last Successful Run Time: %s", dynamicInfo.LastSuccessfulRunTime.String())

	return nil
}

func init() {
	registerSubcommand(kong.DynamicCommand("dynamicinfo", "Dump the DynamicInfo of a given Task.", "", &dynamicInfoCommand{}))
}
