package providers

import (
	"errors"
	"time"
)

type DataProvider interface {
	GetActions(taskId string) ([]byte, error)
	GetTriggers(taskId string) ([]byte, error)
	GetDynamicInfo(taskId string) ([]byte, error)

	GetStringField(taskId, fieldName string) (string, error)
	GetBytesField(taskId, fieldName string) ([]byte, error)
	GetDwordField(taskId, fieldName string) (uint32, error)
	GetDateField(taskId, fieldName string) (*time.Time, error)

	GetTaskIdList() ([]string, error)

	Close()
}

var (
	ErrFileNotFound = errors.New("file not found")
)

const (
	taskKeyBase = `Microsoft\Windows NT\CurrentVersion\Schedule\TaskCache\`
)
