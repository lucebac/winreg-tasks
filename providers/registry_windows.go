package providers

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/sys/windows/registry"
)

type WindowsRegistryProvider struct {
}

func NewWindowsRegistryProvider() (DataProvider, error) {
	return &WindowsRegistryProvider{}, nil
}

func (p WindowsRegistryProvider) GetActions(taskId string) ([]byte, error) {
	key := openTaskKey(taskId)
	if key == 0 {
		return nil, fmt.Errorf("cannot open task key")
	}
	defer key.Close()

	actionsRaw, _, err := key.GetBinaryValue("Actions")
	if err != nil {
		return nil, fmt.Errorf("cannot get actions for task (%v)", err)
	}

	return actionsRaw, nil
}

func (p WindowsRegistryProvider) GetTriggers(taskId string) ([]byte, error) {
	key := openTaskKey(taskId)
	if key == 0 {
		return nil, fmt.Errorf("cannot open task key")
	}
	defer key.Close()

	actionsRaw, _, err := key.GetBinaryValue("Triggers")
	if err != nil {
		return nil, fmt.Errorf("cannot get actions for task (%v)", err)
	}

	return actionsRaw, nil
}

func (p WindowsRegistryProvider) GetDynamicInfo(taskId string) ([]byte, error) {
	key := openTaskKey(taskId)
	if key == 0 {
		return nil, fmt.Errorf("cannot open task key")
	}
	defer key.Close()

	actionsRaw, _, err := key.GetBinaryValue("DynamicInfo")
	if err != nil {
		return nil, fmt.Errorf("cannot get actions for task (%v)", err)
	}

	return actionsRaw, nil
}

func (p WindowsRegistryProvider) GetTaskIdList() ([]string, error) {
	taskDir, err := openKey(`Tasks`)
	if err != nil {
		return nil, err
	}
	defer taskDir.Close()

	taskList, err := taskDir.ReadSubKeyNames(-1)
	if err != nil {
		return nil, fmt.Errorf("cannot get task list from registry (%v)", err)
	}

	return taskList, nil
}

func (p WindowsRegistryProvider) Close() {
}

func getUUIDFromTaskPath(path string) (string, error) {
	key, err := openKey(`Tree\` + path)
	if err != nil {
		return "", err
	}

	val, _, err := key.GetStringValue("Id")
	if err != nil {
		return "", err
	}

	return val, nil
}

func openKey(subKey string) (registry.Key, error) {
	return registry.OpenKey(registry.LOCAL_MACHINE, taskKeyBase+subKey, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
}

func openTaskKey(keyId string) registry.Key {
	var err error

	switch {
	case strings.HasPrefix(keyId, `\`):
		keyId, err = getUUIDFromTaskPath(keyId)
		if err != nil {
			log.Printf("cannot convert task path to uuid: %v\n", err)
			return 0
		}
		fallthrough

	case strings.HasPrefix(keyId, `{`):
		key, err := openKey(`Tasks\` + keyId)
		if err != nil {
			log.Printf("cannot open key %s: %v\n", keyId, err)
			return 0
		}
		return key

	default:
		log.Printf("task id unknown. must start with \\ or {")
		return 0
	}
}
