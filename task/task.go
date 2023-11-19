// SPDX-License-Identifier: MIT
package task

import (
	"time"

	"github.com/google/uuid"
	"github.com/lucebac/winreg-tasks/actions"
	"github.com/lucebac/winreg-tasks/dynamicinfo"
	"github.com/lucebac/winreg-tasks/providers"
	"github.com/lucebac/winreg-tasks/triggers"
	"github.com/rs/zerolog/log"
)

type Task struct {
	provider providers.DataProvider
	stringId string

	Id uuid.UUID

	DynamicInfo *dynamicinfo.DynamicInfo

	Actions  *actions.Actions
	Triggers *triggers.Triggers

	Author             string     `json:",omitempty"`
	Data               []byte     `json:",omitempty"`
	Date               *time.Time `json:",omitempty"`
	Description        string     `json:",omitempty"`
	Documentation      string     `json:",omitempty"`
	Hash               []byte     `json:",omitempty"`
	Path               string     `json:",omitempty"`
	Schema             uint32     `json:",omitempty"`
	SecurityDescriptor string     `json:",omitempty"`
	Source             string     `json:",omitempty"`
	URI                string     `json:",omitempty"`
	Version            string     `json:",omitempty"`
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
		log.Debug().Err(err).Str("taskId", t.stringId).Msg("cannot get Actions")
	}

	if _, err := t.GetTriggers(tz); err != nil {
		log.Debug().Err(err).Str("taskId", t.stringId).Msg("cannot get Triggers")
	}

	if _, err := t.GetDynamicInfo(); err != nil {
		log.Debug().Err(err).Str("taskId", t.stringId).Msg("cannot get DynamicInfo")
	}

	var err error

	if t.Author, err = t.provider.GetStringField(t.stringId, "Author"); err != nil {
		log.Debug().Err(err).Str("taskId", t.stringId).Msg("cannot get Author")
	}
	if t.Data, err = t.provider.GetBytesField(t.stringId, "Data"); err != nil {
		log.Debug().Err(err).Str("taskId", t.stringId).Msg("cannot get Data")
	}
	if t.Date, err = t.provider.GetDateField(t.stringId, "Date"); err != nil {
		log.Debug().Err(err).Str("taskId", t.stringId).Msg("cannot get Date")
	}
	if t.Description, err = t.provider.GetStringField(t.stringId, "Description"); err != nil {
		log.Debug().Err(err).Str("taskId", t.stringId).Msg("cannot get Description")
	}
	if t.Documentation, err = t.provider.GetStringField(t.stringId, "Documentation"); err != nil {
		log.Debug().Err(err).Str("taskId", t.stringId).Msg("cannot get Documentation")
	}
	if t.Hash, err = t.provider.GetBytesField(t.stringId, "Hash"); err != nil {
		log.Debug().Err(err).Str("taskId", t.stringId).Msg("cannot get Hash")
	}
	if t.Path, err = t.provider.GetStringField(t.stringId, "Path"); err != nil {
		log.Debug().Err(err).Str("taskId", t.stringId).Msg("cannot get Path")
	}
	if t.Schema, err = t.provider.GetDwordField(t.stringId, "Schema"); err != nil {
		log.Debug().Err(err).Str("taskId", t.stringId).Msg("cannot get Schema")
	}
	if t.SecurityDescriptor, err = t.provider.GetStringField(t.stringId, "SecurityDescriptor"); err != nil {
		log.Debug().Err(err).Str("taskId", t.stringId).Msg("cannot get SecurityDescriptor")
	}
	if t.Source, err = t.provider.GetStringField(t.stringId, "Source"); err != nil {
		log.Debug().Err(err).Str("taskId", t.stringId).Msg("cannot get Source")
	}
	if t.URI, err = t.provider.GetStringField(t.stringId, "URI"); err != nil {
		log.Debug().Err(err).Str("taskId", t.stringId).Msg("cannot get URI")
	}
	if t.Version, err = t.provider.GetStringField(t.stringId, "Version"); err != nil {
		log.Debug().Err(err).Str("taskId", t.stringId).Msg("cannot get Version")
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
