// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/lucebac/winreg-tasks/actions"
	"github.com/lucebac/winreg-tasks/dynamicinfo"
	"github.com/lucebac/winreg-tasks/triggers"
	"github.com/rs/zerolog/log"
)

func readAndParse(taskId string, rawData []byte, value string, parserCallback func(data []byte) (string, error), quiet bool) error {
	result, err := parserCallback(rawData)
	if err != nil {
		return err
	}

	if !quiet {
		log.Info().Str("taskId", taskId).Msg(fmt.Sprintf("%s: %s", value, result))
	}

	return nil
}

func parseTriggers(data []byte) (string, error) {
	triggs, err := triggers.FromBytes(data, time.Local)
	if err != nil {
		return "", err
	}

	if len(triggs.Triggers) == 0 {
		return "<no triggers>", nil
	}

	var triggers []string

	for _, trigger := range triggs.Triggers {
		triggers = append(triggers, trigger.Name())
	}

	return strings.Join(triggers, ", "), nil
}

func parseActions(data []byte) (string, error) {
	actions, err := actions.FromBytes(data)
	if err != nil {
		return "", err
	}

	if len(actions.Properties) == 0 {
		return "<no actions>", nil
	}

	var propCollection []string

	for _, props := range actions.Properties {
		propCollection = append(propCollection, props.Name())
	}

	return strings.Join(propCollection, ", "), nil
}

func parseDynamicInfo(data []byte) (string, error) {
	dynamicInfo, err := dynamicinfo.FromBytes(data)
	if err != nil {
		return "", err
	}

	ret := fmt.Sprintf(
		"Creation Time: %s, Last Run Time: %s, Last Error Code: 0x%08x",
		dynamicInfo.CreationTime.String(), dynamicInfo.LastRunTime.String(),
		dynamicInfo.LastErrorCode,
	)

	return ret, nil
}

type parseallCommand struct {
	Quiet bool `help:"Don't output the task list, just check for parsing errors" optional:"" default:"false" short:"q"`
}

func (p *parseallCommand) Run(ctx *context) error {
	tasks, err := ctx.provider.GetTaskIdList()
	if err != nil {
		return fmt.Errorf("cannot get task list from registry (%v)", err)
	}

	for _, taskId := range tasks {
		if data, err := ctx.provider.GetActions(taskId); err == nil {
			if err := readAndParse(taskId, data, "Actions", parseActions, p.Quiet); err != nil {
				log.Error().Err(err).Str("taskId", taskId).Msg("cannot parse Actions")
			}
		} else {
			log.Error().Err(err).Str("taskId", taskId).Msg("cannot retrieve Actions data")
		}

		if data, err := ctx.provider.GetTriggers(taskId); err == nil {
			if err := readAndParse(taskId, data, "Triggers", parseTriggers, p.Quiet); err != nil {
				log.Error().Err(err).Str("taskId", taskId).Msg("cannot parse Triggers")
			}
		} else {
			log.Error().Err(err).Str("taskId", taskId).Msg("cannot retrieve Triggers data")
		}

		if data, err := ctx.provider.GetDynamicInfo(taskId); err == nil {
			if err := readAndParse(taskId, data, "DynamicInfo", parseDynamicInfo, p.Quiet); err != nil {
				log.Error().Err(err).Str("taskId", taskId).Msg("cannot parse DynamicInfo")
			}
		} else {
			log.Error().Err(err).Str("taskId", taskId).Msg("cannot retrieve DynamicInfo data")
		}
	}

	return nil
}

func init() {
	registerSubcommand(kong.DynamicCommand("parseall", "Parses all existing Tasks", "", &parseallCommand{}))
}
