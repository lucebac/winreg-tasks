// SPDX-License-Identifier: MIT

package main

import (
	"github.com/alecthomas/kong"
)

var subcommands []kong.Option

func registerSubcommand(f kong.Option) {
	subcommands = append(subcommands, f)
}

type cli struct{}

func main() {
	ctx := kong.Parse(&cli{}, append(subcommands, kong.UsageOnError())...)
	ctx.FatalIfErrorf(ctx.Run())
}
