// SPDX-License-Identifier: MIT

package utils

import (
	"fmt"
	"strings"
)

// Hexdump converst the byte array into an ascii hex dump with lineLength elements per line
func Hexdump(data []byte, lineLength int) string {
	str := ""

	for i, c := range data {
		if i > 0 && i%lineLength == 0 {
			str = str[:len(str)-1] + "\n"
		}

		str += fmt.Sprintf("%02x ", c)
	}

	return strings.TrimSuffix(strings.TrimSuffix(str, "\n"), " ")
}
