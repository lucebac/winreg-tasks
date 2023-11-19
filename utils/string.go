package utils

import (
	"strings"

	"golang.org/x/text/encoding/unicode"
)

func ConvertBytesToStringUTF16(b []byte) (string, error) {
	enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	utf8bytes, err := enc.Bytes(b)
	if err != nil {
		return "", err
	}
	s := string(utf8bytes)
	return strings.TrimRight(s, "\x00"), nil
}
