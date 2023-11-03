// SPDX-License-Identifier: MIT

package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/lucebac/winreg-tasks/providers"
)

var subcommands []kong.Option

func registerSubcommand(f kong.Option) {
	subcommands = append(subcommands, f)
}

type cli struct {
	File *os.File `help:"If provided, use this Hive file instead of the System's live one." short:"f" type:"existingfile" optional:""`
}

type context struct {
	provider providers.DataProvider
}

func main() {
	var err error

	c := &cli{}
	ctx := kong.Parse(c, append(subcommands, kong.UsageOnError())...)

	context := &context{}

	if c.File != nil {
		context.provider, err = providers.NewFileProvider(c.File)
		ctx.FatalIfErrorf(err)
	} else {
		context.provider, err = providers.GetNativeSystemProvider()
		ctx.FatalIfErrorf(err)
	}

	ctx.FatalIfErrorf(ctx.Run(context))
}
