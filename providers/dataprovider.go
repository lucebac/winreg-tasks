package providers

import "errors"

type DataProvider interface {
	GetActions(taskId string) ([]byte, error)
	GetTriggers(taskId string) ([]byte, error)
	GetDynamicInfo(taskId string) ([]byte, error)

	GetTaskIdList() ([]string, error)

	Close()
}

var (
	ErrFileNotFound = errors.New("file not found")
)

const (
	taskKeyBase = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Schedule\TaskCache\`
)
