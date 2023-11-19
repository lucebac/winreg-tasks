package providers

import (
	"errors"
	"os"

	"www.velocidex.com/golang/regparser"
)

type FileProvider struct {
	f           *os.File
	rebuiltHive *os.File
	logFiles    []*os.File
	reg         *regparser.Registry
}

var (
	ErrCannotOpenKey = errors.New("cannot open key")
	ErrValueNotFound = errors.New("value not found")
)

func NewFileProvider(f *os.File, logFiles ...*os.File) (DataProvider, error) {
	rebuiltHive, err := regparser.RecoverHive(f, logFiles...)
	if err != nil {
		return nil, err
	}

	reg, err := regparser.NewRegistry(rebuiltHive)
	if err != nil {
		return nil, err
	}

	return &FileProvider{
		rebuiltHive: rebuiltHive,
		reg:         reg,
		f:           f,
		logFiles:    logFiles,
	}, nil
}

func (p *FileProvider) Close() {
	p.reg = nil

	p.f.Close()
	p.f = nil

	for _, f := range p.logFiles {
		f.Close()
	}
	p.logFiles = nil

	if p.rebuiltHive != nil {
		p.rebuiltHive.Close()
		os.Remove(p.rebuiltHive.Name())
		p.rebuiltHive = nil
	}
}

func (p FileProvider) getValueData(taskId, valueName string) ([]byte, error) {
	key := p.reg.OpenKey(taskKeyBase + `Tasks\` + taskId)
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
