package providers

import (
	"encoding/binary"
	"errors"
	"os"
	"time"

	"github.com/lucebac/winreg-tasks/utils"
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

func (p FileProvider) GetStringField(taskId, fieldName string) (string, error) {
	raw, err := p.getValueData(taskId, fieldName)
	if err != nil {
		return "", err
	}

	return utils.ConvertBytesToStringUTF16(raw)
}

func (p FileProvider) GetBytesField(taskId, fieldName string) ([]byte, error) {
	return p.getValueData(taskId, fieldName)
}

func (p FileProvider) GetDwordField(taskId, fieldName string) (uint32, error) {
	raw, err := p.getValueData(taskId, fieldName)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(raw[:4]), nil
}

func (p FileProvider) GetDateField(taskId, fieldName string) (*time.Time, error) {
	dateString, err := p.GetStringField(taskId, fieldName)
	if err != nil {
		return nil, err
	}
	return utils.ParseWindowsTimestamp(dateString)
}
