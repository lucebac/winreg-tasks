//go:build debug
// +build debug

// SPDX-License-Identifier: MIT

package main

import (
	"fmt"

	"github.com/alecthomas/kong"
)

type dummyCommand struct {
	Quiet  bool   `optional:"" default:"false" short:"q"`
	String string `arg:""`
}

func (d *dummyCommand) Run(ctx *kong.Context) error {
	fmt.Println(d.String)
	fmt.Printf("Bool: %v\n", d.Quiet)

	if d.Quiet {
		return fmt.Errorf("bool was true")
	}

	return nil
}

func init() {
	registerSubcommand(kong.DynamicCommand("dummy", "dummy command", "", &dummyCommand{}))
}
