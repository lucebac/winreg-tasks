package providers

import (
	"errors"
	"os"

	"www.velocidex.com/golang/regparser"
)

type FileProvider struct {
	f   *os.File
	reg *regparser.Registry
}

var (
	ErrCannotOpenKey = errors.New("cannot open key")
	ErrValueNotFound = errors.New("value not found")
)

func NewFileProvider(f *os.File) (DataProvider, error) {
	reg, err := regparser.NewRegistry(f)
	if err != nil {
		return nil, err
	}

	return &FileProvider{
		reg: reg,
		f:   f,
	}, nil
}

func (p *FileProvider) Close() {
	p.reg = nil

	p.f.Close()
	p.f = nil
}

func (p FileProvider) getValueData(taskId, valueName string) ([]byte, error) {
	key := p.reg.OpenKey(taskKeyBase + taskId)
	if key == nil {
		return nil, ErrCannotOpenKey
	}

	values := key.Values()

	for _, val := range values {
		if val == nil {
			continue
		}

		if val.ValueName() != valueName {
			continue
		}

		return val.ValueData().Data, nil
	}

	return nil, ErrValueNotFound
}

func (p FileProvider) GetActions(taskId string) ([]byte, error) {
	return p.getValueData(taskId, "Actions")
}

func (p FileProvider) GetTriggers(taskId string) ([]byte, error) {
	return p.getValueData(taskId, "Triggers")
}

func (p FileProvider) GetDynamicInfo(taskId string) ([]byte, error) {
	return p.getValueData(taskId, "DynamicInfo")
}

func (p FileProvider) GetTaskIdList() ([]string, error) {
	key := p.reg.OpenKey(taskKeyBase + `Tasks`)
	if key == nil {
		return nil, ErrCannotOpenKey
	}

	subkeys := key.Subkeys()
	taskList := make([]string, len(subkeys))

	for i, k := range key.Subkeys() {
		taskList[i] = k.Name()
	}

	return taskList, nil
}
