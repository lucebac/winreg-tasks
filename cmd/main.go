// SPDX-License-Identifier: MIT

package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/lucebac/winreg-tasks/providers"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var subcommands []kong.Option

func registerSubcommand(f kong.Option) {
	subcommands = append(subcommands, f)
}

type cli struct {
	File     *os.File   `help:"If provided, use this Hive file instead of the System's live one." short:"f" optional:""`
	LogFiles []*os.File `help:"If provided, these log files will be applied to the Hive file." short:"x" optional:""`

	LogLevel  string `help:"Set log level" short:"l" optional:"" default:"info" enum:"debug,info,warn,error"`
	LogFormat string `help:"Set the log format" optional:"" default:"plain" enum:"plain,json"`
}

type context struct {
	provider providers.DataProvider
}

func main() {
	var err error

	args := &cli{}
	kongContext := kong.Parse(args, append(subcommands, kong.UsageOnError())...)

	context := &context{}

	// logging setup
	if args.LogFormat == "plain" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	logLevel, err := zerolog.ParseLevel(args.LogLevel)
	kongContext.FatalIfErrorf(err)
	zerolog.SetGlobalLevel(logLevel)

	// setup data provider
	if args.File != nil {
		context.provider, err = providers.NewFileProvider(args.File, args.LogFiles...)
		kongContext.FatalIfErrorf(err)
	} else {
		context.provider, err = providers.GetNativeSystemProvider()
		kongContext.FatalIfErrorf(err)
	}

	// run command
	kongContext.FatalIfErrorf(kongContext.Run(context))
}
