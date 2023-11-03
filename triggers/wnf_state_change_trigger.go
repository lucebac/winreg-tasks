// SPDX-License-Identifier: MIT

package triggers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lucebac/winreg-tasks/generated"
	"github.com/lucebac/winreg-tasks/utils"
)

const WnfStateChangeTriggerMagic TriggerMagic = 0x6666

type WnfStateChangeTrigger struct {
	GenericData *GenericTriggerData
	StateName   []byte
	Data        []byte
}

func NewWnfStateChangeTrigger(gen *generated.Triggers_WnfStateChangeTrigger, tz *time.Location) (*WnfStateChangeTrigger, error) {
	generic, err := NewGenericTriggerData(gen.GenericData, tz)
	if err != nil {
		return nil, err
	}

	return &WnfStateChangeTrigger{
		GenericData: generic,
		StateName:   gen.StateName[:],
		Data:        gen.Data[:],
	}, nil
}

func IsWnfStateChangeTrigger(trigger Trigger) bool {
	return trigger.Magic() == WnfStateChangeTriggerMagic
}

func (t WnfStateChangeTrigger) Id() string {
	return t.GenericData.TriggerId
}

func (t WnfStateChangeTrigger) Magic() TriggerMagic {
	return WnfStateChangeTriggerMagic
}

func (t WnfStateChangeTrigger) Name() string {
	return "WnfStateChange"
}

func (t WnfStateChangeTrigger) String() string {
	return fmt.Sprintf(
		`<WnfStateChange state_name="%s" data="%s">`,
		utils.Hexdump(t.StateName, len(t.StateName)), utils.Hexdump(t.Data, len(t.Data)),
	)
}

func (t WnfStateChangeTrigger) MarshalJSON() ([]byte, error) {
	var s struct {
		GenericData *GenericTriggerData
		StateName   string
		Data        string
	}

	s.GenericData = t.GenericData
	s.StateName = utils.Hexdump(t.StateName, len(t.StateName))
	s.Data = string(t.Data)

	return json.Marshal(s)
}
