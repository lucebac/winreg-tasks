// SPDX-License-Identifier: MIT

package actions

import (
	"fmt"

	"github.com/lucebac/winreg-tasks/generated"
)

const MessageboxPropertiesMagic PropertiesMagic = 0x9999

type MessageboxProperties struct {
	Id string

	Caption string
	Content string
}

func NewMessageboxProperties(id string, gen *generated.Actions_MessageboxTaskProperties) (*MessageboxProperties, error) {
	return &MessageboxProperties{
		Id:      id,
		Caption: gen.Caption.Str,
		Content: gen.Content.Str,
	}, nil
}

func IsMessageboxProperties(properties Properties) bool {
	return properties.Magic() == MessageboxPropertiesMagic
}

func (m MessageboxProperties) Magic() PropertiesMagic {
	return MessageboxPropertiesMagic
}

func (e MessageboxProperties) Name() string {
	return "Messagebox"
}

func (m MessageboxProperties) String() string {
	return fmt.Sprintf(
		`<Messagebox id="%s" caption="%s" content="%s">`,
		m.Id, m.Caption, m.Content,
	)
}
