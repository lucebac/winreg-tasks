// SPDX-License-Identifier: MIT
package task

import (
	"log"
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
		log.Printf("cannot get Actions of Task %s: %v", t.stringId, err)
	}

	if _, err := t.GetTriggers(tz); err != nil {
		log.Printf("cannot get Triggers of Task %s: %v", t.stringId, err)
	}

	if _, err := t.GetDynamicInfo(); err != nil {
		log.Printf("cannot get DynamicInfo of Task %s: %v", t.stringId, err)
	}

	var err error

	if t.Author, err = t.provider.GetStringField(t.stringId, "Author"); err != nil {
		log.Printf("cannot get Author of task %s: %v", t.stringId, err)
	}
	if t.Data, err = t.provider.GetBytesField(t.stringId, "Data"); err != nil {
		log.Printf("cannot get Data of task %s: %v", t.stringId, err)
	}
	if t.Date, err = t.provider.GetDateField(t.stringId, "Date"); err != nil {
		log.Printf("cannot get Date of task %s: %v", t.stringId, err)
	}
	if t.Description, err = t.provider.GetStringField(t.stringId, "Description"); err != nil {
		log.Printf("cannot get Description of task %s: %v", t.stringId, err)
	}
	if t.Documentation, err = t.provider.GetStringField(t.stringId, "Documentation"); err != nil {
		log.Printf("cannot get Documentation of task %s: %v", t.stringId, err)
	}
	if t.Hash, err = t.provider.GetBytesField(t.stringId, "Hash"); err != nil {
		log.Printf("cannot get Hash of task %s: %v", t.stringId, err)
	}
	if t.Path, err = t.provider.GetStringField(t.stringId, "Path"); err != nil {
		log.Printf("cannot get Path of task %s: %v", t.stringId, err)
	}
	if t.Schema, err = t.provider.GetDwordField(t.stringId, "Schema"); err != nil {
		log.Printf("cannot get Schema of task %s: %v", t.stringId, err)
	}
	if t.SecurityDescriptor, err = t.provider.GetStringField(t.stringId, "SecurityDescriptor"); err != nil {
		log.Printf("cannot get SecurityDescriptor of task %s: %v", t.stringId, err)
	}
	if t.Source, err = t.provider.GetStringField(t.stringId, "Source"); err != nil {
		log.Printf("cannot get Source of task %s: %v", t.stringId, err)
	}
	if t.URI, err = t.provider.GetStringField(t.stringId, "URI"); err != nil {
		log.Printf("cannot get URI of task %s: %v", t.stringId, err)
	}
	if t.Version, err = t.provider.GetStringField(t.stringId, "Version"); err != nil {
		log.Printf("cannot get Version of task %s: %v", t.stringId, err)
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
