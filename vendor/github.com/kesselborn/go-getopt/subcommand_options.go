// Copyright (c) 2011, SoundCloud Ltd., Daniel Bornkessel
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/kesselborn/go-getopt

package getopt

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type SubCommands map[string]Options

type SubCommandOptions struct {
	Global      Options
	SubCommands SubCommands
}

func (sco SubCommandOptions) flattenToOptions(subCommand string) (options Options, err *GetOptError) {
	genericOptions := sco.Global

	if subCommandOptions, present := sco.SubCommands[subCommand]; present == true {

		if subCommand != "*" {
			for _, option := range genericOptions.Definitions {
				options.Definitions = append(options.Definitions, option)
			}
		}

		for _, option := range subCommandOptions.Definitions {
			options.Definitions = append(options.Definitions, option)
			options.Description = subCommandOptions.Description
		}
	} else {
		err = &GetOptError{UnknownSubCommand, "Unknown command: " + subCommand}
	}

	return
}

func (sco SubCommandOptions) findSubCommand() (subCommand string, err *GetOptError) {
	options := sco.Global
	subCommand = "*"

	_, arguments, _, _ := options.ParseCommandLine()

	if len(arguments) < 1 {
		err = &GetOptError{NoSubCommand, "Couldn't find sub command"}
	} else {
		subCommand = arguments[0]
	}

	return
}

func (sco SubCommandOptions) ParseCommandLine() (subCommand string, options map[string]OptionValue, arguments []string, passThrough []string, err *GetOptError) {
	if subCommand, err = sco.findSubCommand(); err == nil {
		options, arguments, passThrough, err = sco.parseCommandLineImpl(subCommand, os.Args[1:], mapifyEnvironment(os.Environ()), 0)
	} else {
		options, arguments, passThrough, err = sco.Global.parseCommandLineImpl(os.Args[1:], mapifyEnvironment(os.Environ()), 0)
	}

	return
}

func (sco SubCommandOptions) parseCommandLineImpl(subCommand string, args []string, environment map[string]string, flags int) (options map[string]OptionValue, arguments []string, passThrough []string, err *GetOptError) {

	var flattenedOptions Options
	if flattenedOptions, err = sco.flattenToOptions(subCommand); err == nil {
		options, arguments, passThrough, err = flattenedOptions.parseCommandLineImpl(args, environment, flags)
		arguments = arguments[1:]
	}

	return
}

func (sco SubCommandOptions) Usage() (output string) {
	return sco.UsageCustomArg0(filepath.Base(os.Args[0]))
}

func (sco SubCommandOptions) UsageCustomArg0(arg0 string) (output string) {
	subCommand, err := sco.findSubCommand()
	flattenedOptions, foundSubCommand := sco.SubCommands[subCommand]

	if err != nil || !foundSubCommand {
		output = sco.Global.UsageCustomArg0(arg0)
	} else {
		output = flattenedOptions.UsageCustomArg0(arg0 + " " + subCommand)
	}

	return
}

func (sco SubCommandOptions) Help() (output string) {
	return sco.HelpCustomArg0(filepath.Base(os.Args[0]))
}

func (sco SubCommandOptions) HelpCustomArg0(arg0 string) (output string) {
	subCommand, err := sco.findSubCommand()
	flattenedOptions, foundSubCommand := sco.SubCommands[subCommand]

	if err != nil || !foundSubCommand {
		subCommand = "*"
		flattenedOptions = sco.Global
	} else {
		arg0 = arg0 + " " + subCommand
	}

	output = flattenedOptions.HelpCustomArg0(arg0)

	if subCommand == "*" {
		// TODO: centralize format strings
		fmtStr := fmt.Sprintf("    %%-%ds       %%s\n", flattenedOptions.calculateLongOptTextLenght())

		output = output + "Available commands:\n"

		keys := make([]string, len(sco.SubCommands))
		i := 0

		for k := range sco.SubCommands {
			keys[i] = k
			i = i + 1
		}
		sort.Strings(keys)

		for _, key := range keys {
			output = output + fmt.Sprintf(fmtStr, key, sco.SubCommands[key].Description)
		}
		output = output + "\n"
	}

	return
}
