// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidUtcOffset = errors.New(`invalid UTC offset`)
)

func toGolangHexBytes(data []byte) string {
	str := ""

	for _, c := range data {
		str += fmt.Sprintf(`0x%02x, `, c)
	}

	return strings.TrimSuffix(str, ` `)
}

func parseLocationRelativeUTC(s string) (*time.Location, error) {
	if len(s) < 5 {
		return nil, ErrInvalidUtcOffset
	}

	modifier := s[3]
	s_offset := s[4:]

	if s[:3] != "UTC" || (modifier != '+' && modifier != '-') {
		return nil, ErrInvalidUtcOffset
	}

	var hours, minutes int
	var err error

	switch len(s_offset) {
	case 1:
		hours, err = strconv.Atoi(s_offset[0:1])
		if err != nil {
			return nil, err
		}

	case 4:
		hours, err = strconv.Atoi(s_offset[0:2])
		if err != nil {
			return nil, err
		}

		minutes, err = strconv.Atoi(s_offset[2:4])
		if err != nil {
			return nil, err
		}

	default:
		return nil, ErrInvalidUtcOffset
	}

	if modifier == '+' {
		return time.FixedZone(s, hours*60+minutes), nil
	} else {
		return time.FixedZone(s, -(hours*60 + minutes)), nil
	}
}
