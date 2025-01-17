// SPDX-License-Identifier: MIT

package actions

import (
	"fmt"

	"github.com/lucebac/winreg-tasks/generated"
)

const ExecutionPropertiesMagic PropertiesMagic = 0x6666

type ExecutionProperties struct {
	Id string

	Arguments        string
	Command          string
	WorkingDirectory string

	Flags uint16
}

func NewExecutionProperties(id string, gen *generated.Actions_ExeTaskProperties) (*ExecutionProperties, error) {
	return &ExecutionProperties{
		Id:               id,
		Arguments:        gen.Arguments.Str,
		Command:          gen.Command.Str,
		WorkingDirectory: gen.WorkingDirectory.Str,
		Flags:            gen.Flags,
	}, nil
}

func IsExecutionProperties(properties Properties) bool {
	return properties.Magic() == ExecutionPropertiesMagic
}

func (e ExecutionProperties) Magic() PropertiesMagic {
	return ExecutionPropertiesMagic
}

func (e ExecutionProperties) Name() string {
	return "Execution"
}

func (e ExecutionProperties) String() string {
	return fmt.Sprintf(
		`<Execution id="%s" command="%s" arguments="%s" workingDirectory="%s" flags=0x{%02x}`,
		e.Id, e.Command, e.Arguments, e.WorkingDirectory, e.Flags,
	)
}
