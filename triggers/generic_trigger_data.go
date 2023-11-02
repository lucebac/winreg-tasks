// SPDX-License-Identifier: MIT

package triggers

import (
	"time"

	"github.com/lucebac/winreg-tasks/generated"
	"github.com/lucebac/winreg-tasks/utils"
)

type GenericTriggerData struct {
	StartBoundary       time.Time
	EndBoundary         time.Time
	Delay               utils.Duration
	Timeout             utils.Duration
	RepetitionInterval  utils.Duration
	RepetitionDuration  utils.Duration
	RepetitionDuration2 utils.Duration
	StopAtDurationEnd   bool
	Enabled             bool
	Unknown             []byte
	TriggerId           string
}

func NewGenericTriggerData(gen *generated.Triggers_GenericTriggerData, tz *time.Location) (*GenericTriggerData, error) {
	triggerId := ""
	if gen.TriggerId != nil {
		triggerId = gen.TriggerId.Str
	}

	return &GenericTriggerData{
		StartBoundary:       utils.TimeFromTSTime(gen.StartBoundary, tz),
		EndBoundary:         utils.TimeFromTSTime(gen.EndBoundary, tz),
		Delay:               utils.SecondsToDuration(gen.DelaySeconds),
		Timeout:             utils.SecondsToDuration(gen.TimeoutSeconds),
		RepetitionInterval:  utils.SecondsToDuration(gen.RepetitionIntervalSeconds),
		RepetitionDuration:  utils.SecondsToDuration(gen.RepetitionDurationSeconds),
		RepetitionDuration2: utils.SecondsToDuration(gen.RepetitionDurationSeconds2),
		StopAtDurationEnd:   gen.StopAtDurationEnd != 0,
		Enabled:             gen.Enabled.Value != 0,
		Unknown:             gen.Unknown[:],
		TriggerId:           triggerId,
	}, nil
}
