// SPDX-License-Identifier: MIT
package task

import (
	"time"

	"github.com/google/uuid"
	"github.com/lucebac/winreg-tasks/actions"
	"github.com/lucebac/winreg-tasks/dynamicinfo"
	"github.com/lucebac/winreg-tasks/providers"
	"github.com/lucebac/winreg-tasks/triggers"
)

type Task struct {
	provider providers.DataProvider
	stringId string

	Id uuid.UUID

	DynamicInfo *dynamicinfo.DynamicInfo

	Actions  *actions.Actions
	Triggers *triggers.Triggers
}

func NewTask(id string, provider providers.DataProvider) Task {
	return Task{
		provider: provider,
		stringId: id,
		Id:       uuid.MustParse(id),
	}
}

func (t *Task) ParseAll(tz *time.Location) error {
	if _, err := t.GetActions(); err != nil {
		return err
	}

	if _, err := t.GetTriggers(tz); err != nil {
		return err
	}

	if _, err := t.GetDynamicInfo(); err != nil {
		return err
	}

	return nil
}

func (t *Task) GetActions() (*actions.Actions, error) {
	if t.Actions == nil {
		rawData, err := t.provider.GetActions(t.stringId)
		if err != nil {
			return nil, err
		}

		t.Actions, err = actions.FromBytes(rawData)
		if err != nil {
			return nil, err
		}
	}

	return t.Actions, nil
}

func (t *Task) GetTriggers(tz *time.Location) (*triggers.Triggers, error) {
	if t.Triggers == nil {
		rawData, err := t.provider.GetTriggers(t.stringId)
		if err != nil {
			return nil, err
		}

		t.Triggers, err = triggers.FromBytes(rawData, tz)
		if err != nil {
			return nil, err
		}
	}

	return t.Triggers, nil
}

func (t *Task) GetDynamicInfo() (*dynamicinfo.DynamicInfo, error) {
	if t.DynamicInfo == nil {
		rawData, err := t.provider.GetDynamicInfo(t.stringId)
		if err != nil {
			return nil, err
		}

		t.DynamicInfo, err = dynamicinfo.FromBytes(rawData)
		if err != nil {
			return nil, err
		}
	}
	return t.DynamicInfo, nil
}
